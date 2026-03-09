package base

import (
	"fmt"
	"testing"
)

/*
channel 底层数据结构：buf核心数组（环形数组） + sendq、recvq两个阻塞队列（双向链表） + 一个互斥锁
type hchan struct {
	qcount   uint           // 当前缓冲区中的数据个数
	dataqsiz uint           // 环形缓冲区的大小（make 时定义的容量）
	buf      unsafe.Pointer // 指向环形缓冲区的指针（只对有缓冲 channel 有效）：环形数组（Ring Buffer），不是链表。数组在内存上是连续的，通过 sendx（写下标）和 recvx（读下标）来循环移动，效率极高。
	elemsize uint16         // 元素大小
	closed   uint32         // 是否关闭
	elemtype *_type        // 元素类型
	sendx    uint           // 写入数据的指针下标（环形索引）
	recvx    uint           // 读取数据的指针下标（环形索引）
	recvq    waitq          // 等待读的协程队列（Sudog 链表）双向链表
	sendq    waitq          // 等待写的协程队列（Sudog 链表）双向链表
	lock mutex // 互斥锁！保护 hchan 的所有字段
}

内存模型：不要通过共享内存来通信，而要通过通信来共享内存。
底层结构：hchan 结构体，包含环形队列、等待发送/接收的 sudog 链表、互斥锁。
关闭原则：由发送方关闭，不要在接收方关闭，也不要在有多个发送方时关闭。

三种 panic 的情况
操作	nil Channel (未 make)	Closed Channel (已关闭)  	Normal Channel (正常)
Close	Panic	               Panic					  正常关闭
Send	永久阻塞			     Panic	                    阻塞或进入缓冲区
Receive	永久阻塞				 读完旧数据后返回零值	        阻塞或取走数据

面试题：
1 channel 有什么要注意的？
 -给一个 nil 的 channel 发送或接收数据，会永久阻塞。
 -给一个已关闭（closed）的 channel 发送数据，会直接 panic。
 -从一个已关闭的 channel 接收数据，会先读完剩余数据，然后永远返回“零值”和 false。

2 Channel 为什么是线程安全的？
 -因为它的底层结构体 hchan 里自带了一把 Mutex（互斥锁）。每一次入队、出队或者操作等待队列时，Go 运行时都会先获取这把锁，保证同一时间只有一个协程能修改 Channel 的状态。

3 Channel 底层数据结构？
 -Channel 底层是一个叫 hchan 的结构体。它主要由三部分组成：
 - 环形缓冲区（buf）：这是一个数组，用来存带缓冲的 channel 数据。
 - 两个等待队列（recvq 和 sendq）：这是两个双向链表，里面存的是那些因为没数据读或者没地方发而阻塞的协程（sudog）。
 - 互斥锁（lock）：没错，Channel 内部也是靠锁来保证并发安全的，只不过它把锁封装得非常好，让我们用起来像是在直接传值。”

4 Channel 缓冲区 buf 底层数据结构？这么设计有什么好处？
 -数组实现的环形缓冲区可以复用内存。通过 sendx 和 recvx 两个索引移动，避免了像普通数组那样删除元素后需要移动后面所有数据的问题，保证了 O(1) 的处理速度

5 Channel 怎么会导致协程泄漏？
 -如果协程在等待一个永远没人发的 chan，或者发向一个没人收且没缓冲的 chan，该协程会永久阻塞在后台，无法被 GC 回收。

6 如何优雅地关闭 Channel？有多个发送者，一个接收者，怎么关闭 Channel？
 -直接关闭会 panic（因为可能有其他发送者还在发）。
 -方案 A：用 sync.Once 确保只关一次。
 -方案 B（更专业）：引入一个额外的 stopCh。接收者想停的时候，关闭 stopCh。所有的发送者通过 select 监听 stopCh，一旦发现 stopCh 关了，就停止发送。
 -核心原则：谁的数据源枯竭了，谁负责关闭。永远不要在接收端关闭，除非你是唯一的发送者。”

7 有缓冲和无缓冲的 channel 分别是啥？
 -无缓冲（Unbuffered）：是同步的。发送方和接收方必须“手递手”交接，性能开销最小（直接内存拷贝，不经过缓冲区）。
  -同步阻塞：它是强同步的。如果没有协程在听（接收），发送方会永久阻塞；反之亦然。
  -死锁高发区：在同一个 Goroutine 里又发又收，必然死锁。
  -性能优势：在底层，无缓冲 Channel 有时会触发直接内存拷贝（从发送协程的栈直接拷到接收协程的栈），不需要经过中间缓冲区，延迟极低。
 -有缓冲（Buffered）：是异步的。数据先拷贝进 hchan.buf，再由接收方拷走。

8 协程间的直接拷贝有什么好处？
 -对于无缓冲 Channel，Go 优化到了极致：直接从 A 协程的内存拷到 B 协程的内存，中间不落地。
 -按照常规逻辑，数据传递应该是： 发送者协程 -> 拷贝到 Channel 缓冲区 -> 接收者协程从缓冲区拷贝走。 这涉及到 2 次 内存拷贝，效率不高。
 -Go 的 Channel 设计成：发送者协程 -> 直接拷贝到 接收者协程。当发现recvq又协程等待时，直接拷贝数据到等待协程的内存地址，不经过缓冲区，这样只需要 1 次 内存拷贝，极大提升了性能，尤其是对于大数据量的传递。

9 select 是什么？
 -select 是一种 多路复用（I/O Multiplexing） 机制。它会阻塞当前协程，直到监听的多个 Channel 中至少有一个可以进行读写为止。
 -随机性：多个 case 同时就绪时，随机选一个分支执行，然后整个 select 就结束了，随机是为了保证公平性，防止接
 -原子性：select 对 Channel 的操作是带锁的（hchan.lock），保证了并发安全。
 -单次触发：一次 select 只要有一个 case 成功，其他的 case 就会被撤销。

10 select 底层是怎么从 Channel 读数据的？
 -1 先把 case 顺序随机打乱，保证随机读取和防止饥饿
 -2 按照随机后的顺序，遍历所有的 case，检查 channel，如果有数据则执行分支并返回结果，没数据如果有 default 则执行 default
 -3 如果所有 case 都没有数据，且没有 default 分支：
    -1）打包并入队：当前协程（G）会被打包成多个 sudog 结构，同时放入所有 case 对应的 Channel 的等待接收队列（recvq）中。
    -2）休眠：协程调用 gopark 进入休眠状态，出让 CPU。
	-3）唤醒：只要这几个 Channel 中任何一个有数据进来了，那个 Channel 的发送者就会把数据交给对应的 sudog，并唤醒该协程。
    -4）清理：协程醒来后，会把自已从其他那几个还没动静的 Channel 的等待队列中移除。

11 空的 select 会怎样？(select {})
 -这行代码会直接导致当前协程永久阻塞。因为里面一个 case 都没有，它没法被唤醒。
 -如果这段代码写在 main 函数里，Go 还会触发恐慌报 fatal error: all goroutines are asleep - deadlock!。
 -有时候我们在写测试或者简单的 Demo 时，想让主协程不退出，会随手写个这，但生产环境里一定要慎用

12 select 如何实现“带优先级”的读取？
 -这就得用‘嵌套 select’的小技巧了。
 -单个 select 保证不了顺序，所以我们可以先用一个单路的 select 尝试读 ch1，带上 default 防止阻塞；如果 ch1 没东西，再进第二个正常的 select 同时监听 ch1 和 ch2。
 -这样即使两个都有数据，第一个 select 也会稳稳地把 ch1 截获。”
*/

// 标准遍历：展示 for range 的正确用法
func TestChanBasics(t *testing.T) {
	ch := make(chan string, 3)

	go func() {
		ch <- "Golang"
		ch <- "Channel"
		ch <- "Practice"
		// 发送完一定要关，否则 main 的 for range 会永久阻塞，导致死锁。
		close(ch)
	}()

	// range 会一直读，直到 ch 被关闭且数据取完
	for v := range ch {
		fmt.Println("接收到:", v)
	}
	fmt.Println("遍历正常结束")
}

// 关闭后的读取：展示 (value, ok) 模式
func TestChanReadAfterClosed(t *testing.T) {
	ch := make(chan int, 1)
	ch <- 99
	close(ch)

	// 第一次读：能读到缓冲区剩下的值
	val1, ok1 := <-ch
	fmt.Printf("  第一次读 - 值: %d, ok: %v (说明通道虽关，但数据还在)\n", val1, ok1)

	// 第二次读：缓冲区空了，且通道已关
	val2, ok2 := <-ch
	fmt.Printf("  第二次读 - 值: %d, ok: %v (说明已关且空，读到零值)\n", val2, ok2)
}

// 各种会导致 Panic 的雷区 (手动测试用)
func TestChanSpecialCases(t *testing.T) {
	// --- 情况 A: 向已关闭的通道发数据 ---
	// c1 := make(chan int)
	// close(c1)
	// c1 <- 1 // 会引发 panic: send on closed channel

	// --- 情况 B: 重复关闭通道 ---
	// c2 := make(chan int)
	// close(c2)
	// close(c2) // 会引发 panic: close of closed channel

	// --- 情况 C: 关闭 nil 通道 ---
	// var c3 chan int
	// close(c3) // 会引发 panic: close of nil channel

	// --- 情况 D: 读 nil 通道 (永久阻塞导致死锁) ---
	// var c4 chan int
	// <-c4 // fatal error: all goroutines are asleep - deadlock!

	fmt.Println("  (此函数内的 Panic 示例代码已注释，可手动解开测试)")
}

// 如何实现非阻塞读写，多路监听器
func TestChanSelect(t *testing.T) {
	ch := make(chan string)

	// 1. 非阻塞尝试读取
	select {
	case msg := <-ch:
		fmt.Println("读到了:", msg)
	default:
		fmt.Println("没数据，我不等，直接走")
	}

	// 2. 非阻塞尝试发送
	select {
	case ch <- "hello":
		fmt.Println("发送成功")
	default:
		fmt.Println("没人收，发不出去，我也直接走")
	}

	// 3. 超时设定
	// select {
	// case res := <-ch:
	// 	fmt.Println("收到结果:", res)
	// case <-time.After(1 * time.Second): // 这就是“超时限制”
	// 	fmt.Println("超时退出：等了 1 秒还没好")
	// }
}

// 优先级访问
func TestPrioritySelect(t *testing.T) {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)
	ch1 <- 1
	ch2 <- 2

	// 需求：优先读 ch1
	select {
	case v := <-ch1:
		fmt.Println("优先拿到了 ch1:", v)
	default:
		// 如果 ch1 没好，才公平竞争
		select {
		case v := <-ch1:
			fmt.Println("拿到了 ch1:", v)
		case v := <-ch2:
			fmt.Println("拿到了 ch2:", v)
		}
	}
}

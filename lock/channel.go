package main

import (
	"fmt"
	"time"
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

三种 panic 的情况
操作	nil Channel (未 make)	Closed Channel (已关闭)  	Normal Channel (正常)
Close	Panic	               Panic					  正常关闭
Send	永久阻塞			     Panic	                    阻塞或进入缓冲区
Receive	永久阻塞				  读完旧数据后返回零值	        阻塞或取走数据

Channel 为什么是线程安全的？
 -因为它的底层结构体 hchan 里自带了一把 Mutex（互斥锁）。每一次入队、出队或者操作等待队列时，Go 运行时都会先获取这把锁，保证同一时间只有一个协程能修改 Channel 的状态。

Channel 缓冲区 buf 底层数据结构？这么设计有什么好处？
 -数组实现的环形缓冲区可以复用内存。通过 sendx 和 recvx 两个索引移动，避免了像普通数组那样删除元素后需要移动后面所有数据的问题，保证了 O(1) 的处理速度

协程间的直接拷贝有什么好处？
 -按照常规逻辑，数据传递应该是： 发送者协程 -> 拷贝到 Channel 缓冲区 -> 接收者协程从缓冲区拷贝走。 这涉及到 2 次 内存拷贝，效率不高。
 -Go 的 Channel 设计成：发送者协程 -> 直接拷贝到 接收者协程。当发现recvq又协程等待时，直接拷贝数据到等待协程的内存地址，不经过缓冲区，这样只需要 1 次 内存拷贝，极大提升了性能，尤其是对于大数据量的传递。

*/

func main() {
	fmt.Println("=== 1. 标准遍历演示 ===")
	rangeNormal()

	fmt.Println("\n=== 2. 有缓冲 vs 无缓冲阻塞演示 ===")
	compareBlocking()

	fmt.Println("\n=== 3. 已关闭 Channel 的读取行为 ===")
	readAfterClosed()

	fmt.Println("\n=== 4. 特殊情况演示（Panic 或 阻塞） ===")
	channelSpecialCases()
}

// 1. 标准遍历：展示 for range 的正确用法
func rangeNormal() {
	ch := make(chan string, 3)

	go func() {
		ch <- "Golang"
		ch <- "Channel"
		ch <- "Practice"
		// 发送完一定要关，否则 main 的 for range 会死锁
		close(ch)
	}()

	// range 会一直读，直到 ch 被关闭且数据取完
	for v := range ch {
		fmt.Println("接收到:", v)
	}
	fmt.Println("遍历正常结束")
}

// 2. 阻塞对比：展示无缓冲的“同步”和有缓冲的“异步”
func compareBlocking() {
	// 无缓冲：必须手递手
	unbuf := make(chan int)
	go func() {
		fmt.Println("  [无缓冲] 协程准备发送...")
		unbuf <- 1 // 阻塞，直到 main 读它
		fmt.Println("  [无缓冲] 发送成功")
	}()

	time.Sleep(500 * time.Millisecond)
	fmt.Println("  [无缓冲] 主协程准备接收")
	<-unbuf

	// 有缓冲：缓冲区没满就不阻塞
	buf := make(chan int, 2)
	fmt.Println("  [有缓冲] 发送前两个数据...")
	buf <- 10
	buf <- 20
	fmt.Println("  [有缓冲] 前两个发送成功，没阻塞")
}

// 3. 关闭后的读取：展示 (value, ok) 模式
func readAfterClosed() {
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

// 4. 各种会导致 Panic 的雷区 (手动测试用)
func channelSpecialCases() {
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

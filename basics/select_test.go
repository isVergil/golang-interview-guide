package basics

import (
	"fmt"
	"testing"
)

/*
面试题：
1 select 是什么？
 -select 是一种 多路复用（I/O Multiplexing） 机制。它会阻塞当前协程，直到监听的多个 Channel 中至少有一个可以进行读写为止。
 -随机性：多个 case 同时就绪时，随机选一个分支执行，然后整个 select 就结束了，随机是为了保证公平性，防止接
 -原子性：select 对 Channel 的操作是带锁的（hchan.lock），保证了并发安全。
 -单次触发：一次 select 只要有一个 case 成功，其他的 case 就会被撤销。

2 select 底层是怎么从 Channel 读数据的？
 -1 先把 case 顺序随机打乱，保证随机读取和防止饥饿
 -2 按照随机后的顺序，遍历所有的 case，检查 channel，如果有数据则执行分支并返回结果，没数据如果有 default 则执行 default
 -3 如果所有 case 都没有数据，且没有 default 分支：
    -1）打包并入队：当前协程（G）会被打包成多个 sudog 结构，同时放入所有 case 对应的 Channel 的等待接收队列（recvq）中。
    -2）休眠：协程调用 gopark 进入休眠状态，出让 CPU。
	-3）唤醒：只要这几个 Channel 中任何一个有数据进来了，那个 Channel 的发送者就会把数据交给对应的 sudog，并唤醒该协程。
    -4）清理：协程醒来后，会把自已从其他那几个还没动静的 Channel 的等待队列中移除。

3 空的 select 会怎样？(select {})
 -这行代码会直接导致当前协程永久阻塞。因为里面一个 case 都没有，它没法被唤醒。
 -如果这段代码写在 main 函数里，Go 还会触发恐慌报 fatal error: all goroutines are asleep - deadlock!。
 -有时候我们在写测试或者简单的 Demo 时，想让主协程不退出，会随手写个这，但生产环境里一定要慎用

4 select 如何实现"带优先级"的读取？
 -这就得用'嵌套 select'的小技巧了。
 -单个 select 保证不了顺序，所以我们可以先用一个单路的 select 尝试读 ch1，带上 default 防止阻塞；如果 ch1 没东西，再进第二个正常的 select 同时监听 ch1 和 ch2。
 -这样即使两个都有数据，第一个 select 也会稳稳地把 ch1 截获。"
*/

// 如何实现非阻塞读写，多路监听器
func TestSelect(t *testing.T) {
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
	// case <-time.After(1 * time.Second): // 这就是"超时限制"
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

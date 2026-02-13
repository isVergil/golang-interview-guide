package main

/*
waitgroup 的底层结构:状态变量+信号量+等待队列
它的源码在 sync/waitgroup.go 中，结构如下：
type WaitGroup struct {
    noCopy noCopy // 1. 看到没？这里又出现了 noCopy，说明 WaitGroup 不能拷贝！
    state  atomic.Uint64 // 2. 这是一个 64 位的值，高 32 位存“计数”，低 32 位存“等待者数量”
    sema   uint32        // 3. 又是它！信号量，用于阻塞和唤醒主协程
}

底层原理：用 state 这个 64 位整数同时存储两个计算器 任务计数器 和 等待者数量
1 计数器：高 32 位存储“待完成任务数”，也就是 Add 和 Done 操作修改的值。
2 等待者数量：低 32 位存储“等待者数量”，也就是调用 Wait 的协程数量。通常只有主协程会调用 Wait，所以这个值一般是 0 或 1。
3 信号量：当计数器归零时，Wait 会通过信号量唤醒等待的协程。

为什么要维护两个计数器呢：state是atomic.Uint64类型的，通过一次原子操作同时管理“运行的任务Counter”和“等待的任务Waiter Counter”。
为什么 Mutex 锁的 state 是 32 位的 int32 而 WaitGroup 的 state 是 atomic.Uint64 类型：
	-Mutex 很省空间，它把加锁状态、饥饿状态等都压缩在一个 int32 的不同‘位’上，4 字节就足够了，而且位运算非常快。
    -WaitGroup 必须同时维护两个信息：‘干活的人数’和‘等结果的人数’。为了保证这两个数在修改时是绝对同步的（防止出现这边刚干完活，那边还没记录上等待者的情况），Go 把它们并排放在了一个 64 位整数里。
    -这样，通过一次 Uint64 的原子操作，就能同时更新这两个计数器，避免了使用额外的锁来保护计数器，性能达到了最优。”

三大核心方法
Add(delta int)：把“待完成任务数”加 n。通常在启动协程前调用。
Done()：把“待完成任务数”减 1。通常在协程内部逻辑结束时调用（内部其实就是调了 Add(-1)）。
Wait()：阻塞当前协程，直到计数器归零。

WaitGroup 使用注意事项
1 计数器为负（Panic）：如果调 Done() 的次数比 Add() 多，直接报 panic: sync: negative WaitGroup counter。
2 Add 位置不对：如果在 go worker() 里面才调 Add(1)，可能主协程运行太快，还没来得及加 1 就执行到 Wait() 了，直接就退出了。必须在 go 语句外 Add。
3 没有传指针：如果你把 WaitGroup 按值传给协程，协程内部改的是副本，主协程的计数器永远不会变，结果就是永久死锁。
4 重用 WaitGroup 的风险：在一个 Wait() 还没结束时，又尝试 Add()，会触发 Panic。
*/

import (
	"fmt"
	"sync"
	"time"
)

func worker(id int, wg *sync.WaitGroup) {
	// 3. 任务结束时调用 Done
	defer wg.Done()

	fmt.Printf("工人 %d 开始干活...\n", id)
	time.Sleep(time.Second)
	fmt.Printf("工人 %d 干完活了！\n", id)
}

func main() {
	var wg sync.WaitGroup

	for i := 1; i <= 3; i++ {
		// 1. 启动协程前 Add，表示点名册加一人
		wg.Add(1)
		go worker(i, &wg) // 2. 注意：必须传指针，否则会发生锁拷贝！
	}

	fmt.Println("主协程：等着那帮家伙干完...")

	// 4. 阻塞在这里，不再是 select{} 的死等，而是等计数器归零
	wg.Wait()

	fmt.Println("主协程：全干完了，收工回家！")
}

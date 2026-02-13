package main

/*
1. 什么是 GMP？各个组件的作用是什么？
GMP 是 Go 语言核心的并发调度模型，它解决了原生线程（Kernel Thread）开销过大、上下文切换成本高的问题，支撑了 Go 语言的高并发特性。
G (Goroutine)：协程。它是轻量级的线程，只占用 2KB 左右的栈空间。它存储了运行所需的栈、状态和指针。
M (Machine)：内核线程。它是真正的执行单元。M 不直接运行 G，而是通过绑定 P 来执行 G。
P (Processor)：处理器/上下文。它包含了运行 G 所需的资源和本地队列。P 的数量决定了系统并行的 G 的上限（由 runtime.GOMAXPROCS 决定）。

2. 为什么要有 P？只用 G 和 M 行不行？
早期 Go（1.1 之前）只有 GM。引入 P 是为了解决：
	单一全局锁限制：GM 模型下，所有 M 争抢全局 G 队列，锁竞争剧烈。
	Cache 亲和性差：当 M 阻塞后，与其相关的 G 可能会被其他 M 领走，导致 CPU 缓存失效。
	引入 P 后：每个 P 维护本地队列，大大减少了全局锁争抢。

3. 调度器的策略有哪些？
 -复用线程 (Work Stealing)：当某个 P 的本地队列为空时，它会从其他 P 的队列中“偷”一半的 G 过来执行，避免线程闲置。
 -利用率优化、任务剥离 (Hand Off)：当 M 执行 G 时发生了系统调用阻塞，M 会释放绑定的 P，让 P 去寻找（或创建）其他的 M 继续工作，保证 CPU 不被闲置。
 -抢占式调度和协作式调度
   1.14 之前：基于协作。如果一个 G 运行超过 10ms 且发生函数调用，调度器会尝试将其挂起。
   1.14 之后：基于信号（Signal）的异步抢占。即使是死循环，系统也会通过信号强制挂起 G，解决了老版本中死循环导致程序卡死的问题。
 -全局队列 (Global Queue)：作为本地队列的补充。当本地队列满员或 P 无法从其他地方偷到任务时，会访问全局队列。为了保证公平性，P 每执行 61 次调度，就会强制检查一次全局队列，防止其中的 G 被“饿死”。

4. 为什么 GMP 性能高？
 -减少了内核态/用户态切换：G 的调度完全在用户态完成，不涉及 CPU 的模式切换。
 -更小的栈消耗：相比原生线程数 MB 的栈，G 的动态扩容栈极其节省内存。
 -降低了锁竞争：由于引入了 P，大部分 G 在 P 的本地队列中被处理，避免了多个线程同时争抢同一个全局任务锁的尴尬。

5. 一个 G 的生命周期是怎样的？
 -创建：go func() 创建 G，优先放入当前 P 的本地队列。
 -执行：M 从 P 获取 G 并运行。
 -阻塞：若发生通道 (Channel) 阻塞，G 进入等待状态，M 寻找下一个 G。
 -系统调用：M 阻塞，P “剥离” M 并寻找新 M。
 -销毁：执行完毕，进入 P 的自由列表（free list）待复用。

6. M 的数量和 P 的数量有限制吗？
 -P 的数量：默认等于 CPU 核数。可以通过 runtime.GOMAXPROCS 动态修改。
 -M 的数量：Go 默认限制最大 10,000 个。但通常只有活跃的 M 等于 P 的数量，多出的 M 往往处于休眠或执行系统调用的状态。

7. G 什么时候会被放入全局队列？
 -本地队列已满（默认 256 个）。
 -新建 G 时，若本地队列放不下，会把一半 G 打包放进全局队列。
 -从系统调用返回后，找不到空闲 P 绑定，该 G 也会进入全局队列。

8. 调度器如何保证公平性？
 -每 61 次调度，M 会强制从全局队列获取 G，防止全局队列中的 G 被“饿死”。

*/

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	// 设置使用4个逻辑处理器
	runtime.GOMAXPROCS(4)

	fmt.Printf("GOMAXPROCS: %d\n", runtime.GOMAXPROCS(0))
	fmt.Printf("NumCPU: %d\n", runtime.NumCPU())
	fmt.Printf("初始Goroutine数量: %d\n", runtime.NumGoroutine())

	// 演示1: 基础Goroutine创建和调度
	demoBasicScheduling()

	// 演示2: 工作窃取机制
	demoWorkStealing()

	// 演示3: 系统调用阻塞
	demoSyscallBlocking()

	// 持续监控
	monitorGPM()
}

func demoBasicScheduling() {
	fmt.Println("\n=== 演示1: 基础Goroutine创建和调度 ===")

	var wg sync.WaitGroup

	// 创建多个Goroutine
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Printf("Goroutine %d 开始执行 - P数量: %d, M数量: %d\n",
				id, getPCount(), getMCount())

			// 模拟计算工作
			sum := 0
			for j := 0; j < 1000000; j++ {
				sum += j
			}

			fmt.Printf("Goroutine %d 执行完成\n", id)
		}(i)
	}

	wg.Wait()
	fmt.Printf("演示1结束 - 当前Goroutine数量: %d\n", runtime.NumGoroutine())
}

func demoWorkStealing() {
	fmt.Println("\n=== 演示2: 工作窃取机制演示 ===")

	// 创建一个P负载不均衡的场景
	var wg sync.WaitGroup
	ch := make(chan int, 20)

	// 生产者 - 快速创建任务
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 20; i++ {
			ch <- i
			fmt.Printf("生产任务 %d - Goroutine数量: %d\n", i, runtime.NumGoroutine())
			time.Sleep(10 * time.Millisecond) // 让调度器有时间调度
		}
		close(ch)
	}()

	// 消费者 - 多个Goroutine处理任务
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for task := range ch {
				fmt.Printf("Worker %d 处理任务 %d - P: %d, M: %d\n",
					workerID, task, getPCount(), getMCount())

				// 模拟处理时间差异
				if workerID%2 == 0 {
					time.Sleep(50 * time.Millisecond)
				} else {
					time.Sleep(20 * time.Millisecond)
				}
			}
		}(i)
	}

	wg.Wait()
	fmt.Println("工作窃取演示完成")
}

func demoSyscallBlocking() {
	fmt.Println("\n=== 演示3: 系统调用阻塞处理 ===")

	var wg sync.WaitGroup
	startMCount := getMCount()

	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			if id%2 == 0 {
				// 模拟系统调用阻塞
				fmt.Printf("Goroutine %d 即将执行系统调用 - 当前M数量: %d\n", id, getMCount())
				time.Sleep(200 * time.Millisecond) // 模拟系统调用
				fmt.Printf("Goroutine %d 系统调用完成\n", id)
			} else {
				// CPU密集型任务
				fmt.Printf("Goroutine %d 执行CPU密集型任务\n", id)
				computeIntensiveWork()
			}
		}(i)
	}

	wg.Wait()
	endMCount := getMCount()
	fmt.Printf("系统调用演示完成 - M数量变化: %d -> %d\n", startMCount, endMCount)
}

func computeIntensiveWork() {
	sum := 0
	for i := 0; i < 10000000; i++ {
		sum += i
	}
}

func monitorGPM() {
	fmt.Println("\n=== GPM状态监控 (持续10秒) ===")

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	timeout := time.After(10 * time.Second)

	for {
		select {
		case <-ticker.C:
			fmt.Printf("状态 - Goroutines: %d, P: %d, M: %d\n",
				runtime.NumGoroutine(), getPCount(), getMCount())

		case <-timeout:
			fmt.Println("监控结束")
			return
		}
	}
}

// 获取当前活跃的P数量
func getPCount() int {
	var count int
	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		// 通过创建微小任务来探测P状态
		done := make(chan bool)
		go func() {
			done <- true
		}()
		<-done
		count++
	}
	return count
}

// 获取当前活跃的M数量（近似值）
func getMCount() int {
	// 通过创建多个短暂Goroutine来观察M创建
	before := runtime.NumGoroutine()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
		}()
	}
	wg.Wait()

	// 这是一个简化的估算，实际M数量需要更复杂的监控
	return runtime.NumGoroutine() - before + 1
}

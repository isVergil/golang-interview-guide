package basics

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

/*
Q1: goroutine 是什么？和线程有什么区别？
Q2: GMP 模型是什么？G、M、P 分别是什么？
Q3: goroutine 的调度时机？什么时候会切换？
Q4: goroutine 的栈是怎么管理的？
Q5: 抢占式调度是怎么实现的？
Q6: Work Stealing 机制是什么？
Q7: goroutine 泄漏是什么？怎么排查和避免？
Q8: GOMAXPROCS 是什么？设置多少合适？
Q9: goroutine 和 channel 怎么配合实现并发模式？
Q10: goroutine 的生命周期？状态有哪些？
Q11: Hand Off 机制是什么？G 什么时候会放入全局队列？
Q12: sysmon 是什么？调度器怎么保证公平性？

---
Q1: goroutine 是什么？和线程有什么区别？
【理解】
goroutine 是 Go 的用户态轻量级线程，由 Go runtime 调度，不是 OS 线程。

对比 OS 线程：
  维度          goroutine              OS 线程
  创建成本      ~2KB 栈，纳秒级创建     ~1MB 栈，微秒级创建
  切换成本      用户态切换（~几十ns）    内核态切换（~几μs）
  数量          轻松百万级               通常几千就是极限
  调度          Go runtime 调度         OS 内核调度
  栈大小        动态伸缩（2KB~1GB）     固定大小（默认 1-8MB）

为什么这么轻量？
  1. 栈小：初始只有 2KB（线程 1MB+），按需扩容
  2. 用户态调度：不需要陷入内核，没有系统调用开销
  3. 复用线程：M 个 goroutine 映射到 N 个线程（M:N 模型），线程数远少于协程数
  4. 切换只保存少量寄存器（SP、PC、BP），不需要保存完整 CPU 上下文

本质：goroutine 是 runtime 管理的协程，OS 线程是内核管理的执行单元。
Go 用少量线程驱动大量协程，兼顾了并发能力和系统资源效率。
【回答】
goroutine 是 Go runtime 管理的用户态轻量级协程，初始栈只有 2KB，创建和切换都是纳秒级，轻松开百万个。
和 OS 线程的核心区别：线程由内核调度，栈固定 1MB+，切换要陷入内核（微秒级）；goroutine 由 Go runtime 在用户态调度，切换只保存几个寄存器，极其轻量。
Go 用 M:N 模型，少量 OS 线程驱动大量 goroutine，兼顾并发能力和资源效率。

---
Q2: GMP 模型是什么？G、M、P 分别是什么？
【理解】
GMP 是 Go 调度器的核心模型：

G（Goroutine）：协程，包含栈、PC、状态等信息。就是 runtime.g 结构体。
M（Machine）：OS 线程，真正执行代码的载体。一个 M 同一时刻只能执行一个 G。
P（Processor）：逻辑处理器，持有本地运行队列（local run queue）。
  P 的数量 = GOMAXPROCS（默认等于 CPU 核数）。
  M 必须绑定一个 P 才能执行 G。

关系：
  G 挂在 P 的本地队列上等待执行
  M 绑定 P 后，从 P 的队列取 G 执行
  P 的数量决定了真正的并行度

为什么需要 P？（Go 1.1 引入）
  Go 1.0 只有 G 和 M，全局队列 + 全局锁，多线程竞争严重。
  引入 P 后：
    1. 每个 P 有本地队列，M 优先从本地取 G，减少锁竞争
    2. P 持有 mcache（内存分配缓存），减少内存分配锁竞争
    3. M 阻塞时可以把 P 让给其他 M，不浪费 CPU

调度循环：
  M 绑定 P -> 从 P 本地队列取 G -> 执行 G -> G 完成/阻塞 -> 取下一个 G
  本地队列空了 -> 从全局队列偷 -> 从其他 P 偷（work stealing）
【回答】
GMP 是 Go 调度器的核心：
G 是 goroutine，M 是 OS 线程，P 是逻辑处理器（持有本地运行队列）。
M 必须绑定 P 才能执行 G，P 的数量等于 GOMAXPROCS（默认 CPU 核数），决定了真正的并行度。
为什么需要 P：Go 1.0 只有全局队列，锁竞争严重。引入 P 后每个 P 有本地队列，M 优先从本地取 G，大幅减少锁竞争。M 阻塞时 P 可以转给其他 M，不浪费 CPU。

---
Q3: goroutine 的调度时机？什么时候会切换？
【理解】
Go 调度器是协作式 + 抢占式混合调度。

主动让出（协作式）：
  1. channel 操作阻塞（发送/接收）
  2. 系统调用（文件 IO、网络 IO）
  3. time.Sleep
  4. mutex 锁竞争
  5. runtime.Gosched() 主动让出
  6. select 阻塞
  7. GC 的 STW

被动抢占（抢占式，Go 1.14+）：
  1. sysmon 监控线程发现 G 运行超过 10ms，设置抢占标记
  2. Go 1.14 之前：等 G 执行函数调用时检查标记（协作式抢占）
     问题：纯计算的 for 循环没有函数调用，永远不会被抢占
  3. Go 1.14+：基于信号的异步抢占（发送 SIGURG 信号）
     即使没有函数调用也能被抢占，解决了"死循环饿死"问题

系统调用时的特殊处理：
  G 进入系统调用 -> M 被阻塞 -> P 和 M 解绑 -> P 找另一个空闲 M 继续执行其他 G
  系统调用返回 -> M 尝试重新获取 P -> 没有空闲 P 则 G 放入全局队列，M 休眠
【回答】
两种调度方式：
协作式：channel 阻塞、系统调用、锁竞争、time.Sleep、runtime.Gosched() 等都会触发切换。
抢占式（Go 1.14+）：sysmon 发现 G 运行超过 10ms，发送 SIGURG 信号强制抢占，即使纯计算循环也能被打断。
系统调用时 M 被阻塞，P 会解绑转给其他 M 继续工作，不浪费 CPU。

---
Q4: goroutine 的栈是怎么管理的？
【理解】
goroutine 栈是动态伸缩的：初始 2KB，最大可达 1GB。

■ 栈扩容（连续栈，Go 1.4+）：
  每个函数入口有栈溢出检查（比较 SP 和 stackguard0）。
  如果栈空间不够：
    1. 分配一个 2 倍大小的新栈
    2. 把旧栈内容拷贝到新栈
    3. 调整所有指向旧栈的指针（因为地址变了）
    4. 释放旧栈

  Go 1.3 之前用"分段栈"（segmented stack）：
    栈不够时分配新段，用链表连接。
    问题：hot split——函数在栈边界反复调用导致频繁分配/释放段，性能差。
  Go 1.4+ 改为"连续栈"（contiguous stack）：
    直接拷贝到更大的连续内存，避免 hot split。

■ 栈缩容：
  GC 时检查栈使用率，如果使用不到 1/4，缩容到一半。
  防止曾经用过大栈的 goroutine 一直占着内存。

为什么初始只有 2KB？
  大多数 goroutine 不需要大栈，2KB 够用。
  百万 goroutine × 2KB = 2GB，可以接受。
  如果像线程一样 1MB，百万级就需要 1TB，不可能。
【回答】
goroutine 栈初始 2KB，动态伸缩，最大 1GB。
扩容：函数入口检查栈空间，不够就分配 2 倍新栈，拷贝旧内容并调整指针。Go 1.4+ 用连续栈替代了分段栈，避免 hot split 性能问题。
缩容：GC 时检查栈使用率，低于 1/4 就缩到一半，防止内存浪费。
初始 2KB 的设计让百万 goroutine 只需 2GB 内存，如果像线程 1MB 就需要 1TB。

---
Q5: 抢占式调度是怎么实现的？
【理解】
Go 的抢占经历了两个阶段：

■ Go 1.2~1.13：协作式抢占
  sysmon 后台线程每 10ms 检查一次，发现 G 运行超过 10ms：
    设置 g.stackguard0 = stackPreempt（一个特殊值）
  G 在下次函数调用时检查栈空间，发现 stackguard0 是抢占标记：
    主动调用 schedule() 让出 CPU
  问题：纯计算循环（for { i++ }）没有函数调用，永远不检查，永远不让出。

■ Go 1.14+：基于信号的异步抢占
  sysmon 发现 G 运行超过 10ms：
    向 M 发送 SIGURG 信号
  M 收到信号后：
    1. 信号处理函数保存当前 G 的寄存器状态
    2. 把当前 G 的 PC 修改为 asyncPreempt 函数地址
    3. 信号处理返回后，G "恢复执行"时实际跳到 asyncPreempt
    4. asyncPreempt 调用 schedule()，G 被放回队列
  效果：即使没有函数调用也能被抢占，解决了死循环饿死问题。

sysmon 是什么？
  独立的监控线程，不绑定 P，不参与调度。
  职责：抢占检测、网络轮询（netpoll）、强制 GC、定时器检查。
【回答】
Go 1.14 之前是协作式抢占：sysmon 设置抢占标记，G 在函数调用时检查并让出。问题是纯计算循环没有函数调用就永远不让出。
Go 1.14+ 改为基于信号的异步抢占：sysmon 发现 G 运行超过 10ms，向 M 发送 SIGURG 信号，信号处理函数修改 G 的 PC 指向 asyncPreempt，G "恢复"时实际跳到调度函数让出 CPU。
这样即使死循环也能被抢占，彻底解决了饿死问题。

---
Q6: Work Stealing 机制是什么？
【理解】
Work Stealing = 工作窃取，是 P 本地队列为空时的负载均衡策略。

流程：
  当 M 绑定的 P 本地队列为空时：
    Step1: 先从全局队列取一批 G（全局队列有锁，取一批减少锁竞争）
    Step2: 全局队列也空了 -> 随机选一个其他 P，偷走它本地队列一半的 G
    Step3: 所有地方都没有 G -> M 和 P 解绑，M 进入休眠（放入空闲 M 列表）

偷取策略：
  随机选目标 P（避免总偷同一个导致不均衡）
  偷一半（不是全偷，保持两边都有活干）

本地队列的结构：
  大小固定 256 个槽位的环形数组（lock-free）
  超过 256 的 G 放入全局队列

为什么需要 Work Stealing？
  如果只有全局队列：所有 M 竞争一把锁，性能差。
  有了本地队列：大部分时候无锁操作，只在偷取时才涉及同步。
  Work Stealing 保证了负载均衡——忙的 P 不会堆积，闲的 P 不会空转。
【回答】
Work Stealing 是 P 本地队列为空时的负载均衡策略。
流程：本地队列空 -> 先从全局队列取一批 -> 全局也空则随机偷另一个 P 的一半 G -> 都没有则 M 休眠。
偷一半而不是全偷，保持两边都有活干。本地队列是 256 槽位的 lock-free 环形数组，大部分操作无锁。
意义：避免全局队列的锁竞争，同时保证负载均衡——忙的不堆积，闲的不空转。

---
Q7: goroutine 泄漏是什么？怎么排查和避免？
【理解】
goroutine 泄漏 = goroutine 启动后永远不退出，持续占用内存和调度资源。

常见泄漏场景：
  1. channel 阻塞：发送/接收永远等不到对端
  2. 锁未释放：goroutine 永远等锁
  3. 无限循环：没有退出条件
  4. select 无 default 且所有 case 永远不就绪
  5. 网络连接未设超时：goroutine 卡在 IO 上

排查方法：
  runtime.NumGoroutine()：监控协程数量，持续增长说明有泄漏
  pprof goroutine profile：go tool pprof http://localhost:6060/debug/pprof/goroutine
  goleak 库（uber 出品）：在测试中检测 goroutine 泄漏

避免方法：
  1. 用 context 控制生命周期，确保有退出路径
  2. channel 操作配合 select + timeout/done
  3. 网络操作设置超时（SetDeadline）
  4. 启动 goroutine 前想清楚"它什么时候退出"
【回答】
goroutine 泄漏就是协程启动后永远不退出，持续占用资源。
常见原因：channel 阻塞无对端、锁未释放、无限循环无退出条件、IO 无超时。
排查：runtime.NumGoroutine() 监控数量、pprof goroutine profile 看堆栈、goleak 库在测试中检测。
避免：用 context 控制生命周期、channel 配合 select+timeout、IO 设超时、启动前想清楚退出条件。

---
Q8: GOMAXPROCS 是什么？设置多少合适？
【理解】
GOMAXPROCS = P 的数量 = 同时执行 goroutine 的最大并行数。

默认值：runtime.NumCPU()（CPU 逻辑核数）。
设置方式：
  环境变量：GOMAXPROCS=4 ./myapp
  代码：runtime.GOMAXPROCS(4)

设置多少合适？
  CPU 密集型：默认值（= CPU 核数）最优，设更大反而增加切换开销
  IO 密集型：可以适当调大（如 2×CPU），因为很多 M 在等 IO，P 经常空闲
  容器环境：注意容器 CPU 限制，Go 1.5+ 默认读的是宿主机核数
    解决：用 uber/automaxprocs 库自动适配容器 CPU 配额

注意：
  GOMAXPROCS 限制的是 P 的数量，不是 M（线程）的数量。
  M 的数量可以超过 P（因为有些 M 在做系统调用被阻塞）。
  默认 M 上限是 10000（runtime/debug.SetMaxThreads 可调）。
【回答】
GOMAXPROCS 是 P 的数量，决定了同时并行执行 goroutine 的上限，默认等于 CPU 核数。
CPU 密集型用默认值最优；IO 密集型可以适当调大。
容器环境要注意：Go 默认读宿主机核数，可能远大于容器配额，建议用 automaxprocs 库自动适配。
注意它限制的是 P 不是 M，M 可以超过 P（阻塞在系统调用的 M 不占 P）。

---
Q9: goroutine 和 channel 怎么配合实现并发模式？
【理解】
常见并发模式：

1. Fan-out/Fan-in（扇出/扇入）：
   一个任务分发给多个 worker 并行处理，结果汇总到一个 channel。

2. Pipeline（流水线）：
   多个阶段串联，每个阶段是一组 goroutine，通过 channel 传递数据。

3. Worker Pool（工作池）：
   固定数量的 worker goroutine 从任务 channel 取任务执行。

4. Or-Done（任一完成）：
   多个 goroutine 竞争，第一个完成的结果胜出。

5. Rate Limiting（限流）：
   用 buffered channel 或 time.Ticker 控制并发速率。

选择原则：
  需要并行加速 -> Fan-out + Worker Pool
  数据流处理 -> Pipeline
  超时竞争 -> Or-Done + select
  控制并发度 -> buffered channel 做信号量
【回答】
常见模式：
Fan-out/Fan-in：任务分发给多个 worker 并行，结果汇总。
Pipeline：多阶段串联，channel 传递数据流。
Worker Pool：固定 N 个 worker 从任务 channel 取活干，控制并发度。
Or-Done：多个 goroutine 竞争，select 取第一个完成的结果。
核心思路：goroutine 负责并发执行，channel 负责通信和同步，context 负责生命周期控制。

---
Q10: goroutine 的生命周期？状态有哪些？
【理解】
goroutine 的状态机（runtime/runtime2.go 中的 g.atomicstatus）：

_Gidle(0)：刚分配，还没初始化
_Grunnable(1)：在运行队列中，等待被调度执行
_Grunning(2)：正在某个 M 上执行
_Gsyscall(3)：正在执行系统调用，M 被阻塞
_Gwaiting(4)：被阻塞（channel/锁/sleep/IO），不在运行队列中
_Gdead(6)：执行完毕或刚被分配还没使用
_Gcopystack(8)：栈正在被拷贝（扩容/缩容）
_Gpreempted(9)：被抢占，等待重新调度

生命周期：
  go func() 创建 -> _Grunnable（放入 P 的本地队列）
  被 M 调度执行 -> _Grunning
  遇到 channel 阻塞 -> _Gwaiting（挂到等待队列）
  被唤醒 -> _Grunnable（放回运行队列）
  函数执行完毕 -> _Gdead（G 结构体可被复用）

G 的复用：
  goroutine 执行完后不会立即释放 g 结构体，而是放入 P 的 gFree 列表。
  下次 go func() 时优先从 gFree 取，避免频繁分配。
【回答】
goroutine 主要状态：Runnable（在队列等调度）、Running（正在执行）、Waiting（阻塞中）、Syscall（系统调用中）、Dead（执行完毕）。
生命周期：go func() 创建进入 Runnable -> 被调度变 Running -> 阻塞变 Waiting -> 唤醒回 Runnable -> 执行完变 Dead。
Dead 状态的 G 不会立即释放，而是放入 gFree 列表复用，下次 go func() 优先取已有的 g 结构体，减少分配开销。

---
Q11: Hand Off 机制是什么？G 什么时候会放入全局队列？
【理解】
■ Hand Off（任务剥离）：
  当 M 执行 G 时发生系统调用阻塞（如文件 IO、网络 IO、CGO 调用）：
    Step1: M 进入阻塞，无法继续执行其他 G
    Step2: runtime 把 P 从阻塞的 M 上剥离
    Step3: P 寻找一个空闲的 M（或创建新 M）绑定，继续执行队列中的其他 G
    Step4: 原来的 M 系统调用返回后，尝试重新获取一个 P
           -> 有空闲 P：绑定继续执行
           -> 没有空闲 P：把当前 G 放入全局队列，M 进入休眠

  意义：不让一个系统调用阻塞整个 P 的工作，保证 CPU 利用率。

■ G 放入全局队列的时机：
  1. 本地队列已满（256 个）：新建 G 时本地放不下，把本地队列一半的 G 打包放入全局队列
  2. 系统调用返回后找不到空闲 P：该 G 进入全局队列等待
  3. 被抢占的 G：超时被抢占后放入全局队列
  4. 网络轮询就绪的 G：netpoll 返回的就绪 G 放入全局队列

■ 从全局队列取 G 的时机：
  本地队列空时，从全局队列取一批（min(全局队列长度/P数量+1, 本地队列容量/2)）
  每 61 次调度强制检查一次全局队列（保证公平性）
【回答】
Hand Off：M 阻塞在系统调用时，runtime 把 P 剥离给其他 M 继续工作，不浪费 CPU。M 恢复后如果找不到空闲 P，就把 G 放入全局队列自己休眠。
G 放入全局队列的时机：本地队列满（256 个）时把一半放入全局、系统调用返回找不到 P、被抢占的 G。
从全局取：本地空时取一批，每 61 次调度强制检查一次全局队列防止饥饿。

---
Q12: sysmon 是什么？调度器怎么保证公平性？
【理解】
■ sysmon（系统监控线程）：
  独立的后台线程，不绑定任何 P，不参与正常调度。
  职责：
    1. 抢占检测：发现 G 运行超过 10ms，发送抢占信号
    2. 网络轮询（netpoll）：检查网络 IO 是否就绪，唤醒等待的 G
    3. 强制 GC：超过 2 分钟没 GC 就强制触发
    4. 定时器检查：检查 timer 是否到期
    5. Hand Off 检测：发现 M 阻塞在系统调用超过一定时间，剥离 P

  sysmon 的运行频率：
    初始每 20μs 检查一次
    如果连续没有发现需要处理的事件，逐渐降低频率（最慢 10ms 一次）
    一旦发现事件，频率重新提高

■ 调度公平性保证：
  1. 每 61 次调度强制检查全局队列：
     P 每次从本地队列取 G 时有一个计数器（schedtick）
     schedtick % 61 == 0 时，优先从全局队列取
     防止全局队列中的 G 被"饿死"

  2. Work Stealing 保证负载均衡：
     闲的 P 偷忙的 P 的一半 G

  3. 抢占保证不会独占：
     单个 G 运行超过 10ms 被强制切换

  4. Hand Off 保证 P 不被浪费：
     M 阻塞时 P 转给其他 M

为什么是 61？
  61 是质数，和本地队列大小 256 互质，能更均匀地分散检查时机，避免和其他周期性行为产生共振。
【回答】
sysmon 是独立的监控线程，不绑定 P，职责包括：抢占检测（>10ms）、网络轮询、强制 GC（>2min）、定时器检查、Hand Off 检测。
公平性保证：每 61 次调度强制从全局队列取 G（防饥饿）、Work Stealing 负载均衡、抢占防独占、Hand Off 防 P 浪费。
为什么是 61：质数，和队列大小 256 互质，分散检查时机避免共振。

*/

// TestGoroutineBasic 基础创建和调度
func TestGoroutineBasic(t *testing.T) {
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Printf("goroutine %d 运行在线程上\n", id)
		}(i)
	}
	wg.Wait()
}

// TestGoroutineNum 查看 goroutine 数量
func TestGoroutineNum(t *testing.T) {
	fmt.Println("初始 goroutine 数:", runtime.NumGoroutine())

	for i := 0; i < 10; i++ {
		go func() {
			time.Sleep(time.Second)
		}()
	}

	fmt.Println("启动后 goroutine 数:", runtime.NumGoroutine())
}

// TestWorkerPool Worker Pool 模式
func TestWorkerPool(t *testing.T) {
	tasks := make(chan int, 20)
	var wg sync.WaitGroup

	// 启动 3 个 worker
	for w := 0; w < 3; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for task := range tasks {
				fmt.Printf("worker %d 处理任务 %d\n", id, task)
			}
		}(w)
	}

	// 发送 10 个任务
	for i := 0; i < 10; i++ {
		tasks <- i
	}
	close(tasks)
	wg.Wait()
}

// TestGOMAXPROCS 查看和设置 GOMAXPROCS
func TestGOMAXPROCS(t *testing.T) {
	fmt.Println("CPU 核数:", runtime.NumCPU())
	fmt.Println("当前 GOMAXPROCS:", runtime.GOMAXPROCS(0)) // 0 表示只查询不修改
}

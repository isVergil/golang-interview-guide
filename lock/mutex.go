package main

/*
Mutex 特点
1 完全互斥：Mutex 不区分读写，只要有一方持有锁，其他所有尝试 Lock 的协程都会被挂起（排队）。
2 不可重入：一旦加锁，在解锁前不能再次加锁，否则直接死锁。
3 自旋锁优化：为了避免系统调用导致的上下文切换（Context Switch）开销，Mutex 会先尝试自旋（Spinning），自旋若干次失败后，进入休眠（阻塞）状态，放入等待队列。
4 混合模式：自旋若干次失败后，协程会进入休眠（阻塞）状态，放入等待队列。
5 公平性：正常、饥饿模式来保证，底层维护的双向链表等待队列来处理
	正常模式（性能至上）：新来的协程（正在 CPU 上运行）和队列头部的协程（刚被唤醒）一起竞争。由于新协程已经在 CPU 上了，所以它大概率能赢。队列里的老协程只能继续等，并挪到队列头部。
	饥饿模式（公平至上）：如果队列里的协程等待时间超过 1ms，锁就会进入“饥饿模式”。锁的所有权会直接从解锁的协程移交给队列头的协程。新来的协程完全不能抢锁，只能乖乖去排队。

Mutex 底层结构
type Mutex struct {
	_ noCopy

	mu isync.Mutex
}
type Mutex struct {
    state int32  // 核心字段：记录锁的状态
    sema  uint32 // 信号量：用于控制协程的阻塞和唤醒（等待队列）
}

通过CAS：CompareAndSwapInt32实现具体的加锁和解锁操作。
1 极致的性能如果不用 CAS，修改 state 字段可能还得再加一把“底层的锁”，那就陷入了“为了实现锁而需要另一把锁”的死循环。CAS 直接利用 CPU 指令，不需要上下文切换，性能极高。
2 乐观锁思想，它假设大部分情况下没人跟我抢。我先尝试一下（CAS），成功了就赚了；失败了，我再通过 for 循环不断重试（自旋），直到成功。
3 在 Mutex 里，有多个位（Locked, Starving, WaitersCount）都挤在一个 int32 里。通过 CAS，我们可以确保：只有当我看到的锁状态和我预期的一模一样时，我才更新它。 只要中间有一丁点变动，我就重新读取再试。

state 是一个 int32 整数，它被分成了四个区域来存储不同的信息：
第 0 位 (1)mutexLocked加锁标志：0 为未加锁，1 为已加锁。
第 1 位 (2)mutexWoken唤醒标志：是否有协程正被唤醒去抢锁（防止多个唤醒浪费性能）。
第 2 位 (4)mutexStarving饥饿标志：是否处于饥饿模式。
剩余 29 位waitersCount等待者计数：当前有多少个协程正在队列里排队。

noCopy 作用：
它是一个语法警察， 因为 Go 语言默认是值拷贝的，如果我们定义的结构体里有锁（Mutex），一旦被拷贝，锁的状态就会乱掉。
_ noCopy 字段本身没有任何运行时的功能，它不占空间（空结构体）。它的唯一作用是让 go vet 这种静态检查工具识别出来。如果有人尝试拷贝这个结构体，go vet 就会报错提醒开发者：‘嘿，这个东西不能拷，请使用指针！’
这是一个典型的通过接口契约实现静态约束的设计方案。”

*/

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// 模拟源码中的 noCopy 结构
type noCopy struct{}

func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}

// Account 银行账户结构体
type Account struct {
	_       noCopy // 1. 防止拷贝：如果有人拷贝 Account，go vet 会报警
	mu      sync.Mutex
	balance int32 // 2. 使用 int32 方便展示原子操作
}

// Deposit 存款：展示 Mutex 的标准用法
func (a *Account) Deposit(amount int) {
	a.mu.Lock()         // 底层会调用 CAS 尝试将 state 从 0 改为 1
	defer a.mu.Unlock() // 解锁时会通过 sema 唤醒排队协程
	a.balance += int32(amount)
}

// CASUpdate 余额更新：展示 CompareAndSwap 的逻辑
func (a *Account) TryUpdateBalance(old, new int32) bool {
	// 3. 直接调用底层 CAS 指令
	// 语义：如果当前余额确实是 old，才更新为 new
	return atomic.CompareAndSwapInt32(&a.balance, old, new)
}

func main() {
	acc := &Account{balance: 100}

	// --- 演示锁拷贝风险 ---
	// 4. 这里的 accBad 是 acc 的一个值拷贝
	// 如果此时 acc 已经加锁了，accBad 里的 mu 状态也会是“已加锁”
	accBad := *acc
	fmt.Printf("执行了拷贝，请运行 'go vet' 查看警告\n")
	_ = accBad // 仅仅为了规避未使用变量错误

	// --- 演示 CAS 逻辑 ---
	success := acc.TryUpdateBalance(100, 150)
	fmt.Printf("第一次 CAS 更新结果: %v, 当前余额: %d\n", success, acc.balance)

	fail := acc.TryUpdateBalance(100, 200) // 此时余额已是 150，预期 100 会失败
	fmt.Printf("第二次 CAS 更新结果: %v (因为旧值对不上)\n", fail)

	// --- 演示并发安全性 ---
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			acc.Deposit(10)
		}()
	}
	wg.Wait()
	fmt.Printf("100次并发存款后余额: %d\n", acc.balance)
}

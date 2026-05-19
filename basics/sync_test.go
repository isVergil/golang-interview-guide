package basics

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

/*
Q1: sync.Mutex 和 sync.RWMutex 有什么区别？
Q2: Mutex 的底层实现？有哪些状态？
Q3: Mutex 的正常模式和饥饿模式是什么？
Q4: sync.WaitGroup 怎么用？底层原理？
Q5: sync.Once 怎么实现的？为什么不能用 CAS 替代？
Q6: sync.Map 适用什么场景？和加锁 map 有什么区别？
Q7: sync.Pool 是什么？怎么用？
Q8: sync.Cond 是什么？什么场景用？
Q9: atomic 包和 Mutex 怎么选？
Q10: 死锁的四个必要条件？Go 里怎么避免？
Q11: RWMutex 的底层原理？写优先怎么实现的？
Q12: noCopy 是什么？为什么 sync 包的类型不能复制？

---
Q1: sync.Mutex 和 sync.RWMutex 有什么区别？
【理解】
Mutex：互斥锁，同一时刻只有一个 goroutine 能持有锁。
RWMutex：读写锁，允许多个读并发，写互斥。

RWMutex 的规则：
  读锁（RLock）：多个 goroutine 可以同时持有读锁
  写锁（Lock）：独占，持有写锁时不能有任何读锁或写锁
  写优先：有写锁等待时，新的读锁请求会被阻塞（防止写饥饿）

选择原则：
  读多写少 -> RWMutex（读并发提升性能）
  读写均衡或写多 -> Mutex（RWMutex 的额外开销反而更慢）
  临界区很短 -> Mutex（RWMutex 的 overhead 不值得）

注意：
  Mutex 零值可用，不需要初始化
  不能复制已使用的 Mutex（go vet 会检查）
  不能重入（同一个 goroutine 连续 Lock 两次会死锁）
【回答】
Mutex 是互斥锁，同一时刻只有一个 goroutine 能持有。RWMutex 是读写锁，允许多读并发，写独占。
选择：读多写少用 RWMutex，读写均衡或临界区很短用 Mutex。
注意：Go 的锁不可重入，同一个 goroutine 连续 Lock 两次会死锁。Mutex 零值可用不需要初始化，但不能复制已使用的锁。

---
Q2: Mutex 的底层实现？有哪些状态？
【理解】
Mutex 结构体只有两个字段：
  type Mutex struct {
      state int32   // 状态位
      sema  uint32  // 信号量，用于阻塞/唤醒
  }

state 的 bit 含义：
  bit0: locked（是否被锁定）
  bit1: woken（是否有唤醒的 goroutine）
  bit2: starving（是否进入饥饿模式）
  bit3~31: waiter count（等待者数量）

加锁流程：
  Fast path: CAS 尝试将 locked 位从 0 设为 1，成功则获得锁
  Slow path: CAS 失败，进入自旋或排队：
    1. 先自旋几次（spin），期望持有者很快释放
    2. 自旋超过次数或条件不满足，通过信号量 sema 阻塞
    3. 被唤醒后重新竞争

解锁流程：
  Fast path: 将 locked 位清零，如果没有等待者直接返回
  Slow path: 有等待者，通过信号量唤醒一个

自旋条件（必须全部满足）：
  - 自旋次数 < 4
  - 多核 CPU（GOMAXPROCS > 1）
  - 当前 P 的本地队列为空（不会饿死其他 G）
【回答】
Mutex 只有 state（int32 状态位）和 sema（信号量）两个字段。
state 的低 3 位分别表示：是否锁定、是否有唤醒者、是否饥饿模式；高位存等待者数量。
加锁：先 CAS 快速尝试，失败则自旋几次（多核+队列空时），还拿不到就通过信号量阻塞。
解锁：清 locked 位，有等待者则通过信号量唤醒一个。

---
Q3: Mutex 的正常模式和饥饿模式是什么？
【理解】
■ 正常模式（Normal）：
  解锁时唤醒等待队列头部的 goroutine，但它不会直接获得锁。
  它需要和新到达的 goroutine 竞争（新来的正在 CPU 上运行，有优势）。
  优点：吞吐量高（正在运行的 G 直接拿锁，不需要切换）
  缺点：等待队列里的 G 可能一直抢不到（饥饿）

■ 饥饿模式（Starvation）：
  触发条件：某个等待者等待超过 1ms
  行为：解锁时直接把锁交给等待队列头部的 G（不竞争）
  新来的 G 不自旋，直接排到队尾
  优点：保证公平，避免尾部延迟过高
  缺点：吞吐量下降（每次都要切换 goroutine）

退出饥饿模式的条件：
  获得锁的 G 是队列中最后一个等待者，或等待时间 < 1ms

设计思路：正常模式追求性能，饥饿模式保证公平，两者动态切换。
【回答】
正常模式：唤醒的等待者和新来的 goroutine 竞争锁，新来的有优势（正在 CPU 上），吞吐量高但可能饥饿。
饥饿模式：等待超过 1ms 触发，解锁时直接把锁交给队首等待者，新来的直接排队尾，保证公平但吞吐下降。
两种模式动态切换：正常模式追求性能，饥饿模式兜底公平，避免尾部延迟过高。

---
Q4: sync.WaitGroup 怎么用？底层原理？
【理解】
WaitGroup 用于等待一组 goroutine 完成。

使用模式：
  var wg sync.WaitGroup
  wg.Add(n)           // 计数器 +n（启动前调用）
  go func() { defer wg.Done(); ... }()  // 完成时 -1
  wg.Wait()           // 阻塞直到计数器归零

底层结构：
  type WaitGroup struct {
      state atomic.Uint64  // 高 32 位: counter，低 32 位: waiter 数量
      sema  uint32         // 信号量
  }

  Add(n): atomic 加 counter
  Done(): Add(-1)，counter 归零时唤醒所有 waiter
  Wait(): waiter++，如果 counter > 0 则通过 sema 阻塞

注意事项：
  - Add 必须在 go 语句之前调用（否则 Wait 可能在 Add 之前返回）
  - 不能在 counter 为 0 后再 Add（会 panic）
  - WaitGroup 不能复制（和 Mutex 一样）
【回答】
WaitGroup 用于等待一组 goroutine 完成：Add 设置计数、Done 减一、Wait 阻塞到归零。
底层用 atomic 操作计数器 + 信号量阻塞/唤醒。counter 归零时唤醒所有调用 Wait 的 goroutine。
注意：Add 必须在 go 语句之前调用，不能复制 WaitGroup，counter 为 0 后不能再 Add。

---
Q5: sync.Once 怎么实现的？为什么不能用 CAS 替代？
【理解】
Once 保证函数只执行一次，常用于单例初始化。

底层实现：
  type Once struct {
      done atomic.Uint32  // 是否已完成
      m    Mutex          // 保护 f 的执行
  }
  func (o *Once) Do(f func()) {
      if o.done.Load() == 0 {  // fast path: 已完成直接返回
          o.doSlow(f)
      }
  }
  func (o *Once) doSlow(f func()) {
      o.m.Lock()
      defer o.m.Unlock()
      if o.done.Load() == 0 {  // double check
          defer o.done.Store(1)
          f()
      }
  }

为什么不能只用 CAS？
  如果用 CAS(0->1) 来标记"我来执行"：
    G1 CAS 成功，开始执行 f()（f 可能很慢）
    G2 看到 done=1，认为已完成，直接使用结果
    但此时 f() 还没执行完！G2 用到了未初始化的数据

  Once 的设计保证：Do 返回时 f 一定已经执行完毕。
  所以必须用 Mutex 让其他 goroutine 等待 f 完成，而不是看到标记就返回。
【回答】
Once 保证函数只执行一次。底层是 atomic 快速检查 + Mutex 保护执行 + double check。
不能只用 CAS 的原因：CAS 只能保证"只有一个人去执行"，但不能保证"其他人等执行完才返回"。如果 f() 很慢，其他 goroutine 看到标记就返回，会用到未初始化的数据。
Once 的语义是：Do 返回时 f 一定已经执行完毕，所以必须用锁让其他人等待。

---
Q6: sync.Map 适用什么场景？和加锁 map 有什么区别？
【理解】
sync.Map 是并发安全的 map，针对特定场景优化。

适用场景（官方文档明确说的）：
  1. key 稳定，读多写少（如缓存）
  2. 多个 goroutine 读写不同的 key（无竞争）

底层结构（双 map）：
  type Map struct {
      mu     Mutex
      read   atomic.Pointer[readOnly]  // 无锁读的 map（只读快照）
      dirty  map[any]*entry            // 需要加锁的 map（最新数据）
      misses int                       // read 未命中次数
  }

工作原理：
  读：先查 read（无锁 atomic），命中直接返回；未命中再加锁查 dirty
  写：加锁写入 dirty
  提升：misses 达到 dirty 长度时，dirty 提升为新的 read（摊销成本）

vs 加锁 map（map + RWMutex）：
  sync.Map：读多写少时性能好（读无锁），但写入和 miss 时有额外开销
  加锁 map：简单直接，读写均衡时性能更好，内存开销更小

大多数场景用 map + RWMutex 就够了，sync.Map 只在特定场景有优势。
【回答】
sync.Map 适用于读多写少且 key 稳定的场景（如缓存），底层是双 map：read（无锁快速读）+ dirty（加锁写入）。
读先查 read 无锁命中直接返回，未命中再加锁查 dirty。miss 次数够多时 dirty 提升为新 read。
vs 加锁 map：sync.Map 读多写少时快（读无锁），但写入有额外开销且内存占用更大。大多数场景 map+RWMutex 更简单更好。

---
Q7: sync.Pool 是什么？怎么用？
【理解】
sync.Pool 是临时对象缓存池，用于复用对象减少 GC 压力。（详见 gc_test.go Q10）

核心 API：
  pool := &sync.Pool{
      New: func() any { return make([]byte, 1024) },
  }
  obj := pool.Get()   // 取对象（池空则调用 New）
  pool.Put(obj)       // 还回去

底层结构（每个 P 有私有缓存）：
  每个 P 有 private 字段（无锁快速存取一个对象）
  每个 P 有 shared 链表（其他 P 可以偷）
  Get: 先取 private -> 再取本地 shared -> 再偷其他 P 的 shared -> 调用 New
  Put: 先存 private -> 存本地 shared

GC 时的行为：
  Go 1.13+: victim cache 机制，对象可存活两轮 GC
  不保证持久性，不能做连接池

典型使用：bytes.Buffer、编解码器、临时 []byte
【回答】
sync.Pool 是临时对象缓存池，Get 取对象、Put 还回去，池空时调用 New 创建。
底层每个 P 有私有缓存（无锁快速存取），还有 shared 链表支持跨 P 偷取。
GC 时可能清理池中对象（Go 1.13+ 有 victim cache 可存活两轮），所以不能做连接池。
适合缓存 bytes.Buffer 等临时对象，减少分配和 GC 压力。

---
Q8: sync.Cond 是什么？什么场景用？
【理解】
sync.Cond 是条件变量，用于 goroutine 之间的条件等待和通知。

核心 API：
  cond := sync.NewCond(&sync.Mutex{})
  cond.Wait()      // 释放锁 + 阻塞等待 + 被唤醒后重新获取锁
  cond.Signal()    // 唤醒一个等待者
  cond.Broadcast() // 唤醒所有等待者

使用模式（必须在锁内使用 Wait）：
  cond.L.Lock()
  for !condition {   // 必须用 for 循环检查条件（防止虚假唤醒）
      cond.Wait()
  }
  // 条件满足，执行逻辑
  cond.L.Unlock()

适用场景：
  - 多个 goroutine 等待同一个条件成立
  - 生产者-消费者（通知消费者有新数据）
  - 等待某个状态变化（如连接池等待可用连接）

vs channel：
  channel 更常用，语义更清晰
  Cond 适合"多个等待者等同一个条件 + 需要 Broadcast 唤醒所有人"的场景
【回答】
sync.Cond 是条件变量：Wait 释放锁并阻塞，Signal 唤醒一个，Broadcast 唤醒所有。
使用时必须在锁内用 for 循环检查条件（防虚假唤醒）。
适用场景：多个 goroutine 等待同一条件成立，需要 Broadcast 唤醒所有人。
大多数场景 channel 更好用，Cond 只在需要"广播唤醒 + 条件检查"时才有优势。

---
Q9: atomic 包和 Mutex 怎么选？
【理解】
atomic：CPU 级别的原子操作，无锁，性能极高（纳秒级）。
Mutex：操作系统级别的锁，有上下文切换开销（几十纳秒~微秒级）。

atomic 适用场景：
  - 单个变量的简单操作（计数器、标志位、指针交换）
  - 读多写少的配置热更新（atomic.Value）
  - 性能极度敏感的热路径

Mutex 适用场景：
  - 保护多个变量的复合操作（需要原子性的一组操作）
  - 临界区有复杂逻辑（不是简单的读写）
  - 需要条件等待（配合 Cond）

注意：
  atomic 只保证单个操作的原子性，不保证多个操作的原子性
  如果需要"读-改-写"多个变量，必须用锁
  atomic.Value 的 Store 和 Load 是原子的，但两次 Load 之间值可能变了
【回答】
atomic 是 CPU 级原子操作，无锁极快，适合单变量简单操作（计数器、标志位、指针交换）。
Mutex 适合保护多变量复合操作或复杂临界区逻辑。
选择原则：单个变量的简单读写用 atomic；多个变量需要一起保证一致性用 Mutex。
atomic 只保证单个操作原子，不保证多个操作组合的原子性。

---
Q10: 死锁的四个必要条件？Go 里怎么避免？
【理解】
死锁四个必要条件（同时满足才会死锁）：
  1. 互斥：资源同一时刻只能被一个 goroutine 持有
  2. 持有并等待：持有一个资源的同时等待另一个
  3. 不可剥夺：已持有的资源不能被强制释放
  4. 循环等待：A 等 B，B 等 A，形成环

Go 里常见死锁场景：
  - 同一个 goroutine 对 Mutex 连续 Lock 两次（不可重入）
  - 两个 goroutine 互相等对方的锁（AB-BA 问题）
  - channel 在同一个 goroutine 里又发又收（无缓冲）
  - WaitGroup Add 和 Done 不匹配

避免方法：
  - 固定加锁顺序（破坏循环等待）
  - 用 defer Unlock 确保释放（破坏持有并等待）
  - 用 TryLock（Go 1.18+）避免无限等待
  - 用 context 超时控制
  - go vet / deadlock 检测工具
【回答】
四个必要条件：互斥、持有并等待、不可剥夺、循环等待。
Go 常见死锁：Mutex 连续 Lock 两次、两个 goroutine 互相等对方的锁、无缓冲 channel 同一 goroutine 收发。
避免方法：固定加锁顺序（破坏循环等待）、defer Unlock（确保释放）、TryLock 避免无限等待、context 超时控制。

---
Q11: RWMutex 的底层原理？写优先怎么实现的？
【理解】
RWMutex 底层结构（sync/rwmutex.go）：
  type RWMutex struct {
      w           Mutex   // 写锁之间的互斥
      writerSem   uint32  // 写等待信号量：读完后唤醒写
      readerSem   uint32  // 读等待信号量：写完后唤醒读
      readerCount int32   // 当前读协程数（负数表示有写在排队）
      readerWait  int32   // 写锁等待的剩余读协程数
  }

写优先的实现（readerCount 减 2^30 技巧）：
  1. 正常读：RLock 时 readerCount++，正数表示有多少读协程
  2. 写锁进入（Lock）：
     - 先获取内部 w Mutex（写写互斥）
     - 把 readerCount 减去 rwmutexMaxReaders（2^30），变成很大的负数
     - 此时 readerCount 为负，新来的 RLock 看到负数就知道有写在排队，乖乖阻塞在 readerSem
     - readerWait 记录写锁之前已有的读协程数，等它们全部 RUnlock 后唤醒写
  3. 写锁释放（Unlock）：
     - readerCount 加回 2^30，恢复正数
     - 唤醒所有阻塞在 readerSem 的读协程

为什么这样设计？
  用一个 int32 的正负号就区分了"有没有写在排队"，不需要额外的标志位。
  一次原子操作同时完成"标记写排队"和"记录当前读数"两件事。
【回答】
RWMutex 底层有 readerCount（读协程数）和 readerWait（写等待的读数）两个关键字段。
写优先实现：写锁进入时把 readerCount 减去 2^30 变成负数，新来的读锁看到负数就阻塞等待。写锁释放时加回 2^30 恢复正数，唤醒所有等待的读协程。
一次原子操作同时完成"标记写排队"和"记录当前读数"，设计非常精巧。

---
Q12: noCopy 是什么？为什么 sync 包的类型不能复制？
【理解】
noCopy 是 Go 标准库中的一个空结构体，用于防止值拷贝：
  type noCopy struct{}
  func (*noCopy) Lock()   {}
  func (*noCopy) Unlock() {}

它实现了 sync.Locker 接口，但 Lock/Unlock 什么都不做。
唯一作用：让 go vet 静态检查工具识别——如果有人拷贝了包含 noCopy 的结构体，go vet 会报警。

为什么 sync 包类型不能复制？
  Mutex：state 字段记录了锁状态（是否加锁、等待者数量）。
    拷贝后副本的 state 和原件一样，如果原件已加锁，副本也是"已加锁"状态，
    但没有人会去 Unlock 副本，导致死锁。
  WaitGroup：state 记录了计数器和等待者数量。
    拷贝后副本的计数器和原件一样，Done 改的是原件，Wait 等的是副本，永远等不到。
  Cond：内部持有 Mutex 指针和等待队列。
    拷贝后队列状态混乱。

CAS 在 Mutex 中的作用：
  Mutex 的 Lock 通过 CAS（CompareAndSwap）修改 state 字段。
  如果不用 CAS，修改 state 还得再加一把"底层锁"，陷入"为了实现锁需要另一把锁"的死循环。
  CAS 利用 CPU 原子指令，一条指令完成"比较+交换"，无需上下文切换。
【回答】
noCopy 是空结构体，实现了 Locker 接口但什么都不做，唯一作用是让 go vet 检测到值拷贝时报警。
sync 类型不能复制的原因：Mutex 拷贝后副本保留了"已加锁"状态但没人 Unlock 导致死锁；WaitGroup 拷贝后计数器和等待者分离导致永久阻塞。
Mutex 底层用 CAS 原子指令修改 state，避免了"为了实现锁需要另一把锁"的递归问题。

*/

// TestMutex 互斥锁基础用法
func TestMutex(t *testing.T) {
	var mu sync.Mutex
	count := 0
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			count++
			mu.Unlock()
		}()
	}
	wg.Wait()
	fmt.Println("count:", count) // 1000
}

// TestRWMutex 读写锁：多读并发
func TestRWMutex(t *testing.T) {
	var rw sync.RWMutex
	data := "initial"
	var wg sync.WaitGroup

	// 多个读者并发
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			rw.RLock()
			fmt.Printf("reader %d: %s\n", id, data)
			rw.RUnlock()
		}(i)
	}

	// 一个写者
	wg.Add(1)
	go func() {
		defer wg.Done()
		rw.Lock()
		data = "updated"
		rw.Unlock()
	}()

	wg.Wait()
}

// TestOnce 只执行一次
func TestOnce(t *testing.T) {
	var once sync.Once
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			once.Do(func() {
				fmt.Println("只执行一次，by goroutine", id)
			})
		}(i)
	}
	wg.Wait()
}

// TestAtomic atomic vs 非原子操作
func TestAtomic(t *testing.T) {
	var count int64
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			atomic.AddInt64(&count, 1)
		}()
	}
	wg.Wait()
	fmt.Println("atomic count:", count) // 1000
}

// TestCond 条件变量演示
func TestCond(t *testing.T) {
	cond := sync.NewCond(&sync.Mutex{})
	ready := false

	// 等待者
	go func() {
		cond.L.Lock()
		for !ready {
			cond.Wait()
		}
		fmt.Println("条件满足，开始工作")
		cond.L.Unlock()
	}()

	// 通知者
	time.Sleep(100 * time.Millisecond)
	cond.L.Lock()
	ready = true
	cond.Signal()
	cond.L.Unlock()

	time.Sleep(50 * time.Millisecond)
}

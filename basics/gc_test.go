package basics

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
)

/*
Q1: Go 的 GC 用的什么算法？为什么选这个？
Q2: 三色标记法是什么？怎么工作的？
Q3: 为什么需要写屏障？没有会怎样？
Q4: Go 用的什么写屏障？插入屏障和删除屏障的区别？
Q5: STW（Stop The World）是什么？Go 怎么减少 STW 时间的？
Q6: GC 的完整流程（各阶段）？
Q7: GC 的触发条件有哪些？
Q8: 如何减少 GC 压力？有哪些调优手段？
Q9: 逃逸分析和 GC 有什么关系？
Q10: sync.Pool 和 GC 的关系？

---
Q1: Go 的 GC 用的什么算法？为什么选这个？
【理解】
Go 使用的是 并发三色标记-清除（Concurrent Tri-color Mark-Sweep）算法。

对比其他 GC 算法：
  引用计数（Python/PHP）：实时回收，但循环引用无法处理，且每次赋值都要更新计数（性能差）
  分代收集（Java）：按对象年龄分代，年轻代频繁收集。需要写屏障追踪跨代引用，实现复杂
  标记-清除（Go）：从根出发标记所有可达对象，清除不可达对象。简单直接

Go 为什么不用分代 GC？
  1. Go 有值类型和栈分配，大量短命对象直接在栈上分配，函数返回即回收，不进堆
  2. 逃逸分析把能留在栈上的对象留在栈上，堆上对象本身就少
  3. 分代 GC 的写屏障开销更大（要追踪所有跨代指针），Go 追求低延迟

Go GC 的设计目标：低延迟优先（STW < 1ms），而不是高吞吐。
代价：吞吐量不如 Java G1/ZGC，但对 Web 服务（延迟敏感）更友好。
【回答】
Go 用的是并发三色标记-清除算法。选这个是因为 Go 的设计目标是低延迟（STW < 1ms），而不是高吞吐。
Go 不用分代 GC 的原因：Go 有值类型和逃逸分析，大量短命对象直接在栈上分配不进堆，分代的收益不大；而且分代 GC 需要更重的写屏障追踪跨代引用，开销更大。
代价是 GC 吞吐量不如 Java，但对 Web 服务这种延迟敏感场景更友好。

---
Q2: 三色标记法是什么？怎么工作的？
【理解】
三色标记把堆上所有对象分成三种颜色：
  白色：未被访问，GC 结束后如果还是白色就回收
  灰色：已被发现但子对象还没扫描完（待处理队列）
  黑色：已扫描完毕，自身和所有子对象都已处理

标记过程（BFS/DFS 遍历对象图）：
  Step1: 初始状态——所有对象为白色
  Step2: 把根对象（栈变量、全局变量、寄存器）标灰，放入灰色队列
  Step3: 从灰色队列取出一个对象，扫描它的所有指针字段：
         - 指向的白色对象标灰（加入队列）
         - 自己标黑（处理完毕）
  Step4: 重复 Step3 直到灰色队列为空
  Step5: 此时剩下的白色对象就是垃圾，清除回收

三色不变式（保证正确性的核心约束）：
  强三色不变式：黑色对象不能直接引用白色对象
  弱三色不变式：黑色对象可以引用白色对象，但该白色对象必须有灰色对象保护（能通过灰色到达）

如果不维护三色不变式，并发标记时可能出现"漏标"：
  黑色 A 新增指向白色 C 的引用，同时灰色 B 删除对 C 的引用
  -> C 不在任何灰色对象的待扫描列表里，被误回收
  -> 程序访问悬挂指针，崩溃
【回答】
三色标记把对象分为白（未访问）、灰（已发现待扫描）、黑（已扫描完毕）三种状态。
从根对象开始，把根标灰放入队列。循环取出灰色对象，扫描它的指针字段，把引用的白色对象标灰，自身标黑。直到灰色队列为空，剩下的白色对象就是垃圾。
核心约束是三色不变式：黑色对象不能直接引用白色对象（否则白色对象会被误回收）。并发标记时通过写屏障来维护这个不变式。

---
Q3: 为什么需要写屏障？没有会怎样？
【理解】
写屏障解决的问题：GC 标记阶段和用户程序（mutator）并发执行时的正确性问题。

没有写屏障会出什么问题（漏标场景）：
  初始状态：灰色 B 引用白色 C
  GC 正在扫描 B 的子对象...

  用户程序并发执行：
    Step1: 黑色 A 新增引用 -> C（A.field = C）
    Step2: 灰色 B 删除引用 -> C（B.field = nil）

  此时 C 的状态：
    - A 是黑色已扫描完，不会再扫描 A 的子对象，看不到 C
    - B 删除了对 C 的引用，扫描 B 时也看不到 C
    - C 没有任何灰色对象能保护它
    - 结果：C 被误判为白色垃圾被回收！
    - 程序通过 A.field 访问 C -> 悬挂指针 -> 崩溃

漏标的两个必要条件（需要同时满足才会出问题）：
  条件 1: 黑色对象新增了对白色对象的引用（赋值器插入）
  条件 2: 所有灰色对象到该白色对象的路径都被破坏（赋值器删除）

写屏障的作用：在指针赋值时插入一小段代码，破坏上述任一条件，保证不会漏标。
【回答】
写屏障解决的是并发标记时的"漏标"问题。
GC 标记和用户程序同时运行时，可能出现：黑色对象新增了对白色对象的引用，同时灰色对象删除了对该白色对象的引用。此时白色对象没有任何灰色保护，会被误回收，导致悬挂指针崩溃。
写屏障就是在指针赋值时插入的一小段代码，保证这种情况不会发生。没有写屏障，并发 GC 就不可能正确工作。

---
Q4: Go 用的什么写屏障？插入屏障和删除屏障的区别？
【理解】
■ 插入写屏障（Dijkstra）：
  触发时机：当黑色对象 A 新增指向白色对象 C 的引用时（A.field = C）
  动作：把 C 标灰（"新引用的对象不会漏掉"）
  保证：破坏漏标条件 1，黑色不会引用白色
  缺点：Go 的栈上不开启写屏障（性能原因），所以标记结束后需要 STW 重新扫描所有栈

■ 删除写屏障（Yuasa）：
  触发时机：当灰色/白色对象 B 删除对白色对象 C 的引用时（B.field = nil，旧值是 C）
  动作：把被删除的旧引用 C 标灰（"快照保护，删之前先保留"）
  保证：破坏漏标条件 2，保证灰色到白色的路径不断
  缺点：保守，可能让本该回收的对象多活一轮

■ 混合写屏障（Go 1.8+，当前使用）：
  结合了插入和删除屏障的优点：
  伪代码：
    writePointer(slot, ptr):
        shade(*slot)  // 删除屏障：旧值标灰
        shade(ptr)    // 插入屏障：新值标灰
        *slot = ptr

  优势：
    1. 栈上对象不需要写屏障（GC 开始时把所有栈对象标黑）
    2. 不需要标记结束后 STW 重新扫描栈
    3. STW 时间进一步缩短到 < 1ms

  GC 开始时的特殊处理：
    把所有 goroutine 栈上的对象直接标黑（栈是根，认为栈上对象都是可达的）
    堆上新分配的对象也标黑（分配即可达）
    这样栈上指针变化就不需要写屏障了

【回答】
Go 1.8+ 使用混合写屏障，结合了插入屏障和删除屏障的优点。
插入屏障：新引用的对象标灰，防止黑色直接引用白色。缺点是栈上不开屏障，标记结束要 STW 重扫栈。
删除屏障：被删除的旧引用标灰，快照保护。缺点是保守，垃圾可能多活一轮。
混合屏障：GC 开始时把所有栈对象直接标黑，堆上赋值时同时对新值和旧值标灰。这样栈不需要写屏障，不需要重扫栈，STW < 1ms。

---
Q5: STW（Stop The World）是什么？Go 怎么减少 STW 时间的？
【理解】
STW = 暂停所有用户 goroutine，只有 GC 线程在运行。

Go GC 中需要 STW 的阶段（Go 1.8+）：
  STW 1（Mark Setup）：开启写屏障，扫描栈根对象。通常 < 100μs
  STW 2（Mark Termination）：关闭写屏障，做一些收尾工作。通常 < 100μs

整个标记阶段（最耗时）是并发的，不需要 STW。

Go 减少 STW 的演进：
  Go 1.0：完全 STW，标记+清除全程暂停，停顿可达秒级
  Go 1.3：标记阶段 STW，清除阶段并发
  Go 1.5：并发标记+并发清除，STW 只在开头和结尾（三色标记引入）
  Go 1.8：混合写屏障，去掉标记结束时的栈重扫描，STW < 1ms
  Go 1.12+：进一步优化，STW 通常 < 100μs

如何做到暂停所有 goroutine？
  设置抢占标记，等每个 G 到达安全点（safe point）后停下来。
  Go 1.14+ 有异步抢占（信号抢占），即使没有函数调用也能停下。
【回答】
STW 就是暂停所有用户 goroutine，只让 GC 运行。Go 目前只在 GC 的开头和结尾有两小段 STW（开启/关闭写屏障），各约 100μs，整个标记和清除阶段都是并发的。
演进过程：Go 1.0 全程 STW（秒级）→ Go 1.5 并发标记（毫秒级）→ Go 1.8 混合写屏障去掉栈重扫（< 1ms）→ 现在通常 < 100μs。
核心思路：把能并发做的工作都移到 STW 外面，只留最小的不可并发部分。

---
Q6: GC 的完整流程（各阶段）？
【理解】
Go GC 分为 4 个阶段：

Phase 1: Mark Setup（标记准备）—— STW
  - 开启写屏障
  - 所有 P（处理器）停下来，扫描各自的栈，把栈上对象标灰/标黑
  - 很快完成，< 100μs

Phase 2: Concurrent Mark（并发标记）—— 并发
  - 用户程序和 GC 同时运行
  - GC worker 从灰色队列取对象，扫描指针字段，标灰子对象，自身标黑
  - 用户程序的指针写入通过写屏障保证正确性
  - 占用约 25% CPU（默认 GOMAXPROCS/4 个 GC worker）
  - 如果 GC 跟不上分配速度，会让分配内存的 goroutine 协助标记（Mark Assist）

Phase 3: Mark Termination（标记终止）—— STW
  - 确保没有遗漏的灰色对象
  - 关闭写屏障
  - 计算下次 GC 触发阈值
  - 很快完成，< 100μs

Phase 4: Concurrent Sweep（并发清除）—— 并发
  - 用户程序和清除同时运行
  - 遍历所有 span，把白色对象对应的内存归还给 mcentral/mheap
  - 清除是惰性的：不是一次性全扫，而是分配新对象时顺带清理对应 span
  - 几乎不影响用户程序延迟

Mark Assist（标记协助）机制：
  如果 GC 标记速度跟不上用户分配速度，分配内存的 goroutine 会被"征召"帮忙标记。
  目的是防止堆无限增长。代价是分配变慢（背压机制）。
【回答】
Go GC 四个阶段：
标记准备（STW）：开启写屏障，扫描栈根对象，约 100μs。
并发标记（并发）：GC worker 和用户程序同时运行，BFS 遍历对象图标记可达对象。占 25% CPU，跟不上分配速度时会征召用户 goroutine 协助标记。
标记终止（STW）：关闭写屏障，计算下次触发阈值，约 100μs。
并发清除（并发）：惰性清除白色对象内存，几乎不影响延迟。
整个过程只有两小段 STW，标记和清除都是并发的。

---
Q7: GC 的触发条件有哪些？
【理解】
三种触发方式：

1. 堆内存增长触发（最常见）：
   当堆内存达到上次 GC 后存活量的一定比例时触发。
   由 GOGC 环境变量控制，默认 100，意思是堆增长 100% 时触发。
   公式：触发阈值 = 上次 GC 后存活堆大小 × (1 + GOGC/100)
   例：上次 GC 后存活 4MB，GOGC=100 → 堆到 8MB 时触发下一次 GC
   GOGC=200 → 堆到 12MB 时触发（更少 GC，更多内存）
   GOGC=50  → 堆到 6MB 时触发（更多 GC，更少内存）
   GOGC=off → 关闭 GC

2. 定时触发：
   如果超过 2 分钟没有触发过 GC，runtime 会强制触发一次。
   防止长时间不分配内存导致垃圾积累。
   源码：runtime/proc.go -> forcegchelper()

3. 手动触发：
   调用 runtime.GC()，会触发一次完整的 GC 并等待完成。
   一般只在测试或特殊场景使用，生产环境不建议手动调用。

Go 1.19+ 新增 GOMEMLIMIT：
  设置内存上限（soft limit），当内存接近上限时 GC 会更激进。
  可以配合 GOGC=off 使用：关闭基于增长比例的触发，纯靠内存上限控制。
  适用于容器环境，防止 OOM Kill。
【回答】
三种触发条件：
堆增长触发（最常见）：堆大小增长到上次 GC 后存活量的 GOGC% 时触发，默认 GOGC=100 即翻倍触发。
定时触发：超过 2 分钟没 GC 就强制触发一次。
手动触发：runtime.GC()，测试用，生产不建议。
Go 1.19+ 新增 GOMEMLIMIT，设置内存软上限，接近上限时 GC 更激进，适合容器环境防 OOM。

---
Q8: 如何减少 GC 压力？有哪些调优手段？
【理解】
核心原则：减少堆分配 = 减少 GC 工作量。

代码层面：
  1. 复用对象：sync.Pool 缓存临时对象，避免频繁分配/回收
  2. 预分配：make([]T, 0, n) 预分配容量，避免 append 多次扩容
  3. 避免逃逸：小对象尽量留在栈上，减少堆分配（见逃逸分析）
  4. 减少指针：用值类型代替指针类型，GC 扫描时不需要追踪
  5. 字符串拼接：用 strings.Builder 代替 + 拼接（减少中间字符串分配）
  6. 用数组代替 map：map 有大量内部指针，GC 扫描开销大

环境变量/参数：
  GOGC=200：降低 GC 频率（用更多内存换更少 GC）
  GOMEMLIMIT=512MiB：设置内存上限，配合 GOGC 使用
  debug.SetGCPercent()：运行时动态调整 GOGC

诊断工具：
  go tool pprof -alloc_space：查看哪里分配最多
  go build -gcflags="-m"：查看逃逸分析结果
  GODEBUG=gctrace=1：打印每次 GC 的详细信息
  runtime.ReadMemStats()：读取内存统计信息

注意：
  过度优化可能损害可读性。先用 pprof 确认 GC 确实是瓶颈再优化。
  大多数 Web 服务 GC 不是瓶颈。
【回答】
核心原则：减少堆分配就是减少 GC 压力。
代码层面：sync.Pool 复用对象、预分配 slice 容量、避免逃逸让对象留在栈上、减少指针（值类型替代指针类型，GC 不用追踪）、strings.Builder 替代 + 拼接。
参数调优：GOGC 调大降低 GC 频率（用内存换 CPU）、GOMEMLIMIT 设内存上限防 OOM。
诊断：pprof 看分配热点、-gcflags="-m" 看逃逸分析、GODEBUG=gctrace=1 看 GC 日志。
先确认 GC 确实是瓶颈再优化，大多数服务 GC 不是瓶颈。

---
Q9: 逃逸分析和 GC 有什么关系？
【理解】
逃逸分析（Escape Analysis）是编译器在编译期决定变量分配在栈还是堆的分析过程。

栈分配 vs 堆分配：
  栈：函数返回自动回收，零成本，不需要 GC 介入
  堆：需要 GC 追踪和回收，有标记/清除开销

逃逸的常见场景（变量必须分配到堆上）：
  1. 函数返回局部变量的指针（外部还要用，不能随函数栈帧销毁）
  2. 变量被闭包捕获（闭包生命周期可能超过当前函数）
  3. 变量赋值给接口类型（接口内部用指针存）
  4. slice/map 太大，超过栈帧大小限制
  5. 发送指针到 channel（跨协程，编译器无法确定生命周期）

不逃逸的场景（留在栈上）：
  1. 局部变量只在函数内使用
  2. 值拷贝传参（即使结构体大一点，也比逃逸到堆好）
  3. 内联函数的局部变量可能不逃逸

查看逃逸分析结果：
  go build -gcflags="-m" ./...
  输出如："moved to heap: x"、"x escapes to heap"

关系总结：
  逃逸到堆的对象越多 → GC 要扫描/回收的对象越多 → GC 压力越大
  减少逃逸 = 减少堆对象 = 减轻 GC 负担
【回答】
逃逸分析决定变量分配在栈还是堆。栈上的对象函数返回自动回收，零成本；堆上的对象需要 GC 追踪和回收。
逃逸到堆的对象越多，GC 工作量越大。所以减少逃逸就是减轻 GC 压力。
常见逃逸场景：返回局部变量指针、闭包捕获、赋值给接口、发送到 channel。
用 go build -gcflags="-m" 可以看哪些变量逃逸了，针对性优化。

---
Q10: sync.Pool 和 GC 的关系？
【理解】
sync.Pool 是一个临时对象缓存池，用来复用对象，减少堆分配和 GC 压力。

核心特性：
  - Pool 里的对象随时可能被 GC 回收（没有持久性保证）
  - 每次 GC 时，Pool 会清空（Go 1.13 之前），或淘汰旧对象（Go 1.13+）
  - 适合缓存临时对象（如 buffer），不适合做连接池等需要持久性的场景

底层机制（Go 1.13+，victim cache）：
  GC 时不是直接清空，而是：
    Step1: 把当前 pool 的对象移到 victim cache（受害者缓存）
    Step2: 清空上一轮的 victim cache
    Step3: Get 时先找 local pool，再找 victim cache
  效果：对象可以存活两轮 GC，减少冷启动时的分配风暴

为什么不能做连接池？
  GC 可能随时回收 Pool 里的连接，导致连接丢失。
  连接池要用自己管理生命周期的数据结构（如带引用计数的队列）。

典型使用场景：
  - bytes.Buffer / []byte 缓冲区
  - 编解码器（json.Encoder 等）
  - fmt 包内部就大量使用 sync.Pool

性能收益：
  减少内存分配次数 → 减少 GC 扫描对象数 → 降低 GC 频率和耗时
  同时减少了内存分配器的锁竞争
【回答】
sync.Pool 是临时对象缓存池，复用对象来减少堆分配和 GC 压力。
核心特点：Pool 里的对象没有持久性保证，GC 时可能被回收。Go 1.13+ 引入 victim cache 机制，对象能存活两轮 GC，减少了冷启动分配风暴。
适用场景：缓存 bytes.Buffer 等临时对象。不适合做连接池（GC 回收会导致连接丢失）。
收益：减少分配次数 → GC 扫描对象变少 → GC 频率和耗时都降低。

*/

// TestGCStats 查看 GC 统计信息
func TestGCStats(t *testing.T) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Printf("堆分配对象数: %d\n", m.HeapObjects)
	fmt.Printf("堆使用内存: %d KB\n", m.HeapAlloc/1024)
	fmt.Printf("GC 次数: %d\n", m.NumGC)
	fmt.Printf("上次 GC 暂停时间: %d μs\n", m.PauseNs[(m.NumGC+255)%256]/1000)
}

// TestGCTrigger 演示手动触发 GC 和观察效果
func TestGCTrigger(t *testing.T) {
	// 分配一些对象
	var s []*[1024]byte
	for i := 0; i < 1000; i++ {
		s = append(s, new([1024]byte))
	}

	var m1 runtime.MemStats
	runtime.ReadMemStats(&m1)
	fmt.Printf("GC 前 - 堆使用: %d KB, GC 次数: %d\n", m1.HeapAlloc/1024, m1.NumGC)

	// 释放引用，让对象变成垃圾
	s = nil
	runtime.GC() // 手动触发 GC

	var m2 runtime.MemStats
	runtime.ReadMemStats(&m2)
	fmt.Printf("GC 后 - 堆使用: %d KB, GC 次数: %d\n", m2.HeapAlloc/1024, m2.NumGC)
	fmt.Printf("回收了约 %d KB\n", (m1.HeapAlloc-m2.HeapAlloc)/1024)
}

// TestEscapeAnalysis 演示逃逸 vs 不逃逸
// 运行 go build -gcflags="-m" 可以看到逃逸分析结果
func TestEscapeAnalysis(t *testing.T) {
	// 不逃逸：局部变量只在函数内使用
	x := 42
	fmt.Println("栈上变量:", x)

	// 逃逸：返回指针，外部要用
	p := newInt(100)
	fmt.Println("堆上变量:", *p)
}

// newInt 返回局部变量指针，变量逃逸到堆
func newInt(n int) *int {
	v := n // v 会逃逸到堆上（因为返回了指针）
	return &v
}

// TestSyncPoolEffect 演示 sync.Pool 减少分配
func TestSyncPoolEffect(t *testing.T) {
	// 不使用 Pool：每次都分配新对象
	var m1 runtime.MemStats
	runtime.ReadMemStats(&m1)

	for i := 0; i < 10000; i++ {
		b := make([]byte, 1024)
		_ = b
	}

	var m2 runtime.MemStats
	runtime.ReadMemStats(&m2)
	fmt.Printf("不用 Pool - 分配次数增加: %d\n", m2.Mallocs-m1.Mallocs)

	// 使用 Pool：复用对象
	pool := &poolExample{}
	runtime.ReadMemStats(&m1)

	for i := 0; i < 10000; i++ {
		b := pool.get()
		pool.put(b)
	}

	runtime.ReadMemStats(&m2)
	fmt.Printf("用 Pool  - 分配次数增加: %d\n", m2.Mallocs-m1.Mallocs)
}

type poolExample struct {
	pool sync.Pool
}

func (p *poolExample) get() []byte {
	if v := p.pool.Get(); v != nil {
		return v.([]byte)
	}
	return make([]byte, 1024)
}

func (p *poolExample) put(b []byte) {
	p.pool.Put(b)
}

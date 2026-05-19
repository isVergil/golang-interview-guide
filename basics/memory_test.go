package basics

import (
	"fmt"
	"testing"
	"unsafe"
)

/*
Q1: Go 的内存分配器是怎么设计的？（TCMalloc）
Q2: mcache、mcentral、mheap 分别是什么？
Q3: 什么是 span（mspan）？size class 是什么？
Q4: 小对象、大对象的分配路径分别是什么？
Q5: 栈内存和堆内存有什么区别？
Q6: 逃逸分析的规则？哪些情况会逃逸到堆？
Q7: 内存对齐是什么？为什么需要？
Q8: Go 的内存模型（Memory Model）是什么？

---
Q1: Go 的内存分配器是怎么设计的？（TCMalloc）
【理解】
Go 的内存分配器基于 Google 的 TCMalloc（Thread-Caching Malloc）设计，核心思想是多级缓存减少锁竞争。

三级结构：
  mcache（P 级缓存）-> mcentral（全局中心缓存）-> mheap（堆）

分配流程（以小对象为例）：
  Step1: 根据对象大小找到对应的 size class
  Step2: 从当前 P 的 mcache 中对应 size class 的 span 分配（无锁）
  Step3: mcache 的 span 用完了 -> 从 mcentral 获取新 span（加锁）
  Step4: mcentral 也没有 -> 从 mheap 分配新 span（加锁）
  Step5: mheap 也不够 -> 向 OS 申请新内存（mmap）

设计目标：
  - 小对象分配尽量无锁（mcache 是 P 私有的）
  - 减少内存碎片（size class 预定义大小）
  - 减少系统调用（批量向 OS 申请）
【回答】
Go 内存分配器基于 TCMalloc 设计，核心是三级缓存：mcache（P 私有，无锁）-> mcentral（全局，有锁）-> mheap（堆，有锁）。
小对象优先从 mcache 分配，完全无锁；mcache 用完从 mcentral 补充；mcentral 不够从 mheap 分配；mheap 不够向 OS 申请。
设计目标：小对象分配无锁化、减少碎片（size class）、减少系统调用（批量申请）。

---
Q2: mcache、mcentral、mheap 分别是什么？
【理解】
■ mcache（每个 P 一个，无锁）：
  - 缓存各种 size class 的 span
  - 分配小对象时直接从这里取，不需要加锁
  - 类似于"每个收银台的零钱盒"

■ mcentral（每个 size class 一个，有锁）：
  - 管理某个 size class 的所有 span
  - 分为两个链表：有空闲对象的 span、已满的 span
  - mcache 用完时从这里补充
  - 类似于"银行柜台，各收银台来这里换零钱"

■ mheap（全局唯一，有锁）：
  - 管理所有未分配的内存页
  - mcentral 不够时从这里切割新 span
  - 大对象（>32KB）直接从 mheap 分配
  - 负责向 OS 申请/归还内存
  - 类似于"金库"

层级关系：
  P.mcache -> mcentral[sizeClass] -> mheap -> OS
  越往下锁越重，但访问频率越低
【回答】
mcache：每个 P 私有，缓存各 size class 的 span，小对象分配无锁直接取。
mcentral：每个 size class 一个，管理该规格的所有 span，mcache 用完来这补充（有锁）。
mheap：全局唯一，管理所有内存页，mcentral 不够时切割新 span，大对象直接从这分配。
层级越深锁越重但访问越少，大部分分配在 mcache 层就完成了。

---
Q3: 什么是 span（mspan）？size class 是什么？
【理解】
■ mspan：
  Go 内存管理的基本单位，是一组连续的内存页（page，8KB/页）。
  一个 span 被切割成多个相同大小的对象槽位。

  type mspan struct {
      startAddr uintptr   // 起始地址
      npages    uintptr   // 页数
      freeindex uintptr   // 下一个空闲对象的索引
      allocBits *gcBits   // 位图，标记哪些槽位已分配
      spanclass spanClass // size class 编号
  }

  例：一个 span 有 1 页（8KB），size class 是 32B
      -> 这个 span 被切成 8192/32 = 256 个 32 字节的槽位

■ size class：
  Go 预定义了 ~67 种大小规格（8B、16B、32B、48B、64B...32KB）。
  分配对象时，向上取整到最近的 size class。
  例：申请 25 字节 -> 分配 32 字节的槽位（浪费 7 字节，但减少碎片种类）

  对象大小分类：
    tiny（<16B 且无指针）：用 tiny allocator，多个小对象共享一个 16B 槽位
    small（16B~32KB）：按 size class 从 mcache 分配
    large（>32KB）：直接从 mheap 分配，不走 size class
【回答】
mspan 是内存管理的基本单位，一组连续内存页被切成相同大小的对象槽位。
size class 是预定义的 ~67 种大小规格（8B~32KB），分配时向上取整到最近的规格，用少量内部碎片换取管理简单。
对象按大小分三类：tiny（<16B 无指针，多个共享一个槽位）、small（16B~32KB，按 size class 分配）、large（>32KB，直接从 mheap 分配）。

---
Q4: 小对象、大对象的分配路径分别是什么？
【理解】
■ Tiny 对象（<16B 且无指针）：
  mcache 有一个 tiny 指针，指向当前 tiny 块（16B）。
  多个 tiny 对象共享同一个 16B 块，通过偏移量分配。
  例：分配 3 个 bool（各 1B）-> 都在同一个 16B 块里
  好处：极大减少小对象的内存浪费和分配次数。

■ Small 对象（16B~32KB）：
  Step1: 计算 size class（向上取整）
  Step2: 从 P.mcache 对应 size class 的 span 取一个空闲槽位（无锁）
  Step3: span 满了 -> 从 mcentral 换一个有空闲的 span
  Step4: mcentral 没有 -> 从 mheap 分配新 span

■ Large 对象（>32KB）：
  直接从 mheap 分配，按页对齐。
  不走 mcache 和 mcentral（太大了缓存没意义）。
  分配时需要加锁。

分配器优化：
  - mcache 无锁（P 私有）
  - span 用位图（allocBits）追踪空闲槽位，O(1) 找到下一个空闲位
  - tiny allocator 把多个微小对象打包，减少 overhead
【回答】
Tiny（<16B 无指针）：多个对象共享一个 16B 块，极大减少浪费。
Small（16B~32KB）：按 size class 从 mcache 的 span 取空闲槽位，无锁；span 满了从 mcentral 补充。
Large（>32KB）：直接从 mheap 按页分配，不走缓存，需要加锁。
大部分对象是 small，在 mcache 层无锁完成分配，性能很好。

---
Q5: 栈内存和堆内存有什么区别？
【理解】
栈（Stack）：
  - 每个 goroutine 私有，初始 2KB，动态伸缩
  - 分配/回收极快：只需移动 SP 指针
  - 函数返回自动回收，不需要 GC
  - 存放：局部变量、函数参数、返回值

堆（Heap）：
  - 所有 goroutine 共享
  - 分配需要走内存分配器（mcache/mcentral/mheap）
  - 回收需要 GC（标记-清除）
  - 存放：逃逸的变量、make 创建的 slice/map/channel

性能差异：
  栈分配：~1ns（移动指针）
  堆分配：~25ns（走分配器）+ GC 回收开销

编译器通过逃逸分析决定变量放栈还是堆：
  能留在栈上就留在栈上（零成本回收）
  必须逃逸才放堆上（需要 GC 管理）
【回答】
栈：goroutine 私有，分配只需移动 SP 指针（~1ns），函数返回自动回收不需要 GC。
堆：所有 goroutine 共享，分配走内存分配器（~25ns），回收需要 GC。
编译器通过逃逸分析决定放哪：能留栈上就留栈上（零成本），必须逃逸才放堆上。
减少堆分配 = 减少 GC 压力 = 提升性能。

---
Q6: 逃逸分析的规则？哪些情况会逃逸到堆？
【理解】
逃逸分析是编译器在编译期判断变量生命周期的静态分析。

会逃逸的场景：
  1. 返回局部变量的指针：函数返回后还要用，不能随栈帧销毁
  2. 闭包捕获：闭包可能比当前函数活得久
  3. 赋值给接口：interface 内部用指针存数据
  4. 发送到 channel：跨协程，编译器无法确定生命周期
  5. slice/map 太大：超过栈帧限制
  6. 动态大小分配：make([]byte, n)，n 是变量时编译器不知道大小

不逃逸的场景：
  1. 局部变量只在函数内使用
  2. 值拷贝传参（即使结构体大一点）
  3. 编译器能确定生命周期不超过当前函数

查看逃逸分析：
  go build -gcflags="-m" ./...
  go build -gcflags="-m -m" ./...  // 更详细的原因
【回答】
逃逸分析是编译器判断变量放栈还是堆的静态分析。
会逃逸：返回局部变量指针、闭包捕获、赋值给接口、发送到 channel、动态大小分配。
不逃逸：局部变量只在函数内用、值拷贝传参、编译器能确定生命周期。
用 go build -gcflags="-m" 查看逃逸结果，针对热路径优化减少逃逸。

---
Q7: 内存对齐是什么？为什么需要？
【理解】
内存对齐 = 变量的地址必须是其对齐值的整数倍。

Go 的对齐规则：
  类型        大小    对齐值
  bool        1B      1
  int8        1B      1
  int16       2B      2
  int32       4B      4
  int64       8B      8
  float64     8B      8
  pointer     8B      8（64位系统）
  string      16B     8（内部是指针+长度）

结构体对齐：
  type Bad struct {   // 占 24 字节（有 padding）
      a bool    // 1B + 7B padding
      b int64   // 8B
      c bool    // 1B + 7B padding
  }
  type Good struct {  // 占 16 字节（紧凑排列）
      b int64   // 8B
      a bool    // 1B
      c bool    // 1B + 6B padding
  }

为什么需要对齐？
  1. CPU 访问对齐地址只需一次内存读取，未对齐可能需要两次
  2. 某些 CPU 架构不支持未对齐访问（直接 panic）
  3. atomic 操作要求地址对齐（否则不是原子的）

优化：把大字段放前面，小字段放后面，减少 padding。
工具：fieldalignment（go vet 的一部分）可以自动检查。
【回答】
内存对齐要求变量地址是其对齐值的整数倍。CPU 访问对齐地址只需一次读取，未对齐可能需要两次或直接报错。
结构体字段顺序影响大小：把大字段放前面小字段放后面可以减少 padding。
例：bool+int64+bool 占 24B（有 padding），int64+bool+bool 只占 16B。
用 fieldalignment 工具自动检查和优化结构体字段顺序。

---
Q8: Go 的内存模型（Memory Model）是什么？
【理解】
Go Memory Model 定义了在多 goroutine 环境下，一个 goroutine 的写操作何时对另一个 goroutine 可见。

核心概念：happens-before
  如果事件 A happens-before 事件 B，那么 A 的效果对 B 可见。
  没有 happens-before 关系的操作，可见性不保证（可能被重排）。

Go 保证的 happens-before 关系：
  1. 同一个 goroutine 内：按代码顺序
  2. channel 发送 happens-before 对应的接收完成
  3. channel 关闭 happens-before 因关闭而返回的接收
  4. sync.Mutex Unlock happens-before 下一次 Lock
  5. sync.Once Do 的执行 happens-before 任何 Do 的返回
  6. sync.WaitGroup Done happens-before Wait 返回

不保证的场景：
  var x, y int
  go func() { x = 1; y = 2 }()
  go func() { fmt.Println(y, x) }()
  // 可能输出 "2 0"！编译器/CPU 可能重排 x=1 和 y=2

解决方法：用 channel、sync 包、atomic 包建立 happens-before 关系。
【回答】
Go 内存模型定义了多 goroutine 环境下写操作何时对其他 goroutine 可见，核心是 happens-before 关系。
保证的：channel 收发、Mutex Lock/Unlock、Once Do、WaitGroup Done/Wait 都建立 happens-before。
不保证的：没有同步原语的跨 goroutine 读写，编译器和 CPU 可能重排指令，导致看到"不可能"的中间状态。
规则：用 channel、sync、atomic 建立 happens-before 关系，不要依赖裸变量的可见性。

*/

// TestMemoryAlign 内存对齐演示
func TestMemoryAlign(t *testing.T) {
	type Bad struct {
		a bool
		b int64
		c bool
	}
	type Good struct {
		b int64
		a bool
		c bool
	}

	fmt.Printf("Bad  size: %d, align: %d\n", unsafe.Sizeof(Bad{}), unsafe.Alignof(Bad{}))
	fmt.Printf("Good size: %d, align: %d\n", unsafe.Sizeof(Good{}), unsafe.Alignof(Good{}))
}

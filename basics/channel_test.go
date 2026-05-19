package basics

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

/*
Q1: Channel 底层数据结构？
Q2: Channel 是线程安全的吗？
Q3: Channel 缓冲区 buf 底层数据结构？这么设计有什么好处？
Q4: 有缓冲和无缓冲的 channel 分别是啥？
Q5: channel 有什么要注意的？（三种 panic）
Q6: Channel 怎么会导致协程泄漏？
Q7: 如何优雅地关闭 Channel？有多个发送者，一个接收者，怎么关闭？
Q8: channel 发送数据的完整流程？
Q9: channel 接收数据的完整流程？
Q10: 协程间的直接拷贝有什么好处？
Q11: select 多路复用的实现原理？
Q12: close 操作底层做了什么？
Q13: channel 是值类型还是引用类型？传参会拷贝吗？
Q14: channel 和 Mutex 怎么选？
Q15: channel 怎么保证内存可见性（happens-before）？
Q16: 单向 channel 有什么意义？
Q17: 如何用 channel 实现并发限流（信号量）？
Q18: sudog 是什么？为什么需要它？
Q19: for range channel 的退出条件？
Q20: channel 常见的死锁场景有哪些？
Q21: goroutine 底层的 g 结构体由什么组成？
Q22: gopark 和 goready 是什么？协程阻塞和唤醒的底层机制？

---
Q1: Channel 底层数据结构？
【理解】
channel 底层数据结构：buf 核心数组（环形数组）+ sendq、recvq 两个阻塞队列（双向链表）+ 一个互斥锁
type hchan struct {
    qcount   uint           // 当前缓冲区中的数据个数
    dataqsiz uint           // 环形缓冲区的大小（make 时定义的容量）
    buf      unsafe.Pointer // 指向环形缓冲区的指针（只对有缓冲 channel 有效）
    elemsize uint16         // 元素大小
    closed   uint32         // 是否关闭
    elemtype *_type         // 元素类型
    sendx    uint           // 写入数据的指针下标（环形索引）
    recvx    uint           // 读取数据的指针下标（环形索引）
    recvq    waitq          // 等待读的协程队列（Sudog 双向链表）
    sendq    waitq          // 等待写的协程队列（Sudog 双向链表）
    lock     mutex          // 互斥锁，保护 hchan 的所有字段
}
buf 是环形数组（Ring Buffer），不是链表。数组在内存上是连续的，通过 sendx 和 recvx 循环移动，效率极高。
内存模型：不要通过共享内存来通信，而要通过通信来共享内存。
关闭原则：由发送方关闭，不要在接收方关闭，也不要在有多个发送方时关闭。
【回答】
Channel 底层是一个叫 hchan 的结构体，主要由三部分组成：
一是环形缓冲区 buf，这是一个数组，用来存带缓冲的 channel 数据；
二是两个等待队列 recvq 和 sendq，是双向链表，存的是因为没数据读或没地方发而阻塞的协程（sudog）；
三是互斥锁 lock，Channel 内部也是靠锁来保证并发安全的，只不过封装得很好，让我们用起来像在直接传值。

---
Q2: Channel 是线程安全的吗？
【理解】
hchan 自带 mutex，所有对 buf、sendq、recvq 的操作都在锁内完成。
【回答】
因为它的底层结构体 hchan 里自带了一把 Mutex 互斥锁。每一次入队、出队或者操作等待队列时，Go 运行时都会先获取这把锁，保证同一时间只有一个协程能修改 Channel 的状态。

---
Q3: Channel 缓冲区 buf 底层数据结构？这么设计有什么好处？
【理解】
数组实现的环形缓冲区，通过 sendx 和 recvx 两个索引移动实现循环读写。
对比普通数组：出队后不需要挪动后续元素，O(1) 入队出队。
对比链表：内存连续，CPU 缓存友好，无额外指针开销。
【回答】
buf 是一个数组实现的环形缓冲区，可以复用内存。通过 sendx 和 recvx 两个索引移动，避免了像普通数组那样删除元素后需要移动后面所有数据的问题，保证了 O(1) 的入队出队速度。而且数组内存连续，对 CPU 缓存友好。

---
Q4: 有缓冲和无缓冲的 channel 分别是啥？
【理解】
无缓冲：dataqsiz=0，buf 为空。发送方和接收方必须"手递手"交接，底层可能触发直接内存拷贝（跳过 buf）。
有缓冲：dataqsiz=N，buf 有 N 个槽位。数据先拷进 buf 再由接收方拷走，异步解耦。
在同一个 Goroutine 里对无缓冲 channel 又发又收，必然死锁。
【回答】
无缓冲 channel 是同步的，发送方和接收方必须"手递手"交接，如果没有协程在对端等待就会阻塞。底层有时会触发直接内存拷贝，从发送协程的栈直接拷到接收协程的栈，延迟极低。
有缓冲 channel 是异步的，数据先拷贝进 hchan.buf，再由接收方拷走。只要缓冲区没满，发送方就不会阻塞。

---
Q5: channel 有什么要注意的？（三种 panic）
【理解】
操作        nil Channel    Closed Channel       Normal Channel
Close       Panic          Panic                正常关闭
Send        永久阻塞       Panic                阻塞或进入缓冲区
Receive     永久阻塞       读完旧数据后返回零值  阻塞或取走数据
【回答】
三个要注意的点：
给一个 nil 的 channel 发送或接收数据，会永久阻塞（不是 panic）。
给一个已关闭的 channel 发送数据，会直接 panic。
从一个已关闭的 channel 接收数据，会先读完剩余缓冲区数据，然后永远返回零值和 ok=false。

---
Q6: Channel 怎么会导致协程泄漏？
【理解】
goroutine 被阻塞在 channel 操作上且永远不会被唤醒时，GC 无法回收它（因为 channel 等待队列中有对该 g 的引用）。
常见场景：生产者发完忘关 channel，消费者 range 永远等待；或者消费者不存在了，生产者向无缓冲/满 channel 发送。
【回答】
如果协程在等待一个永远没人发的 chan，或者发向一个没人收且没缓冲的 chan，该协程会永久阻塞在后台，无法被 GC 回收，这就是协程泄漏。本质是 channel 等待队列持有对 goroutine 的引用，导致它永远不会被回收。

---
Q7: 如何优雅地关闭 Channel？有多个发送者，一个接收者，怎么关闭？
【理解】
核心原则：谁的数据源枯竭了，谁负责关闭。永远不要在接收端关闭，除非你是唯一的发送者。
多发送者场景不能直接 close dataCh（可能有人还在发，会 panic）。
方案 A：sync.Once 包一层，确保只 close 一次。
方案 B：引入独立的 stopCh，接收者想停时 close(stopCh) 做广播，发送者 select 监听 stopCh 后主动退出。
【回答】
直接关闭会 panic，因为可能有其他发送者还在发。
方案 A：用 sync.Once 确保只关一次。
方案 B 更专业：引入一个额外的 stopCh。接收者想停的时候关闭 stopCh，所有的发送者通过 select 监听 stopCh，一旦发现 stopCh 关了就停止发送。等所有发送者退出后，再安全关闭 dataCh。
核心原则：谁的数据源枯竭了，谁负责关闭。永远不要在接收端关闭，除非你是唯一的发送者。

---
Q8: channel 发送数据的完整流程？（发送 = 往 channel 写入数据，即 ch <- value）
【理解】
源码路径：runtime/chan.go -> chansend 函数。
"发送"就是执行 ch <- value 时，runtime 内部走的流程。
分有缓冲和无缓冲两种情况：

■ 有缓冲 channel（ch := make(chan int, N)）：
  Step1: 获取锁 lock(&c.lock)
  Step2: 检查 channel 是否已关闭，closed=1 则 panic("send on closed channel")
  Step3: 看 recvq 是否有等待的接收者（有人正卡在 <-ch 等数据）
         -> 有：直接把数据拷贝到接收者的内存地址，唤醒它，不经过 buf（最快路径）
         -> 没有：往下走
  Step4: 看 buf 是否还有空位（qcount < dataqsiz）
         -> 有空位：把数据写入 buf[sendx]，sendx = (sendx+1) % dataqsiz，qcount++
         -> 没有空位：往下走
  Step5: buf 满了，也没人在等接收
         -> 把当前 goroutine 包装成 sudog{g: curG, elem: &data}，挂到 sendq 链表
         -> 调用 gopark 让出 CPU，进入阻塞
         -> 直到某个接收方从 buf 取走数据后唤醒自己
  Step6: 释放锁 unlock(&c.lock)

  场景推演：ch := make(chan int, 2)
    ch<-1：recvq 空，buf 有空位 -> buf=[1,_], sendx=1, qcount=1
    ch<-2：recvq 空，buf 有空位 -> buf=[1,2], sendx=0(环形回绕), qcount=2
    ch<-3：recvq 空，buf 满 -> 当前 goroutine 挂到 sendq 阻塞

■ 无缓冲 channel（ch := make(chan int)）：
  Step1: 获取锁 lock(&c.lock)
  Step2: 检查 channel 是否已关闭，closed=1 则 panic("send on closed channel")
  Step3: 看 recvq 是否有等待的接收者
         -> 有：直接把数据拷贝到接收者的内存地址（memmove），唤醒它
              这就是"手递手"直接交付，只需 1 次内存拷贝，性能最优
         -> 没有：往下走
  Step4: 没有 buf（dataqsiz=0），无处可放
         -> 把当前 goroutine 包装成 sudog 挂到 sendq 链表
         -> 调用 gopark 让出 CPU，进入阻塞
         -> 直到某个接收方执行 <-ch 时发现 sendq 有自己，直接拷走数据并唤醒
  Step5: 释放锁 unlock(&c.lock)

  场景推演：ch := make(chan int)
    goroutine B 先执行 v := <-ch -> B 被挂到 recvq 等待
    goroutine A 执行 ch <- 42 -> 发现 recvq 有 B，直接把 42 拷到 B 的变量 v 地址，唤醒 B
    如果 A 先到，没人等接收 -> A 挂到 sendq 阻塞，直到 B 来 <-ch 把 A 唤醒

【回答】
发送就是 ch <- value，往 channel 写入数据。分有缓冲和无缓冲两种情况：
有缓冲的情况：首先获取 hchan 的锁，然后检查 channel 是否关闭（关闭则 panic）。
接着看 recvq 有没有等待的接收者，有的话直接把数据拷到接收者内存里并唤醒它，不走 buf。
如果没有等待的接收者，就看 buf 有没有空位，有就写入 buf[sendx]，sendx 环形自增，qcount++。
如果 buf 满了，就把自己包装成 sudog 挂到 sendq 阻塞，等接收方取走数据后才会被唤醒。最后释放锁。
无缓冲的情况：同样先获取锁、检查关闭。然后看 recvq 有没有等待的接收者，有就直接把数据拷贝到接收者的内存地址并唤醒它（手递手直接交付，1 次内存拷贝）。
如果没有接收者在等，因为无缓冲没有 buf 可放，直接把自己挂到 sendq 阻塞，等接收方到了才会被唤醒。最后释放锁。

---
Q9: channel 接收数据的完整流程？（接收 = 从 channel 读取数据，即 v := <-ch）
【理解】
源码路径：runtime/chan.go -> chanrecv 函数。
"接收"就是执行 v := <-ch 或 v, ok := <-ch 时，runtime 内部走的流程。
举例：ch := make(chan int, 2)，goroutine B 执行 v := <-ch
流程（按优先级依次判断）：
  Step1: 获取锁 lock(&c.lock)
  Step2: 检查 channel 是否已关闭 且 buf 为空
         -> 是：返回零值（v=0）和 ok=false，表示"没数据且永远不会有了"
  Step3: 看 sendq 是否有等待的发送者（有人正卡在 ch<- 等着发数据）
         -> 无缓冲 channel：直接从发送者内存拷贝数据到 v，唤醒发送者
         -> 有缓冲 channel：先从 buf[recvx] 取数据给 v，再把 sendq 队首发送者
            的数据搬到 buf[sendx]（填补刚空出的位置），唤醒发送者
  Step4: 看 buf 是否有数据（qcount > 0）
         -> 有：从 buf[recvx] 取数据赋给 v，recvx = (recvx+1) % dataqsiz，qcount--
  Step5: buf 空了，也没人在等发送
         -> 把当前 goroutine B 包装成 sudog{g: B, elem: &v}，挂到 recvq 链表
         -> 调用 gopark 让出 CPU，B 进入阻塞状态
         -> 直到某个发送方 ch<-value 时发现 recvq 有 B，把数据拷到 &v 并唤醒 B
  Step6: 释放锁 unlock(&c.lock)
具体场景推演：
场景 A（buf 有数据）：ch := make(chan int, 2); ch<-1; ch<-2; v := <-ch
  buf=[1,2], recvx=0
  执行 <-ch：v=1, buf=[_,2], recvx=1, qcount=1

场景 B（buf 满 + 有发送者阻塞）：ch := make(chan int, 1); ch<-1; go ch<-2(阻塞在sendq)
  执行 v := <-ch：
    先从 buf[0] 取出 1 给 v
    再把 sendq 队首（值为2）搬到 buf[0]（填补空位），唤醒发送者
    结果：v=1, buf=[2], 发送者被唤醒继续执行
  这样保证了 FIFO：先发的 1 先被收到

场景 C（无缓冲直接交付）：ch := make(chan int)
  goroutine A 先执行 ch <- 42，A 被挂到 sendq 等待
  goroutine B 执行 v := <-ch：发现 sendq 有 A，直接从 A 的内存拷贝 42 到 v，唤醒 A

场景 D（channel 已关闭）：close(ch)，然后 v, ok := <-ch
  如果 buf 还有残留数据：正常取出，ok=true
  如果 buf 空了：v=0（零值），ok=false
【回答】
接收就是 v := <-ch，从 channel 读取数据。底层流程按优先级：
首先获取锁。然后检查 channel 是否关闭且 buf 为空，是的话直接返回零值和 ok=false。
接着看 sendq 有没有阻塞的发送者：无缓冲场景直接从发送者内存拷数据；有缓冲场景先从 buf 取数据，再把 sendq 队首的数据搬进 buf 填补空位，唤醒发送者（保证 FIFO 顺序）。
如果没有等待的发送者，就看 buf 有没有数据，有就从 buf[recvx] 取。
如果 buf 也空了，就把自己包装成 sudog 挂到 recvq 阻塞，等发送方写入数据后唤醒自己。
最后释放锁。

---
Q10: 协程间的直接拷贝有什么好处？
【理解】
常规路径：发送者 -> buf -> 接收者，2 次内存拷贝。
优化路径：发送者 -> 直接拷贝到接收者，1 次内存拷贝。
触发条件：无缓冲 channel 发送时 recvq 有等待者，或接收时 sendq 有等待者。
本质是 runtime 直接操作目标 goroutine 栈上的变量地址。
【回答】
对于无缓冲 Channel，Go 优化到了极致：直接从发送协程的内存拷到接收协程的内存，中间不落地。
常规逻辑需要两次拷贝：发送者 -> 缓冲区 -> 接收者。Go 的设计是当发现等待队列有协程时，直接拷贝数据到等待协程的内存地址，只需 1 次内存拷贝，极大提升了性能，尤其是对于大数据量的传递。

---
Q11: select 多路复用的实现原理？
【理解】
源码路径：runtime/select.go -> selectgo 函数。
核心设计：随机打乱避免饥饿 + 按 channel 地址排序加锁避免死锁。
一个 g 可同时挂在多个 channel 的 waitq 上，唤醒后需要从其他队列摘除。
空 select{} 永久阻塞当前 goroutine。

完整流程（以 3 个 case 为例）：
  select {
  case v := <-ch1:  // case0
  case v := <-ch2:  // case1
  case ch3 <- data: // case2
  }

  第一阶段（快速路径）：
    1. 把所有 case 收集成 scase 数组
    2. 随机打乱遍历顺序（避免饥饿，不是总偏向第一个 case）
    3. 依次检查每个 case 能否立即执行（buf有数据/对端有等待者/已关闭）
    4. 如果有就绪的 case，直接执行并返回
    5. 如果有 default，直接走 default 返回

  第二阶段（阻塞路径，无 default 才走）：
    1. 按 channel 地址排序，依次对所有 channel 加锁（避免多个 select 互相死锁）
    2. 创建多个 sudog，分别挂到各 channel 的等待队列：
       sudog1 -> ch1.recvq
       sudog2 -> ch2.recvq
       sudog3 -> ch3.sendq
    3. gopark 阻塞当前 G，让出 CPU

  第三阶段（唤醒后）：
    假设 ch2 先有数据到达，ch2 的发送方发现 recvq 有 sudog2：
    1. 把数据拷贝到 sudog2.elem 指向的地址，标记 sudog2.success = true
    2. goready 唤醒 G
    3. G 醒来后，遍历其他 case 对应的 channel 等待队列：
       - 从 ch1.recvq 中摘除 sudog1
       - 从 ch3.sendq 中摘除 sudog3
    4. 执行 case1 分支的代码

  为什么必须摘除其他队列上的 sudog？
    如果不摘，G 已经醒了在执行别的代码，但 ch1 来数据时又发现队列里有这个 G 的 sudog，
    再次往它的 elem 地址写数据，就会出现数据竞争或内存破坏（幽灵唤醒）。

  关键字段：
    type sudog struct {
        isSelect bool  // 标记来自 select，唤醒时需要走摘除逻辑
        success  bool  // 标记是否是"赢家"（被选中的 case）
    }

【回答】
所有 case 收集成 scase 数组，先对 case 顺序做随机打乱来避免饥饿。
第一轮按打乱顺序遍历，看哪个 case 能立即执行，如果有就绪的就直接执行。有 default 则走 default。
第二轮（无 default 才走）：按 channel 地址排序加锁，把当前协程打包成多个 sudog，分别挂到所有 case 涉及的 channel 等待队列上，gopark 阻塞。
唤醒后：只有一个 case 能"中奖"执行，G 醒来后立刻遍历其他 channel 的等待队列，把自己的 sudog 全部摘除，避免幽灵唤醒（防止其他 channel 后续再往已失效的地址写数据导致数据竞争）。
一句话：select 阻塞时一个 G 同时排多个队，但只有一个 case 能触发，触发后立刻把其他队的号全撤掉。

---
Q12: close 操作底层做了什么？
【理解】
源码路径：runtime/chan.go -> closechan 函数。
关键点：唤醒所有等待者，接收者收到零值，发送者会 panic。
close nil channel：直接 panic（hchan 指针为空，没法操作）。
close 已关闭 channel：panic（closed 字段已为 1）。
【回答】
加锁，标记 closed=1，把 recvq 里所有等待的接收者全部唤醒（它们读到零值和 ok=false），把 sendq 里所有等待的发送者全部唤醒（它们醒来后会 panic：send on closed channel），最后解锁。
所以重复 close 会 panic（closed=1 已置位），close nil channel 也会 panic（hchan 指针为空）。

---
Q13: channel 是值类型还是引用类型？传参会拷贝吗？
【理解】
make(chan T) 编译后调用 runtime.makechan，返回 *hchan。
所以 channel 变量本质是指针，和 map 一样。
传参时拷贝指针值，所有 goroutine 共享同一个 hchan 实例。
【回答】
make(chan T) 返回的是 *hchan 指针。channel 变量本身就是个指针，传参时拷贝的是指针，所有 goroutine 操作的是同一个底层 hchan，所以才能跨协程通信。
这也是为什么 channel 不需要用 &ch 传递，和 map 同理。

---
Q14: channel 和 Mutex 怎么选？
【理解】
channel 内部本身就有锁，额外还有 goroutine 调度开销，性能不一定优于直接用 Mutex。
官方建议：https://go.dev/wiki/MutexOrChannel
数据所有权转移 -> channel；保护共享状态 -> Mutex。
【回答】
共享状态、保护字段、临界区短、追求性能用 Mutex。
协程间传递数据/事件、表达流水线、解耦生产消费、控制并发度用 channel。
经验法则："共享内存用锁，传递所有权用 channel"。
补充：channel 内部本身就有锁，性能并不一定优于 Mutex，不要无脑用 channel。

---
Q15: channel 怎么保证内存可见性（happens-before）？
【理解】
Go 内存模型定义了哪些操作之间有 happens-before 关系，有这个关系就保证不会被重排、写入一定可见。
Go Memory Model 规定的 channel 相关 happens-before 关系：
- 对 buffered channel 的第 n 次发送 happens-before 第 n 次接收完成。
- 对 unbuffered channel 的第 n 次接收 happens-before 第 n 次发送完成。
- channel 的 close happens-before 因 close 而返回零值的接收。
本质：channel 操作包含 lock/unlock，lock 自带内存屏障。
【回答】
Go 内存模型规定：对 channel 的第 n 次发送 happens-before 对应的第 n 次接收完成；channel 的关闭 happens-before 因关闭而返回的接收。
通俗讲：发送方在发之前对共享变量的所有写，接收方在收到之后都能看到，channel 通信自带内存屏障，无需额外加锁。

---
Q16: 单向 channel 有什么意义？
【理解】
chan<- T 只发，<-chan T 只收。纯编译期约束，运行时底层都是同一个 *hchan。
防御性编程：函数签名用单向 channel 防止误操作（误关、误读/写）。
双向可隐式转单向，反过来不行（编译报错）。
【回答】
chan<- T 只发，<-chan T 只收，是编译期类型约束，没有运行时开销。
最佳实践：函数参数用单向 channel，明确职责，防止函数内部误关 channel 或误向接收方写入。
双向 channel 可以隐式转单向，反过来不行。

---
Q17: 如何用 channel 实现并发限流（信号量）？
【理解】
核心思路：容量为 N 的 buffered channel 天然就是一个计数信号量。
发送 = P 操作（获取信号量），接收 = V 操作（释放信号量）。
优势：无需引入第三方库，语义清晰，天然协程安全。
【回答】
用容量为 N 的 buffered channel 当令牌桶。sem := make(chan struct{}, N)，进入临界区前 sem<-struct{}{} 占坑（满了会阻塞），结束时 <-sem 释放。并发度天然被限制为 N，且不需要额外锁。

---
Q18: sudog 是什么？为什么需要它？
【理解】
sudog 读作 su-dog，是 "pseudo G" 的缩写（pseudo = 伪），意思是"伪 G"、"G 的代理"。
本质是 goroutine 在某个 channel 上的等待记录，是 goroutine 的包装形式。

为什么不直接把 g 挂到等待队列，要多包装一层？
因为一个 g 可以同时等待多个 channel（select 场景）：
  select {
  case v := <-ch1:   // 需要挂到 ch1.recvq
  case v := <-ch2:   // 同时也要挂到 ch2.recvq
  case ch3 <- data:  // 同时也要挂到 ch3.sendq
  }
一个 g 只有一组链表指针，没法同时挂在多个队列上。
所以每个等待位置创建一个 sudog，它们都指向同一个 g。
类比：g 是"人"，sudog 是"排队号"——一个人可以同时在多个窗口取号排队。

核心字段：
  type sudog struct {
      g        *g              // 所属的 goroutine
      elem     unsafe.Pointer  // 要发送/接收的数据地址
      c        *hchan          // 在哪个 channel 上等待
      prev     *sudog          // 双向链表前驱
      next     *sudog          // 双向链表后继
      isSelect bool            // 是否来自 select 语句
  }

sudog 由 P 本地缓存池（sudogcache）复用，避免频繁分配带来的 GC 压力。

【回答】
sudog 读作 su-dog，是 pseudo G（伪 G）的缩写，本质是 goroutine 在某个 channel 上的等待记录。
为什么不直接用 g？因为一个 g 在 select 中可以同时等待多个 channel，需要同时挂在多个等待队列上，但一个 g 只有一组链表指针。所以用 sudog 做中间层，每个等待位置一个 sudog，它们都指向同一个 g。
类比：g 是"人"，sudog 是"排队号"，一个人可以同时在多个窗口取号排队。
核心字段包括 g 指针、收发数据地址 elem、所属 channel、链表前后指针等。sudog 由 P 本地缓存池复用，减少 GC 压力。

---
Q19: for range channel 的退出条件？
【理解】
编译器把 for v := range ch 展开成 for { v, ok := <-ch; if !ok { break } ... }。
退出条件：channel 被 close 且 buf 数据全部被取完（ok=false）。
如果发送方不 close，range 永远阻塞 -> 死锁。
【回答】
range 会一直读，直到 channel 被 close 且 buf 数据全部被取完，才会退出循环。
所以发送方必须在发完后 close，否则 range 永远阻塞导致死锁。
range 自动处理 (v, ok) 中的 ok=false，使用者无感知。

---
Q20: channel 常见的死锁场景有哪些？
【理解】
Go 运行时在所有 goroutine 都阻塞时会触发 fatal error: all goroutines are asleep - deadlock!
注意：如果只是部分 goroutine 泄漏（还有其他活跃的），运行时不会报死锁，只是内存泄漏。
【回答】
主协程对无缓冲 channel 又发又收（同一个 g）：必死锁。
发送方忘记 close，接收方用 for range 等待：永久阻塞。
多个协程互相等对方的 channel：环形依赖死锁。
向 nil channel 收发：fatal error: all goroutines are asleep - deadlock!
select 所有 case 都是 nil channel 且无 default：死锁。

---
Q21: goroutine 底层的 g 结构体由什么组成？
【理解】
goroutine 在 runtime 里就是一个 g 结构体（runtime/runtime2.go），核心字段：
type g struct {
    stack       stack      // 栈内存：lo 和 hi 两个指针，描述栈的起止地址（初始 2KB，按需扩容）
    stackguard0 uintptr    // 栈溢出检查哨兵，用于触发栈扩容

    _panic      *_panic    // 当前 panic 链表
    _defer      *_defer    // 当前 defer 链表

    m           *m         // 当前绑定的 M（OS 线程），nil 表示没在运行
    sched       gobuf      // 调度上下文：保存 SP、PC、BP 等寄存器，用于协程切换和恢复

    goid        uint64     // goroutine ID
    status      uint32     // 状态：_Gidle/_Grunnable/_Grunning/_Gwaiting/_Gdead

    waitreason  waitReason // 阻塞原因（如 "chan receive"、"select"）

    preempt     bool       // 是否被标记为需要抢占

    waiting     *sudog     // 这个 g 当前关联的 sudog 链表头
}

g 的状态机：
  _Gidle → _Grunnable → _Grunning → _Gwaiting → _Grunnable → ...
    创建      进入队列      被M执行     阻塞(如chan)    被唤醒重新排队
                                          ↓
                                      _Gdead (结束)

为什么 g 不能直接挂到 channel 等待队列？
如果想把 g 直接作为链表节点挂到 channel 的等待队列，g 自身需要 prev/next 指针。
但一个结构体的一组 prev/next 只能让它存在于一条链表中。
select 场景中一个 g 要同时等多个 channel（多条队列），所以必须用 sudog 包装：
  ch1.recvq: ... <-> [sudog1 {g: G}] <-> ...
  ch2.recvq: ... <-> [sudog2 {g: G}] <-> ...
  三个 sudog 各自独立在不同链表中，都指向同一个 G，互不冲突。
类比：g 是"人"，sudog 是"排队号"，一个人可以同时在多个窗口取号排队。

【回答】
goroutine 底层是 runtime 里的 g 结构体，核心组成包括：
栈（stack）：记录栈的起止地址，初始 2KB 按需扩容；stackguard0 哨兵用于触发栈扩容。
调度上下文（sched gobuf）：保存 SP、PC 等寄存器，协程切换时靠它保存和恢复现场。
关联的 M：当前绑定的 OS 线程，nil 表示没在运行。
状态（status）：_Grunnable（可运行）、_Grunning（运行中）、_Gwaiting（阻塞中）等。
panic/defer 链表：当前协程的 panic 和 defer 调用链。
waiting 指针：指向当前关联的 sudog 链表。

g 不能直接挂到 channel 等待队列，因为一组 prev/next 指针只能让它在一条链表中。
select 场景中一个 g 要同时等多个 channel，所以需要 sudog 做包装层，每个等待位置一个 sudog，各自独立在不同链表中，都指向同一个 g。

---
Q22: gopark 和 goready 是什么？协程阻塞和唤醒的底层机制？
【理解】
gopark 和 goready 是 Go runtime 中协程阻塞/唤醒的底层函数对，是 channel、select、mutex、timer 等所有阻塞机制的共同底座。

■ gopark（阻塞当前协程）：
  源码路径：runtime/proc.go -> gopark()
  触发场景：channel 发送时 buf 满、接收时 buf 空、select 无就绪 case、mutex 竞争等。
  做了三件事：
    1. 当前 G 状态 _Grunning → _Gwaiting（标记为等待中）
    2. 当前 G 和 M（OS 线程）解绑，G 不再占用任何线程资源
    3. 调用 schedule()，让 M 去 P 的本地队列找下一个可运行的 G 来执行
  关键点：gopark 是协程级阻塞，不是线程级阻塞！
    - OS 线程（M）没有 sleep，它立刻去执行别的协程
    - 阻塞的 G 只是一个挂在等待队列上的结构体，不占线程资源
    - 这就是 Go 能开百万协程的原因

  执行流程：
    G1 执行 ch<-42（buf满）
      → gopark()
      → G1 状态: _Grunning → _Gwaiting
      → G1 与 M1 解绑
      → M1 调用 schedule() 去执行 G2
      → G1 安静地待在 sendq 里，代码停在 ch<-42 这行不动

■ goready（唤醒阻塞的协程）：
  源码路径：runtime/proc.go -> goready()
  触发场景：channel 接收时发现 sendq 有等待者、发送时发现 recvq 有等待者、mutex 解锁等。
  做了两件事：
    1. 目标 G 状态 _Gwaiting → _Grunnable（标记为可运行）
    2. 把 G 放回 P 的本地运行队列，等待被某个 M 调度执行

  执行流程：
    G2 执行 v := <-ch（从 buf 取走数据）
      → 发现 sendq 有 G1，把 G1 的数据搬进 buf
      → goready(G1)
      → G1 状态: _Gwaiting → _Grunnable
      → G1 放回 P 的运行队列
      → 某个 M 拿到 G1 后，从 ch<-42 的下一行继续执行

■ 完整配对流程（以 channel 发送阻塞为例）：
  G1: ch<-42 → buf满 → gopark() → G1 阻塞在 sendq
        ...（G1 停在这，M 去忙别的）...
  G2: v:=<-ch → 取走 buf 数据 → 发现 sendq 有 G1 → goready(G1) → G1 恢复执行

■ 与操作系统线程阻塞的区别：
  OS 线程阻塞：线程 sleep，占用内核资源，上下文切换开销大（微秒级）
  gopark 协程阻塞：只改状态 + 换个 G 执行，纯用户态操作，开销极小（纳秒级）

【回答】
gopark 和 goready 是 Go runtime 中协程阻塞和唤醒的底层函数对，是 channel、select、mutex 等所有阻塞机制的共同底座。
gopark 让当前协程主动让出 CPU 并阻塞，做三件事：把 G 状态从 _Grunning 改为 _Gwaiting；把 G 和 M 解绑；让 M 调用 schedule() 去执行其他协程。阻塞的 G 只是挂在等待队列上的一个结构体，不占用任何线程资源，这就是 Go 能开百万协程的原因。
goready 把阻塞的 G 状态改回 _Grunnable，放回 P 的运行队列等待被调度。
和 OS 线程阻塞的本质区别：gopark 是纯用户态的协程切换（纳秒级），OS 线程不会 sleep；而系统调用导致的线程阻塞需要内核介入（微秒级），会真正占用线程资源。

*/

// TestChanBasics 标准遍历：展示 for range 的正确用法
func TestChanBasics(t *testing.T) {
	ch := make(chan string, 3)

	go func() {
		ch <- "Golang"
		ch <- "Channel"
		ch <- "Practice"
		// 发送完一定要关，否则 main 的 for range 会永久阻塞，导致死锁
		close(ch)
	}()

	// range 会一直读，直到 ch 被关闭且数据取完
	for v := range ch {
		fmt.Println("接收到:", v)
	}
	fmt.Println("遍历正常结束")
}

// TestChanReadAfterClosed 关闭后的读取：展示 (value, ok) 模式
func TestChanReadAfterClosed(t *testing.T) {
	ch := make(chan int, 1)
	ch <- 99
	close(ch)

	// 第一次读：能读到缓冲区剩下的值
	val1, ok1 := <-ch
	fmt.Printf("第一次读 - 值: %d, ok: %v (通道虽关，但数据还在)\n", val1, ok1)

	// 第二次读：缓冲区空了，且通道已关
	val2, ok2 := <-ch
	fmt.Printf("第二次读 - 值: %d, ok: %v (已关且空，读到零值)\n", val2, ok2)
}

// TestChanSpecialCases 各种会导致 Panic 的雷区
func TestChanSpecialCases(t *testing.T) {
	// --- 情况 A: 向已关闭的通道发数据 ---
	// c1 := make(chan int)
	// close(c1)
	// c1 <- 1 // panic: send on closed channel

	// --- 情况 B: 重复关闭通道 ---
	// c2 := make(chan int)
	// close(c2)
	// close(c2) // panic: close of closed channel

	// --- 情况 C: 关闭 nil 通道 ---
	// var c3 chan int
	// close(c3) // panic: close of nil channel

	// --- 情况 D: 读 nil 通道 (永久阻塞导致死锁) ---
	// var c4 chan int
	// <-c4 // fatal error: all goroutines are asleep - deadlock!

	fmt.Println("(Panic 示例代码已注释，可手动解开测试)")
}

// TestChanSelectTimeout select 多路复用 + 超时控制
func TestChanSelectTimeout(t *testing.T) {
	ch := make(chan int)

	go func() {
		time.Sleep(200 * time.Millisecond)
		ch <- 42
	}()

	// select 随机选一个就绪 case；都没就绪则阻塞
	select {
	case v := <-ch:
		fmt.Println("收到数据:", v)
	case <-time.After(100 * time.Millisecond):
		// time.After 返回 <-chan Time，到点后可读，相当于超时信号
		fmt.Println("超时退出，避免协程永久阻塞")
	}
}

// TestChanSemaphore 用 buffered channel 实现信号量限流
func TestChanSemaphore(t *testing.T) {
	const N = 3 // 最多 3 个协程同时执行
	sem := make(chan struct{}, N)
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			sem <- struct{}{}        // 占坑，满了就阻塞
			defer func() { <-sem }() // 离开时让出坑

			fmt.Printf("任务 %d 执行中\n", id)
			time.Sleep(50 * time.Millisecond)
		}(i)
	}
	wg.Wait()
}

// TestChanGracefulClose 多发送者 + 一个接收者的优雅关闭
func TestChanGracefulClose(t *testing.T) {
	dataCh := make(chan int, 10)
	stopCh := make(chan struct{}) // 关闭它 = 广播"停"
	var wg sync.WaitGroup

	// 启动 3 个发送者
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; ; j++ {
				select {
				case <-stopCh:
					return // 收到停止信号，主动退出
				case dataCh <- id*100 + j:
				}
			}
		}(i)
	}

	// 接收者：读够 5 个就通知发送者退出
	go func() {
		count := 0
		for v := range dataCh {
			fmt.Println("接收:", v)
			count++
			if count >= 5 {
				close(stopCh) // 广播停止信号
				break
			}
		}
	}()

	wg.Wait()     // 等所有发送者退出
	close(dataCh) // 此时已无发送者，安全关闭
	for range dataCh {
	} // 排空残留
	fmt.Println("优雅关闭完成")
}

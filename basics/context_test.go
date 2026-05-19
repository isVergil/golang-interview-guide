package basics

import (
	"context"
	"fmt"
	"testing"
	"time"
)

/*
Q1: Context 是什么？解决了什么问题？
Q2: Context 有哪些常见函数？分别什么场景用？
Q3: Context 的底层原理？取消信号是怎么传播的？
Q4: Context 是线程安全的吗？为什么？
Q5: 能不能把所有参数都塞进 Context.Value 里传递？
Q6: WithValue 的查找性能如何？底层数据结构？
Q7: 为什么 context.WithCancel 要返回一个 cancel 函数？不调用会怎样？
Q8: Context 使用有哪些最佳实践？

---
Q1: Context 是什么？解决了什么问题？
【理解】
Context 主要解决 goroutine 树状管理的问题。
当一个请求进来，可能会派生出一棵 goroutine 树（主协程 -> 子协程 -> 孙协程...）。
如果请求取消或超时，需要一种机制通知整棵树上所有协程优雅退出，而不是让它们变成泄漏的孤儿协程。

三大作用：
  1. 退出通知（取消信号）：手动或自动通知子协程停止工作
  2. 超时控制：规定一个操作必须在多少时间内完成
  3. 元数据传递：链路追踪 ID、用户 Token 等横切关注点

Context 的接口定义：
  type Context interface {
      Deadline() (deadline time.Time, ok bool)  // 截止时间
      Done() <-chan struct{}                     // 取消信号 channel
      Err() error                               // 取消原因
      Value(key any) any                        // 取值
  }

【回答】
Context 主要解决 goroutine 树状管理的问题。一个请求可能派生出一棵协程树，如果请求取消或超时，需要通知整棵树上所有协程优雅退出。
它有三大作用：退出通知（取消信号）、超时控制、以及元数据传递。
接口很简单：Done() 返回一个 channel 用于监听取消信号，Err() 返回取消原因，Deadline() 返回截止时间，Value() 用于传递元数据。

---
Q2: Context 有哪些常见函数？分别什么场景用？
【理解】
四个派生函数，每个都基于 parent context 创建子 context：
  WithCancel(parent)   -> 手动取消，返回 cancel 函数
  WithTimeout(parent, duration) -> 超时自动取消（相对时间）
  WithDeadline(parent, time)    -> 到期自动取消（绝对时间）
  WithValue(parent, key, val)   -> 附加键值对，不涉及取消

两个根 context：
  context.Background() -> 整个程序的根，main/init/测试 用
  context.TODO()       -> 还没想好用什么 context 时的占位符

WithTimeout 底层就是 WithDeadline(parent, time.Now().Add(timeout))，本质相同。
【回答】
四个派生函数：
WithCancel：手动取消，典型场景是任务完成后通知辅助协程退出。
WithTimeout：超时取消（相对时间），比如调用第三方 API 规定 2 秒必须返回。
WithDeadline：定时取消（绝对时间），比如任务必须在某个时刻前完成。
WithValue：传值，链路追踪的 TraceID 或用户的 UserID。
两个根：Background() 是整个程序的根 context，TODO() 是占位符。

---
Q3: Context 的底层原理？取消信号是怎么传播的？
【理解】
核心结构（以 cancelCtx 为例）：
  type cancelCtx struct {
      Context                        // 嵌入 parent context
      mu       sync.Mutex            // 保护下面的字段
      done     atomic.Value          // 存 chan struct{}，懒初始化
      children map[canceler]struct{} // 所有子 context
      err      error                 // 取消原因
  }

取消传播机制：
  Step1: 调用 cancel() 时，关闭自己的 done channel（close(c.done)）
  Step2: 所有监听 <-ctx.Done() 的协程收到零值，知道该退出了
  Step3: 遍历 children map，递归调用每个子 context 的 cancel()
  Step4: 整棵子树全部被取消

注册机制：
  子 context 创建时，会向 parent 的 children map 注册自己。
  这样 parent 取消时能找到所有子 context。

done channel 是懒初始化的：第一次调用 Done() 时才创建，节省内存。

【回答】
核心是 cancelCtx 结构体，里面有一个 done channel 和一个 children map。
取消传播流程：调用 cancel() 时，先关闭自己的 done channel，所有监听 <-ctx.Done() 的子协程都会收到信号。然后遍历 children map，递归调用所有子 context 的 cancel，确保整棵树都被清理。
子 context 创建时会向 parent 注册自己（写入 children map），这样 parent 才能在取消时找到所有后代。

---
Q4: Context 是线程安全的吗？为什么？
【理解】
Context 被设计为并发安全的：
  1. Context 接口本身是只读的（Done/Err/Deadline/Value 都是读操作）
  2. cancelCtx 内部用 sync.Mutex 保护 children map 和 err 字段
  3. done channel 用 atomic.Value 存储，保证并发访问安全
  4. WithValue 创建的是新节点（不可变链表），不修改已有节点

所以可以在多个 goroutine 之间传递同一个 ctx 而不需要加锁。
【回答】
Context 是线程安全的。它被设计为不可变的接口——Done/Err/Value 都是只读操作。
内部 cancelCtx 用 sync.Mutex 保护 children map，done channel 用 atomic.Value 存储。
WithValue 创建的是新节点不会修改原节点。所以可以在多个 goroutine 之间传递同一个 ctx 而不需要加锁。

---
Q5: 能不能把所有参数都塞进 Context.Value 里传递？
【理解】
绝对不行。WithValue 应该只传"非业务核心数据"。

问题：
  1. 性能差：查找是 O(n)，链表结构从叶子往根一层层找
  2. 类型不安全：interface{} 类型，取出来要类型断言，编译期抓不住错误
  3. 隐藏依赖：函数签名看不出它需要什么数据，可读性极差
  4. 难以测试：mock 困难，不知道要塞什么 key 进去

正确用法：只传 RequestID、TraceID、UserToken 这种横向贯穿整个链路的元数据。
业务参数应该显式写在函数签名里。
【回答】
绝对不行。WithValue 应该只传非业务核心的横切数据。
三个原因：性能差——查找是 O(n) 的链表遍历；类型不安全——interface{} 要断言，容易出错；隐藏依赖——函数签名看不出需要什么数据。
正确用法：只传 RequestID、TraceID、UserToken 这种横向贯穿整个链路的元数据。业务参数应该显式写在函数签名里。

---
Q6: WithValue 的查找性能如何？底层数据结构？
【理解】
WithValue 的底层是单向链表（不是 map）：
  type valueCtx struct {
      Context         // parent
      key, val any    // 一个节点只存一对 KV
  }

查找过程：
  ctx.Value(targetKey)
    -> 当前节点 key 匹配？返回 val
    -> 不匹配？向 parent 递归查找
    -> 直到根节点，返回 nil

时间复杂度：O(n)，n 是从当前节点到根的深度。
每次 WithValue 都创建新节点挂到链表尾部，不会修改已有节点（不可变）。

如果链路很长（十几层 middleware 都 WithValue），查找性能会退化。
所以不要滥用 WithValue，Key 的数量应该很少。
【回答】
WithValue 底层是单向链表，每个 valueCtx 节点只存一对 KV。
查找时从当前节点往 parent 方向逐层匹配 key，找到就返回，找不到一直到根节点返回 nil。
时间复杂度 O(n)，n 是链表深度。所以不要滥用 WithValue，key 数量应该很少，链路很长时性能会退化。

---
Q7: 为什么 context.WithCancel 要返回一个 cancel 函数？不调用会怎样？
【理解】
cancel 函数的作用：
  1. 关闭 done channel，通知所有子协程退出
  2. 把自己从 parent 的 children map 中移除（断开引用）
  3. 释放关联资源

不调用 cancel 的后果：
  - 子 context 永远留在 parent 的 children map 里（内存泄漏）
  - 关联的 goroutine 可能永远不退出（协程泄漏）
  - go vet 会报警告

所以标准写法是 defer cancel()：
  ctx, cancel := context.WithCancel(parent)
  defer cancel()  // 无论走哪个分支，最后都释放

即使 context 已经超时/被 parent 取消了，再调用 cancel 也没有副作用（幂等）。
【回答】
cancel 函数有两个作用：关闭 done channel 通知子协程退出，以及把自己从 parent 的 children map 中移除释放内存。
不调用的后果：子 context 永远留在 parent 的 children map 里造成内存泄漏，关联的 goroutine 也可能永远不退出。
标准写法是 defer cancel()，即使 context 已经超时了再调用 cancel 也没副作用，它是幂等的。

---
Q8: Context 使用有哪些最佳实践？
【理解】
1. context 作为函数第一个参数，变量名用 ctx
2. 不要把 context 存到结构体里（除非有特殊原因），应该显式传递
3. 不要传 nil context，不知道用什么就用 context.TODO()
4. WithValue 只传请求级别的横切数据，不传业务参数
5. cancel 函数拿到后立即 defer cancel()
6. 服务端收到请求后应该创建 context，客户端发起调用时应该传递 context
7. 长任务要定期检查 ctx.Done()，不要等任务做完才检查

反模式：
  - 把 context 存到结构体成员变量里
  - 用 WithValue 传业务参数
  - 忘记调用 cancel
  - 函数签名里没有 context 参数但内部自己创建 Background（切断了取消链路）
【回答】
核心实践：
context 作为函数第一个参数显式传递，不要存到结构体里。
拿到 cancel 函数后立即 defer cancel()，防止泄漏。
不要传 nil context，不知道用什么就用 context.TODO()。
WithValue 只传 TraceID 等横切数据，业务参数写到函数签名里。
长任务要定期检查 ctx.Done()，及时响应取消信号。
不要在函数内部偷偷创建 Background() 切断取消链路。

*/

// TestContextTimeout 超时控制演示
func TestContextTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel() // 无论是否超时，最后都调用 cancel 释放资源

	resChan := make(chan string)

	go func() {
		time.Sleep(2 * time.Second) // 模拟慢任务
		resChan <- "Success"
	}()

	select {
	case res := <-resChan:
		fmt.Println("得到结果:", res)
	case <-ctx.Done():
		fmt.Println("超时:", ctx.Err()) // context deadline exceeded
	}
}

// TestContextCancel 手动取消演示
func TestContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("子协程收到取消信号, 原因:", ctx.Err())
				return
			default:
				fmt.Println("工作中...")
				time.Sleep(200 * time.Millisecond)
			}
		}
	}(ctx)

	time.Sleep(500 * time.Millisecond)
	cancel() // 手动取消，子协程会收到信号退出
	time.Sleep(100 * time.Millisecond)
}

// TestContextValue WithValue 传值演示
func TestContextValue(t *testing.T) {
	type traceIDKey struct{} // 用空结构体作为 key，避免冲突

	ctx := context.WithValue(context.Background(), traceIDKey{}, "trace-abc-123")

	// 模拟中间件读取 traceID
	traceID := ctx.Value(traceIDKey{}).(string)
	fmt.Println("TraceID:", traceID)
}

// TestContextPropagation 取消传播演示：parent 取消，所有子 context 都会被取消
func TestContextPropagation(t *testing.T) {
	parent, cancelParent := context.WithCancel(context.Background())

	child1, cancelChild1 := context.WithCancel(parent)
	defer cancelChild1()
	child2, cancelChild2 := context.WithCancel(parent)
	defer cancelChild2()

	// 取消 parent，child1 和 child2 都会被取消
	cancelParent()

	fmt.Println("child1 err:", child1.Err()) // context canceled
	fmt.Println("child2 err:", child2.Err()) // context canceled
}

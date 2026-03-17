package base

import (
	"context"
	"fmt"
	"testing"
	"time"
)

/*
Context 主要解决了 Goroutine 树状管理 的问题。它有三大作用：退出通知（取消信号）、超时控制、以及元数据传递。

面试题：
1 context 常见函数？
 -WithCancel	手动取消	任务完成后，通知辅助协程退出。
 -WithTimeout	倒计时取消	调用第三方 API，规定 2 秒必须返回。
 -WithDeadline	定时取消	规定任务必须在 2026-03-09 18:00 前完成。
 -WithValue	    传值	   链路追踪里的 TraceID 或用户的 UserID。

2 Context 的底层原理？
 -核心结构：cancelCtx 结构体里有一个 done chan struct{}。
 -传播机制：
  -当调用 cancel() 时，它会关闭自己的 done channel。
  -此时，所有监听 <-ctx.Done() 的子协程都会收到零值，从而停止工作。
  -更关键的是，cancelCtx 内部维护了一个 children map。它会递归地调用所有子 Context 的 cancel 函数，确保整棵树都被清理。

3 能不能把所有参数都塞进 Context 的 Value 里传递？
 -绝对不行。WithValue 应该是 ‘非业务核心数据’ 的传递。
 -性能损耗：WithValue 查找数据是 O(n) 的。它是一个链表结构，找一个 key 要从叶子节点往根节点一层层找，找多了很慢。
 -类型安全：它是 interface{} 类型，拿出来要类型断言，容易出错。
 -不可见性：它隐藏了接口的显式依赖。
 -正确用法：只传 RequestID、UserToken 这种横向贯穿整个链路的元数据。

4 Context 是线程安全的吗？
 -在 Go 的并发哲学中，Context 被设计为不可变的。这意味着你可以在多个 Goroutine 之间传递同一个 ctx，而不需要加锁。

*/

func TestContextTimeout(t *testing.T) {
	// 1. 创建一个 1 秒超时的 Context
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel() // 养成好习惯：无论是否超时，最后都调用 cancel 释放资源

	resChan := make(chan string)

	go func() {
		// 模拟耗时操作
		time.Sleep(2 * time.Second)
		resChan <- "Success"
	}()

	select {
	case res := <-resChan:
		fmt.Println("得到结果:", res)
	case <-ctx.Done():
		// 2. 监听 Context 的结束信号
		fmt.Println("超时或取消:", ctx.Err()) // 输出: context deadline exceeded
	}
}

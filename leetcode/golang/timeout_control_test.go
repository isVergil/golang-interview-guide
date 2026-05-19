package golang

import (
	"context"
	"fmt"
	"testing"
	"time"
)

/*
题目：超时控制
使用 select + context 实现请求超时控制，避免 goroutine 永久阻塞。

核心思路：
  select 同时监听业务 channel 和超时信号，谁先到就走谁的分支。

关键点：
  1. context.WithTimeout 创建带超时的 context，到期自动 cancel
  2. select 中监听 ctx.Done()，超时后立即返回，不等慢任务
  3. 慢任务的 goroutine 仍在运行，需要在任务内部也检查 ctx 避免泄漏
  4. time.After 适合简单场景，context 适合需要传递取消信号的场景
*/

// TestTimeoutControl context 超时控制
func TestTimeoutControl(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	result := make(chan string, 1)

	// 模拟慢任务
	go func() {
		time.Sleep(200 * time.Millisecond) // 超过超时时间
		result <- "任务完成"
	}()

	select {
	case res := <-result:
		fmt.Println("成功:", res)
	case <-ctx.Done():
		fmt.Println("超时:", ctx.Err()) // context deadline exceeded
	}
}

// TestTimeoutWithRetry 超时重试模式
func TestTimeoutWithRetry(t *testing.T) {
	result := requestWithTimeout(50 * time.Millisecond)
	fmt.Println("结果:", result)
}

// requestWithTimeout 模拟带超时的请求，超时返回降级结果
func requestWithTimeout(timeout time.Duration) string {
	ch := make(chan string, 1)

	go func() {
		// 模拟耗时操作
		time.Sleep(30 * time.Millisecond)
		ch <- "正常响应"
	}()

	select {
	case res := <-ch:
		return res
	case <-time.After(timeout):
		return "超时降级"
	}
}

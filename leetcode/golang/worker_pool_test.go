package golang

import (
	"fmt"
	"sync"
	"testing"
)

/*
题目：Worker Pool（工作池）
固定数量的 worker goroutine 从任务队列取任务执行，控制并发度。

核心思路：
  1. 用 buffered channel 作为任务队列
  2. 启动固定数量的 worker，每个 worker 循环从 channel 取任务
  3. channel 关闭后 worker 自动退出

关键点：
  1. worker 数量 = 最大并发度，不会无限开 goroutine
  2. buffered channel 容量决定了能缓冲多少待处理任务
  3. 适用场景：限制对下游服务的并发请求数、CPU 密集型任务控制并行度
*/

// TestWorkerPool 固定 worker 数处理任务
func TestWorkerPool(t *testing.T) {
	const numWorkers = 3
	const numTasks = 20

	tasks := make(chan int, numTasks)
	var wg sync.WaitGroup

	// 启动 worker
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for task := range tasks {
				fmt.Printf("worker %d 处理任务 %d\n", id, task)
			}
		}(w)
	}

	// 发送任务
	for i := 1; i <= numTasks; i++ {
		tasks <- i
	}
	close(tasks)
	wg.Wait()
}

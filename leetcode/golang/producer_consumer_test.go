package golang

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

/*
题目：生产者消费者模型
多个生产者往 channel 发数据，多个消费者从 channel 取数据处理，生产完毕后优雅退出。

核心思路：
  1. 用 buffered channel 作为任务队列，解耦生产和消费速度
  2. 生产者全部完成后 close(channel)，消费者 range 自动退出
  3. 用 WaitGroup 分别等待生产者和消费者

关键点：
  1. close 只能由发送方调用，接收方 range 会在 close 后自动退出循环
  2. 多个生产者时，需要额外 WaitGroup 等所有生产者完成后再 close
  3. buffered channel 的容量决定了生产者能"跑在前面"多少
*/

// TestProducerConsumer 多生产者多消费者
func TestProducerConsumer(t *testing.T) {
	tasks := make(chan int, 10) // 缓冲队列
	var prodWg, consWg sync.WaitGroup

	// 3 个生产者
	for p := 1; p <= 3; p++ {
		prodWg.Add(1)
		go func(id int) {
			defer prodWg.Done()
			for i := 0; i < 5; i++ {
				task := id*100 + i
				tasks <- task
				fmt.Printf("生产者%d: 发送任务 %d\n", id, task)
			}
		}(p)
	}

	// 2 个消费者
	for c := 1; c <= 2; c++ {
		consWg.Add(1)
		go func(id int) {
			defer consWg.Done()
			for task := range tasks { // close 后自动退出
				fmt.Printf("消费者%d: 处理任务 %d\n", id, task)
				time.Sleep(10 * time.Millisecond)
			}
		}(c)
	}

	prodWg.Wait()  // 等所有生产者完成
	close(tasks)   // 关闭 channel，通知消费者
	consWg.Wait()  // 等所有消费者处理完
}

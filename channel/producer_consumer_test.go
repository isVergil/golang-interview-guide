package channel

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// 生产者-消费者模型（带停止信号）
// 实现一个经典的生产者-消费者模型。要求：生产者发完 100 个数据后通知消费者，消费者处理完所有数据后优雅退出。
func TestProducerConsumer(t *testing.T) {
	dataChan := make(chan int, 5)

	var wg sync.WaitGroup

	// 生产者
	go func() {
		for i := 1; i <= 10; i++ {
			dataChan <- i
			fmt.Printf("生产：%d\n", i)
		}
		close(dataChan)
	}()

	// 消费者
	wg.Add(1)
	go func() {
		defer wg.Done()
		for d := range dataChan {
			fmt.Printf("消费：%d\n", d)
			time.Sleep(time.Millisecond * 10)
		}
	}()

	wg.Wait()
}

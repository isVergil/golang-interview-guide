package golang

import (
	"fmt"
	"sync"
	"testing"
)

/*
题目：多路归并（Fan-in）
多个 channel 各自产生数据，将它们合并到一个 channel 统一消费。

核心思路：
  为每个输入 channel 启动一个 goroutine，都往同一个输出 channel 发送。
  所有输入 channel 关闭后，关闭输出 channel。

关键点：
  1. 用 WaitGroup 追踪所有输入 goroutine，全部完成后 close 输出 channel
  2. select 方案适合少量 channel；goroutine 方案适合动态数量的 channel
  3. 输出 channel 的消费者用 range 即可，close 后自动退出
*/

// TestFanIn 将多个 channel 合并为一个
func TestFanIn(t *testing.T) {
	// 模拟 3 个数据源
	sources := make([]<-chan int, 3)
	for i := 0; i < 3; i++ {
		ch := make(chan int)
		sources[i] = ch
		go func(id int, c chan<- int) {
			for j := 0; j < 5; j++ {
				c <- id*10 + j
			}
			close(c)
		}(i, ch)
	}

	// fan-in：合并到一个 channel
	merged := fanIn(sources...)

	for val := range merged {
		fmt.Println("收到:", val)
	}
}

// fanIn 将多个只读 channel 合并为一个
func fanIn(channels ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	for _, ch := range channels {
		wg.Add(1)
		go func(c <-chan int) {
			defer wg.Done()
			for v := range c {
				out <- v
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(out) // 所有输入关闭后关闭输出
	}()

	return out
}

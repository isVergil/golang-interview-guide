package channel

import (
	"fmt"
	"testing"
	"time"
)

// 多路归并 (Fan-in / 扇入)
// 有两个任务源（Channel），谁先产生数据就先处理谁，并将结果汇总到一个通道中。
func TestFanIn(t *testing.T) {
	ch1 := make(chan string)
	ch2 := make(chan string)

	// 模拟两个异步任务源
	go func() {
		for {
			ch1 <- "来自任务 A 的数据"
			time.Sleep(time.Second)
		}
	}()
	go func() {
		for {
			ch2 <- "来自任务 B 的数据"
			time.Sleep(time.Second * 2)
		}
	}()

	// 汇总处理
	for i := 0; i < 5; i++ {
		select {
		case msg1 := <-ch1:
			fmt.Println("处理:", msg1)
		case msg2 := <-ch2:
			fmt.Println("处理:", msg2)
		case <-time.After(time.Millisecond * 1500): // 容错控制
			fmt.Println("超时：等太久了")
		}
	}
}

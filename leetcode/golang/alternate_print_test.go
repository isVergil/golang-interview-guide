package golang

import (
	"fmt"
	"sync"
	"testing"
)

/*
题目：交替打印
使用 N 个 goroutine 交替打印指定内容，保证输出顺序严格轮转。

核心思路：
  用 N 个无缓冲 channel 形成"令牌环"，每个 goroutine 等待自己的 channel 收到信号后打印，
  然后把信号传给下一个 channel，实现严格轮转。

关键点：
  1. 无缓冲 channel 天然是同步的，收发必须配对，保证顺序
  2. 最后一轮最后一个 goroutine 不再发信号，避免死锁
  3. 主 goroutine 发出第一个信号启动链条
*/

// TestAlternatePrintAB 两个 goroutine 交替打印数字和字母：1A2B3C...26Z
func TestAlternatePrintAB(t *testing.T) {
	numCh := make(chan struct{})
	charCh := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(2)

	// 打印数字
	go func() {
		defer wg.Done()
		for i := 1; i <= 26; i++ {
			<-numCh // 等待令牌
			fmt.Print(i)
			charCh <- struct{}{} // 传给字母协程
		}
	}()

	// 打印字母
	go func() {
		defer wg.Done()
		for i := 0; i < 26; i++ {
			<-charCh // 等待令牌
			fmt.Printf("%c", 'A'+i)
			if i < 25 {
				numCh <- struct{}{} // 传回数字协程（最后一轮不发，避免死锁）
			}
		}
	}()

	numCh <- struct{}{} // 启动第一个信号
	wg.Wait()
	fmt.Println()
}

package golang

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

/*
题目：并发安全计数器
多个 goroutine 同时对计数器进行读写，保证最终结果正确。

三种实现方式对比：
  1. Mutex：通用方案，适合复合操作（读-改-写多个变量）
  2. atomic：无锁方案，适合单变量简单操作，性能最优
  3. channel：Go 风格方案，"通过通信共享内存"，适合需要串行化的场景

关键点：
  1. 不加保护的并发写会导致数据竞争（go test -race 可检测）
  2. atomic 只保证单个操作原子，不保证多操作组合的原子性
  3. Mutex 的 Lock/Unlock 之间是临界区，任意复杂逻辑都能保护
*/

// TestSafeCounterMutex Mutex 方案
func TestSafeCounterMutex(t *testing.T) {
	var mu sync.Mutex
	count := 0
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			count++
			mu.Unlock()
		}()
	}
	wg.Wait()
	fmt.Println("Mutex count:", count) // 1000
}

// TestSafeCounterAtomic atomic 方案（性能更优）
func TestSafeCounterAtomic(t *testing.T) {
	var count int64
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			atomic.AddInt64(&count, 1)
		}()
	}
	wg.Wait()
	fmt.Println("Atomic count:", count) // 1000
}

// TestSafeCounterChannel channel 方案（Go 风格）
func TestSafeCounterChannel(t *testing.T) {
	ch := make(chan int, 1000)
	var wg sync.WaitGroup

	// 1000 个 goroutine 发送增量
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ch <- 1
		}()
	}

	// 等发送完毕后关闭
	go func() {
		wg.Wait()
		close(ch)
	}()

	// 单个 goroutine 汇总（天然串行，无竞争）
	count := 0
	for v := range ch {
		count += v
	}
	fmt.Println("Channel count:", count) // 1000
}

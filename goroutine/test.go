package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done() // 确保在函数退出时调用Done

	fmt.Printf("Worker %d 开始工作\n", id)
	time.Sleep(time.Second * time.Duration(id)) // 模拟不同耗时的工作
	fmt.Printf("Worker %d 工作完成\n", id)
}

func main() {
	var wg sync.WaitGroup

	fmt.Println("=== 启动3个worker goroutine ===")

	// 方法1：逐个添加
	wg.Add(1)
	go worker(1, &wg)

	wg.Add(1)
	go worker(2, &wg)

	wg.Add(1)
	go worker(3, &wg)

	fmt.Println("主goroutine等待所有worker完成...")
	wg.Wait() // 阻塞直到所有worker调用Done()
	fmt.Println("所有worker都已完成！")

	fmt.Println("\n=== 批量添加示例 ===")

	// 方法2：批量添加
	workers := 5
	wg.Add(workers) // 一次性添加所有计数

	for i := 1; i <= workers; i++ {
		go func(workerID int) {
			defer wg.Done()
			fmt.Printf("批量Worker %d 执行中...\n", workerID)
			time.Sleep(time.Millisecond * 500)
		}(i)
	}

	wg.Wait()
	fmt.Println("批量worker全部完成！")
}

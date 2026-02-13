package main

import "fmt"

func main() {
	//无缓冲通道的写入和接收操作必须同时完成，否则会导致死锁。
	//不死锁
	ch := make(chan int) // 无缓冲通道
	go func() {
		fmt.Println("协程准备接收")
		val := <-ch // 接收操作
		fmt.Println("接收到值:", val)
	}()
	fmt.Println("主协程准备写入")
	ch <- 42 // 写入操作
	fmt.Println("写入完成")

	//死锁
	ch1 := make(chan int) // 创建一个管道ch
	ch1 <- 1              // 向管道ch中发送数据v.
	_ = <-ch1             // 从管道中读取数据存储到变量v
	close(ch1)
}

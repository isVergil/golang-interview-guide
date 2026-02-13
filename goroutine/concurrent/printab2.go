package main

import (
	"fmt"
	"sync"
)

// 2 个协程打印，需要注意终止条件
func main() {
	fun1()
	fun2()
}

// 无缓冲的
func fun1() {
	a := make(chan struct{})
	b := make(chan struct{})
	wg := sync.WaitGroup{}

	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 'a'; i <= 'z'; i++ {
			<-a
			fmt.Printf("%c\n", i)
			b <- struct{}{}
		}
	}()

	go func() {
		defer wg.Done()
		for i := 'a'; i <= 'z'; i++ {
			<-b
			fmt.Printf("%d\n", i)
			// 最后一次循环，上面协程已退出，a 不再接收此信号，便会引发死锁
			if i < 'z' {
				a <- struct{}{}
			}
		}
	}()

	a <- struct{}{}

	wg.Wait()
}

// 有缓冲的
// 有缓冲通道允许异步通信，发送方和接收方不需要同时就绪。
func fun2() {
	a := make(chan struct{}, 1)
	b := make(chan struct{}, 1)
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 'a'; i <= 'z'; i++ {
			<-a
			fmt.Printf("%c\n", i)
			b <- struct{}{}
		}
	}()

	go func() {
		defer wg.Done()
		for i := 'a'; i <= 'z'; i++ {
			<-b
			fmt.Printf("%d\n", i)
			a <- struct{}{}
		}
	}()

	a <- struct{}{}

	wg.Wait()
}

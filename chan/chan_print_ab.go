package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	printNoBuffer()
	//printNoBuffer()
}

// 有缓冲的 chan 交替打印 AB
func printBuffer() {
	lock := make(chan bool, 1)
	go func() {
		for {
			fmt.Println("A")
			lock <- true
		}
	}()

	go func() {
		for {
			<-lock
			fmt.Println("B")
		}
	}()

	time.Sleep(time.Second)
}

// 无缓冲的 chan 交替打印 AB
func printNoBuffer() {
	lock := make(chan bool)
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			fmt.Println("A")
			lock <- true
			time.Sleep(100 * time.Millisecond)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			<-lock
			fmt.Println("B")
		}
	}()

	wg.Wait()
}

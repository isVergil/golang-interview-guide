package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	fmt.Println(runtime.NumCPU())

	for i := 0; i < 100000; i++ {
		go time.Sleep(time.Second * 10)
	}
	time.Sleep(time.Second * 5)
	fmt.Println(runtime.NumGoroutine())
	time.Sleep(time.Second * 100)
}

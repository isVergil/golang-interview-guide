package main

import (
	"fmt"
	"time"
)

// 只需要两个 chan 但是主协程得等待 这里用 time.Sleep(time.Second * 5)
func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	go func() {
		i := 1
		for {
			if i > 26 {
				return
			}
			<-ch1
			fmt.Println(i)
			i += 1
			ch2 <- i
		}

	}()

	go func() {
		str := 'A'
		for {
			i := <-ch2
			fmt.Println(string(str))
			str += 1
			ch1 <- i
		}
	}()

	ch1 <- 1
	time.Sleep(time.Second * 3)
}

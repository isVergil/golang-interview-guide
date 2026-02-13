package main

import (
	"fmt"
	"sync"
)

const N = 10

func main() {
	var wg sync.WaitGroup
	ch1, ch2, ch3 := make(chan struct{}), make(chan struct{}), make(chan struct{})
	wg.Add(3)

	go func(s string) {
		defer wg.Done()
		for i := 0; i < N; i++ {
			<-ch1
			fmt.Println(s)
			ch2 <- struct{}{}
		}
	}("A")

	go func(s string) {
		defer wg.Done()
		for i := 0; i < N; i++ {
			<-ch2
			fmt.Println(s)
			ch3 <- struct{}{}
		}
	}("B")

	go func(s string) {
		defer wg.Done()
		for i := 0; i < N; i++ {
			<-ch3
			fmt.Println(s)
			if i < N-1 {
				ch1 <- struct{}{}
			}
		}
	}("C")

	ch1 <- struct{}{}
	wg.Wait()
}

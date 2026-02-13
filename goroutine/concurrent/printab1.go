package main

import (
	"fmt"
	"sync"
)

func main() {
	var mu sync.Mutex
	cond := sync.NewCond(&mu)

	var trun int
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 'a'; i <= 'z'; i++ {
			mu.Lock()
			for trun != 1 {
				cond.Wait()
			}
			fmt.Printf("%c\n", i)
			trun = 0
			cond.Broadcast()
			mu.Unlock()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 'a'; i <= 'z'; i++ {
			mu.Lock()
			for trun != 0 {
				cond.Wait()
			}
			fmt.Printf("%d\n", i)
			trun = 1
			cond.Broadcast()
			mu.Unlock()
		}
	}()

	wg.Wait()
}

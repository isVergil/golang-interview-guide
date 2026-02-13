package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {

	var sum int32 = 0
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			atomic.AddInt32(&sum, 1)
		}()
	}
	wg.Wait()
	fmt.Printf("sum is %d\n", sum)
}

package main

import "fmt"

func main() {
	ch := make(chan int, 10)
	for i := 1; i <= 10; i++ {
		go func(i int) {
			fmt.Printf("send num: %d\n", i)
			ch <- i
		}(i)
	}

	for i := 1; i <= 10; i++ {
		res := <-ch
		fmt.Printf("get num: %d\n", res)
	}

	fmt.Println("end")
}

package main

import (
	"fmt"
)

func main() {
	ch1 := make(chan bool)
	ch2 := make(chan bool)
	done := make(chan bool)

	go func() {
		for i := 1; i < 27; i++ {
			<-ch1
			fmt.Println(i)
			ch2 <- true
		}
	}()

	go func() {
		for i := 'A'; i <= 'Z'; i++ {
			<-ch2
			fmt.Println(string(i))
			ch1 <- true
		}
		done <- true
	}()
	ch1 <- true
	<-done
}

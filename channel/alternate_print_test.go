package channel

import (
	"fmt"
	"sync"
	"testing"
)

// TestPrintABNoBuffer 交替打印 AB
func TestPrintABNoBuffer(t *testing.T) {
	numChan := make(chan struct{})
	charChan := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 1; i <= 26; i++ {
			<-numChan
			fmt.Print(i)
			charChan <- struct{}{}
		}
	}()

	go func() {
		defer wg.Done()
		for i := 1; i <= 26; i++ {
			<-charChan
			fmt.Printf("%c", 'A'+i-1)
			if i < 26 {
				numChan <- struct{}{}
			}
		}
	}()

	numChan <- struct{}{}
	wg.Wait()
	fmt.Println("\n执行完毕")
}

// TestPrintABNoBuffer 交替打印 ABC
func TestPrintABCNoBuffer(t *testing.T) {
	chanA := make(chan struct{})
	chanB := make(chan struct{})
	chanC := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(3)

	max := 10

	go func() {
		defer wg.Done()
		for i := 1; i <= max; i++ {
			<-chanA
			fmt.Print("A")
			chanB <- struct{}{}
		}
	}()

	go func() {
		defer wg.Done()
		for i := 1; i <= max; i++ {
			<-chanB
			fmt.Print("B")
			chanC <- struct{}{}
		}
	}()

	go func() {
		defer wg.Done()
		for i := 1; i <= max; i++ {
			<-chanC
			fmt.Print("C")
			if i < max {
				chanA <- struct{}{}
			}
		}
	}()

	chanA <- struct{}{}

	wg.Wait()
	fmt.Println("\n执行完毕")
}

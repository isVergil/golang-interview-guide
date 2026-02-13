package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	startTime := time.Now()
	dumpPrint()
	fmt.Println("dumpPrint time cost: ", time.Since(startTime))

	startTime = time.Now()
	joinPrint()
	fmt.Println("joinPrint time cost: ", time.Since(startTime))

	// fmt.Println("args name: %s", os.Args[0])
	// fmt.Printf("args counts: %d\n", len(os.Args))
	// for idx, arg := range os.Args[1:] {
	// 	fmt.Printf("args %d: %s\n", idx, arg)
	// }
}

func dumpPrint() {
	s, sep := "", ""
	for i := 1; i < len(os.Args); i++ {
		s += sep + os.Args[i]
		sep = " "
	}
	fmt.Println(s)
}

func joinPrint() {
	fmt.Println(strings.Join(os.Args[1:], " "))
}

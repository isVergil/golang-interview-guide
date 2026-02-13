package main

import (
	"fmt"
)

/*
1.若有多个 case 满足,随机执行
2.每次 select 只判断条件是否满足 case，执行结果不影响其他 case 的判断
3.case <- chan: 语句用于从通道 chan 中读取数据，数据会被丢弃
4.case data := <-chan: 读取数据并赋值给变量 data
*/
func main() {
	chA := make(chan string, 1)
	chB := make(chan string, 1)
	chC := make(chan string, 1)
	chA <- "A"
	for i := 0; i < 10; i++ {
		select {
		case <-chA:
			chB <- "B"
			fmt.Println("A")
		case <-chB:
			chC <- "C"
			fmt.Println("B")
		case <-chC:
			chA <- "A"
			fmt.Println("C")
		}
	}
}

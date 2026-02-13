package main

import "fmt"

// 结构体
type User struct {
	Name string
	Age  int
}

func main() {
	// new(type) 返回 *type
	p1 := new(int)    // *int, 值为0
	p2 := new(string) // *string, 值为""
	p3 := new(bool)   // *bool, 值为false
	p4 := new(User)   // *User, 所有字段为零值

	fmt.Printf("p1: %T = %v\n", p1, p1)
	fmt.Printf("p2: %T = %v\n", p2, p2)
	fmt.Printf("p3: %T = %v\n", p3, p3)
	fmt.Printf("p4: %T = %v\n", p4, p4)

	// slice - 需要指定长度（和可选容量）
	s := make([]int, 5)      // []int, 值为[0 0 0 0 0]
	s2 := make([]int, 3, 10) // 长度3，容量10

	// map - 可指定初始容量
	m := make(map[string]int)       // map[string]int, 空map
	m2 := make(map[string]int, 100) // 预设容量

	// channel - 可指定缓冲区大小
	ch := make(chan int)      // 无缓冲通道
	ch2 := make(chan int, 10) // 缓冲大小为10
	fmt.Printf("s: %T = %#v\n", s, s)
	fmt.Printf("s2: %T = %#v\n", s2, s2)
	fmt.Printf("m: %T = %#v\n", m, m)
	fmt.Printf("m2: %T = %#v\n", m2, m2)
	fmt.Printf("ch: %T = %#v\n", ch, ch)
	fmt.Printf("ch2: %T = %#v\n", ch2, ch2)

}

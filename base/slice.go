package main

import "fmt"

func main() {
	// 1. 初始化
	s := make([]int, 0, 3) // len=0, cap=3
	fmt.Printf("Initial: addr=%p, len=%d, cap=%d\n", s, len(s), cap(s))

	// 2. 扩容观察
	//oldPtr := &s
	for i := 0; i < 4; i++ {
		s = append(s, i)
	}
	// 此时容量翻倍，底层数组地址发生变化
	fmt.Printf("After Append: addr=%p, len=%d, cap=%d\n", s, len(s), cap(s))

	// 3. 截取与共享
	sub := s[1:3]
	sub[0] = 999
	fmt.Println("Original after sub change:", s[1]) // 输出 999，共享内存

	// 4. 解决内存泄露
	huge := make([]int, 1000)
	small := make([]int, 2)
	copy(small, huge[0:2]) // small 拥有独立数组，huge 可被回收
}

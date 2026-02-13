package main

import "fmt"

const (
	A = iota // 0
	B        // 1
	C = "c"  // c
	D        // c，与上一⾏行相同。
	E = iota // 4，显式恢复。注意计数包含了 C、D 两⾏行。
	F        // 5
)

// 作用域限制‌：iota 仅在 const 声明块内有效，离开 const 块后无法使用
// 跨块重置‌：遇到下一个 const 关键字时，iota 的值会被重新置为 0
const (
	G = iota + 1 // 1
	H            // 2
)

func main() {
	fmt.Println("\n格式化输出：")
	fmt.Printf("A: %T = %v\n", A, A)
	fmt.Printf("B: %T = %v\n", B, B)
	fmt.Printf("C: %T = %v\n", C, C)
	fmt.Printf("D: %T = %v\n", D, D)
	fmt.Printf("E: %T = %v\n", E, E)
	fmt.Printf("F: %T = %v\n", F, F)
	fmt.Printf("G: %T = %v\n", G, G)
	fmt.Printf("H: %T = %v\n", H, H)
}

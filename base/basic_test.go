package base

import (
	"fmt"
	"testing"
	"unicode/utf8"
)

/*
面试题：
1 new、make 区别？
 -new 分配内存并置零，返回指针；
 -make 只用于 slice/map/chan，初始化并返回引用。

2 字符串遍历有什么要注意的？
 -字符串不可变性：Go 字符串底层是只读字节数组，修改需转为 []byte 或 []rune
 -for i 遍历的是“字节(byte): uint8
 -for range 遍历的是“字符(rune/unicode) : int32

3 如何修改字符串？
 -转为 []byte 或 []rune

4 iota 是啥？有什么用法？
 -编译器级别的常量行计数器，是编译期的行为，它不占用运行时的 CPU 计算资源，性能极高
 -作用域：只在 const 块中有效，遇到下一个 const 关键字重置为 0。
 -应用：多用于枚举、字节单位换算（KB/MB/GB）、位掩码（Bitmask）。

*/

// TestNewVsMake new 和 make 的区别
func TestNewVsMake(t *testing.T) {
	t.Log("--- new vs make ---")

	// 1. new: 分配空间，所有位清零，返回指针
	p := new(int)
	fmt.Printf("new(int) 类型: %T, 值: %v\n", p, *p) // *int, 0

	// 2. make: 专门用于内置引用类型 (slice, map, chan)
	// 如果用 new 弄一个 map，直接写值会 panic，因为没初始化内部结构
	m := make(map[string]int)
	m["key"] = 1
	fmt.Printf("make(map) 结果: %v\n", m)
}

// TestStringIter 字符串遍历
func TestStringIter(t *testing.T) {
	t.Log("--- 字符串遍历区别 ---")
	s := "Go语言" // "语言" 两个字各占 3 个字节

	// 1. 普通 for 循环 (byte 遍历)
	fmt.Print("for i 遍历 (byte): ")
	for i := 0; i < len(s); i++ {
		fmt.Printf("%x ", s[i]) // 打印的是字节的十六进制
	}
	fmt.Println("\n结论：for i 可能会把中文截断，它是按字节读的。")

	// 2. for range 遍历 (rune 遍历)
	fmt.Print("for range 遍历 (rune): ")
	for _, r := range s {
		fmt.Printf("%c ", r) // 能够正确识别中文
	}
	fmt.Println("\n结论：for range 会自动处理 UTF-8 解码，按字符读。")

	// 3. 统计长度
	fmt.Printf("len(s): %d (字节长度)\n", len(s))
	fmt.Printf("RuneCount: %d (字符数量)\n", utf8.RuneCountInString(s))
}

// TestStringImmutable 细节题：如何修改字符串
func TestStringImmutable(t *testing.T) {
	s := "hello"
	// s[0] = 'H' // 编译报错：cannot assign to s[0]

	// 方案 A：转为 []byte (适用于纯 ASCII)
	b := []byte(s)
	b[0] = 'H'
	s2 := string(b)

	// 方案 B：转为 []rune (适用于含中文)
	r := []rune("中文")
	r[0] = '日'
	s3 := string(r)

	fmt.Println("修改后:", s2, s3)
}

// TestIota 跳值与插值
func TestIota(t *testing.T) {
	const (
		A = iota  // 0
		B         // 1 (隐式重复上行的表达式)
		C = "Gap" // C="Gap", 此时 iota 依然在计数，为 2
		D         // D="Gap", 此时 iota 为 3
		E = iota  // E=4 (显式找回计数器)
	)
	fmt.Printf("A:%d, B:%d, C:%v, D:%v, E:%d\n", A, B, C, D, E)
}

// TestIotaBitmask 位运算定义
func TestIotaBitmask(t *testing.T) {
	// 问：如何定义权限位或存储单位？
	// 答：利用左移运算 (<<) 配合 iota。
	const (
		_  = iota             // 忽略 0
		KB = 1 << (10 * iota) // 1 << (10*1) = 1024
		MB = 1 << (10 * iota) // 1 << (10*2) = 1048576
		GB = 1 << (10 * iota)
	)
	fmt.Printf("KB: %d, MB: %d, GB: %d\n", KB, MB, GB)
}

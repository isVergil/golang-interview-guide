package basics

import (
	"fmt"
	"testing"
)

/*
Q1: new 和 make 的区别？
Q2: iota 是什么？
Q3: Go 是值传递还是引用传递？
Q4: Go 的零值有什么设计意义？好处？
Q5: 什么是 zerobase？

---
Q1: new 和 make 的区别？
【理解】
new(T) 分配零值内存，返回 *T；make 只用于 slice/map/chan，初始化内部结构，返回 T。
slice/map/chan 底层有复杂结构（hmap、hchan），光清零不够，必须 make 初始化才能用。
实际开发中 new 很少用，因为 &T{} 既能取指针又能赋初始值更灵活。
new 偶尔在需要基本类型指针时用（如 protobuf 的可选字段 *int32），大部分场景被 &T{} 和 make 覆盖。
【回答】
两者都会分配内存，区别在于初始化程度和适用类型。
new 给任何类型分配内存并置零，返回指针 *T；make 专门给 slice/map/channel 用，返回 T 值本身，make 不仅分配内存还会初始化底层结构。
比如 new 一个 map 出来不能直接写值会 panic，因为只是清零了，哈希桶没建好；make 出来的 map 哈希桶都建好了，可以直接用。
另外实际开发中 new 很少用，基本被 &T{} 替代了。

---
Q2: iota 是什么？
【理解】
编译期常量行计数器，const 块内从 0 开始每行 +1，遇到新 const 重置。
用途：枚举、位掩码（1<<iota）、字节单位换算。
注意：中间插其他赋值，iota 照样递增，它只认行号。
【回答】
iota 是 Go 的常量计数器，在 const 块里从 0 开始每行自动加 1，是编译期行为，零运行开销。
常用来做枚举和位掩码。有两个细节：一是后续行会继承第一行的整个表达式，只有 iota 在变；二是中间即使插了别的赋值，iota 也不会停，它只认行号不认赋值。"
补充：iota 没有运行时开销，它是编译器在处理 const 块时维护的一个行计数器。编译器遍历每一行时把 iota 替换为当前计数值，然后求值得到最终常量。编译完成后 iota 就不存在了，二进制里只有算好的字面量数字。

---
Q3: Go 是值传递还是引用传递？
【理解】
只有值传递。slice/map/channel 改了外部能看到，是因为它们本身存的就是指针/header，
拷贝的副本里的指针还是指向同一块底层数据。但 append 扩容后外部看不到。
【回答】
Go 只有值传递，没有引用传递。函数传参永远是拷贝一份副本。
slice 传进去改元素外面能看到，是因为拷贝的 header 里的指针还指向同一个底层数组。
但如果 append 触发扩容，产生新数组，外部就看不到了。map 和 channel 同理，底层就是指针。
补充：对于切片来说，值传递拷贝了 header，两份 header 各自独立。扩容只改了函数内那份的 Data 指针，外面那份纹丝不动。

---
Q4: Go 的零值有什么设计意义？好处？
【理解】
数值=0，布尔=false，字符串=""，引用类型=nil。
好处：避免垃圾值（安全）、很多类型零值直接可用（如 sync.Mutex、bytes.Buffer）。
C 语言中局部变量声明后不会自动清零，内存里是上一次残留的随机数据
【回答】
Go 变量声明后没赋值会自动置为零值：数字是 0，布尔是 false，字符串是空串，指针/slice/map/channel 是 nil。
设计意义有两点：一是安全，不会像 C 那样读到未初始化的垃圾数据；二是很多标准库类型零值直接可用，比如 sync.Mutex 不需要 init 就能 Lock，bytes.Buffer 不需要 make 就能 Write，减少了构造函数的必要性。
好处：减少了样板代码。很多语言需要显式构造才能用。这样做的好处是尽量让零值有意义，用户声明完就能直接用，不需要额外的初始化步骤。

---
Q5: 什么是 zerobase？
【理解】
Go 语言规范要求：两个不同变量的地址必须不同（除非编译器能证明它们不会被取地址比较）。但 0 字节对象没有实际内容，分配真实内存是浪费。zerobase 是一个折中方案：
- 所有零字节对象都指向同一地址（节省内存）
- 这个地址是合法的、非 nil 的（满足"指针不为 nil"的语义）
- 不能解引用（因为没有真实存储空间）
【回答】
zerobase 是 Go runtime 中一个全局的、固定地址的占位符变量，专门用于所有大小为 0 的内存分配请求。
当你试图分配一个 0 字节的对象时，Go 不会真的分配内存，而是统一返回 &zerobase 这个地址，用于节约内存。
*/

// TestNewVsMake new 和 make 的区别
func TestNewVsMake(t *testing.T) {
	p := new(int)
	fmt.Printf("new(int): %T, val=%v\n", p, *p)

	m := make(map[string]int)
	m["key"] = 1
	fmt.Printf("make(map): %v\n", m)

	// new 出来的 map 不能直接用：
	// m2 := new(map[string]int)
	// (*m2)["key"] = 1  // panic: assignment to entry in nil map
}

// TestIota 跳值与位掩码
func TestIota(t *testing.T) {
	const (
		A = iota  // 0
		B         // 1
		C = "Gap" // iota=2，但值是 "Gap"
		D         // iota=3，值依然是 "Gap"
		E = iota  // 4，显式找回计数器
	)
	fmt.Printf("A:%d B:%d C:%v D:%v E:%d\n", A, B, C, D, E)

	const (
		_  = iota
		KB = 1 << (10 * iota) // 1024
		MB = 1 << (10 * iota) // 1048576
		GB = 1 << (10 * iota)
	)
	fmt.Printf("KB:%d MB:%d GB:%d\n", KB, MB, GB)
}

// TestValuePass 值传递验证
func TestValuePass(t *testing.T) {
	a := 10
	changeInt(a)
	fmt.Println("基本类型:", a) // 10，不变

	s := []int{1, 2, 3}
	changeSlice(s)
	fmt.Println("改元素:", s) // [99 2 3]，共享底层数组

	s2 := []int{1, 2}
	appendSlice(s2)
	fmt.Println("append:", s2) // [1 2]，扩容后外部看不到
}

func changeInt(x int)     { x = 999 }
func changeSlice(s []int) { s[0] = 99 }
func appendSlice(s []int) { s = append(s, 100) }

// TestZeroValue 零值可用性
func TestZeroValue(t *testing.T) {
	var s []int
	s = append(s, 1, 2, 3)
	fmt.Println("nil slice append:", s) // 直接可用

	// 读返回 0 值，写 panic
	var m map[string]int
	fmt.Println("nil map 读:", m["key"]) // 0，不panic
	// m["key"] = 1  // panic: assignment to entry in nil map
}

package base

import (
	"fmt"
	"strings"
	"testing"
)

/*
闭包：普通函数执行完就“忘事”了，局部变量会被回收；但闭包函数会把外层函数的作用域‘打包’带走。只要闭包还没被销毁，那些被捕获的外部变量就会一直活着。
定义：闭包 = 函数 + 引用环境。函数“捕获”了外层变量，并与之绑定。
底层：闭包是一个结构体，包含函数指针和受引用的变量指针。
逃逸分析：闭包引用的局部变量会从“栈”逃逸到“堆”，确保变量在函数返回后依然存活。
风险：容易引起并发 Bug（循环变量共享）和内存泄漏（长生命周期闭包持有大对象）。

面试题：
1 到底什么是闭包？
 -闭包就是一个函数加上它引用的外部变量。
 -你可以把它想象成一个自带‘私人背包’的函数。普通的函数执行完，它里面的局部变量就销毁了；但闭包不同，它会把定义时抓到的那些变量一直‘背’在身上。只要这个闭包函数还在，它‘背包’里的变量就一直活着，哪怕外层函数已经退出了。

2 闭包引用的变量存在哪？（底层原理）？
 -这涉及到 Go 的逃逸分析。
 -正常情况下，函数局部变量是在‘栈’上的，函数结束就回收。但如果一个变量被闭包引用了，编译器发现这个变量在函数返回后还要用，就会把它‘逃逸’到堆（Heap）上。
 -所以，闭包底层其实是一个结构体，里面存了一个函数指针，还有一堆指向这些堆变量的指针。

3 在循环里创建闭包打印循环变量 i，为什么结果全是最大值？
 -因为闭包捕获的是变量的地址，而不是当时的值。
 -在循环里，所有的闭包其实都盯着同一个 i 的地址。循环跑得飞快，等这些闭包真正开始执行的时候，循环早就结束了，此时 i 已经被加到了最大值。所以它们去地址里一跳，看到的就全是同一个最终值。
 -解决方法：在 Go 1.22 之前，我们需要在循环体里写个 i := i 这种‘变量影子’来制造副本；或者直接把 i 当作参数传进闭包函数，利用值拷贝来解决。

补充：
1 在实际开发中，我常用闭包来做函数工厂（比如后缀检查器）、装饰器（比如给函数加耗时统计）或者回调函数。它是实现‘逻辑封装’和‘状态保持’的神器，但前提是要避开循环变量共享和长时间持有大对象这两个坑。
2 如果闭包捕获了很大的 Slice 或者 Map，而这个闭包又被长期持有（比如存在全局 Map 里），由于逃逸分析的作用，这些大对象会一直留在堆上无法被 GC。这时候我会考虑手动置空或者传参替代闭包，来优化内存性能。
*/

// 闭包能捕获外部变量的核心机制在于它‌保存的是变量的引用（内存地址），而不是变量值的拷贝
// 当闭包形成时，Go语言会在堆上创建一个"引用容器"（funcval结构体），其中存储了被捕获变量的内存地址
// 模拟编译器对闭包的底层实现（伪代码）
//
//	type ClosureStruct struct {
//		F     uintptr // 函数指针
//		addrI *int    // 捕获的变量 i 的地址 (逃逸到堆)
//	}
//
// 当你写 f := func() { fmt.Println(i) } 时
// 编译器底层其实在做：
// c := &ClosureStruct{ F: funcAddr, addrI: &i }
func TestClosure(t *testing.T) {
	next := counter()
	fmt.Println(next()) // 输出: 1
	fmt.Println(next()) // 输出: 2
	fmt.Println(next()) // 输出: 3

	addLog := makeSuffix(".log")
	fmt.Println(addLog("access"))    // access.log
	fmt.Println(addLog("error.log")) // error.log

	closureloop()
}

// 1 状态保持与计数器
func counter() func() int {
	// 累加器 count 逃逸到了堆上
	count := 0
	return func() int {
		count++
		return count
	}
}

// 2 函数工厂模式‌ 通过闭包可以创建具有特定配置的函数
// 使用闭包工厂可以预先配置好一个逻辑，然后到处复用，而不需要每次都传那个重复的后缀。
func makeSuffix(suffix string) func(string) string {
	return func(name string) string {
		if !strings.HasSuffix(name, suffix) {
			return name + suffix
		}
		return name
	}
}

// 3 for 循环闭包
func closureloop() {
	var funcs []func()

	for i := 0; i < 3; i++ {
		// 错误写法：直接引用 i
		// 正确写法：i := i
		funcs = append(funcs, func() {
			fmt.Printf("Incorrect i: %d\n", i)
		})
	}

	for _, f := range funcs {
		f() // 输出均为 3
	}
}

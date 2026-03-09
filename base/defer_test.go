package base

import (
	"fmt"
	"testing"
)

/*
核心用途：资源释放（锁、文件、连接）、异常处理（recover）。
执行顺序：栈结构，后进先出。
重点陷阱：defer 与 return 的先后顺序。执行顺序是：设置返回值 -> 执行 defer -> 彻底退出。
性能：Go 1.14 后引入了开放编码（open-coded defers），开销已经非常小了。

面试题：
1 defer、return、返回值到底谁先谁后？
 -return 并不是一个原子操作。它先给返回值赋值，然后去跑 defer。
 -如果你的返回值是匿名的（比如 func() int），defer 即使改了代码里的变量也动不了返回值；
 -但如果返回值是命名的（比如 func() (res int)），defer 直接操作的就是那个 res 变量，所以能改掉最终结果。”

2 什么是参数预计算、还有 defer 闭包结合？
 -存快照、翻旧账的区别，见 TestDeferTriple 函数

*/

// TestDeferOrder 基础题：执行顺序
func TestDeferOrder(t *testing.T) {
	t.Log("--- LIFO 顺序展示 ---")
	defer fmt.Println("第一顺位定义（最后执行）")
	defer fmt.Println("第二顺位定义")
	defer fmt.Println("第三顺位定义（最先执行）")
	// 输出顺序：三 -> 二 -> 一
}

// TestDeferValue 重点：预计算参数
func TestDeferValue(t *testing.T) {
	t.Log("--- 参数预计算 ---")
	i := 1
	defer fmt.Println("defer 里的 i:", i) // 这里的 i 在此刻已经定死成 1 了

	i++
	fmt.Println("main 里的 i:", i)
	// 输出：main 里的 i: 2，然后才是 defer 里的 i: 1
}

// TestDeferTriple 参数机制
func TestDeferTriple(t *testing.T) {
	x := "初始值"

	// 1. 指定函数传参 (快照)
	defer fmt.Println("1.快照派:", x)

	// 2. 匿名函数 (翻旧账)
	defer func() {
		fmt.Println("2.翻旧账派:", x)
	}()

	// 3. 匿名函数传参 (伪装成闭包的快照)
	defer func(val string) {
		fmt.Println("3.伪装派:", val)
	}(x)

	x = "修改值"
}

// TestDeferReturn 必杀题：defer 如何修改返回值？
func TestDeferReturn(t *testing.T) {
	t.Log("--- 修改返回值陷阱 ---")
	fmt.Println("无名返回值结果:", returnV1()) // 输出 5
	fmt.Println("有名返回值结果:", returnV2()) // 输出 6
}

// 场景 A：无名返回值
func returnV1() int {
	i := 5
	defer func() {
		i++ // 修改的是局部变量 i，不是返回值
	}()
	return i // 返回值已经在这一步被定下来是 5 了
}

// 场景 B：有名返回值
func returnV2() (result int) {
	result = 5
	defer func() {
		result++ // 修改的就是有名变量 result，它直接影响最终结果
	}()
	return result
}

// 配合 recover 处理 panic
func TestDeferWithRecover() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()

	a, b := 1, 0
	fmt.Println("result: ", a/b)
}

package main

import "fmt"

func main() {
	// 1 defer语句的参数在注册时就固定了，后续的变量修改不会影响已注册的defer调用‌
	//   defer表达式中变量的值在defer表达式被定义时就已经明确
	a()

	// 2 如果需要在defer中获取变量的最新值，应该使用闭包函数
	b()

	// 3 调用顺序是按照先进后出的方式 栈stack 结构
	c()

	// 4 表达式中可以修改函数中的命名返回值

	// 5 引用类型在defer语句中传递的是‌底层数据的引用地址‌，而不是数据本身的拷贝
	numbers := []int{1, 2, 3}
	defer fmt.Println("切片:", numbers) // 记录当前切片引用

	numbers[0] = 100 // 修改底层数组
	fmt.Println("修改切片元素后:", numbers)

	withRecover()
}

func a() {
	i := 0
	defer fmt.Println(i)
	i++
	return
}

func b() {
	i := 0
	defer func() {
		fmt.Println(i)
	}()
	i++
	return
}

func c() {
	defer fmt.Print(1)
	defer fmt.Print(2)
	defer fmt.Print(3)
	defer fmt.Print(4)
}

func d() (i int) {
	defer func() { i++ }()
	return 1
}

// 配合 recover 处理 panic
func withRecover() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()

	a, b := 1, 0
	fmt.Println("result: ", a/b)

}

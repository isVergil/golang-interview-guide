package main

import (
	"fmt"
	"strings"
)

func main() {
	next := counter()
	fmt.Println(next()) // 输出: 1
	fmt.Println(next()) // 输出: 2
	fmt.Println(next()) // 输出: 3

	addLog := makeSuffix(".log")
	fmt.Println(addLog("access"))    // access.log
	fmt.Println(addLog("error.log")) // error.log
}

// 闭包能捕获外部变量的核心机制在于它‌保存的是变量的引用（内存地址），而不是变量值的拷贝
// 当闭包形成时，Go语言会在堆上创建一个"引用容器"（funcval结构体），其中存储了被捕获变量的内存地址

// 1 状态保持与计数器
func counter() func() int {
	count := 0
	return func() int {
		count++
		return count
	}
}

// 2 函数工厂模式‌ 通过闭包可以创建具有特定配置的函数
func makeSuffix(suffix string) func(string) string {
	return func(name string) string {
		if !strings.HasSuffix(name, suffix) {
			return name + suffix
		}
		return name
	}
}

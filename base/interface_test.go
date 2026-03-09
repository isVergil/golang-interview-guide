package base

import (
	"fmt"
	"testing"
	"unsafe"
)

/*
接口 interface
底层结构是 Go 语言实现鸭子类型（Duck Typing）和运行时多态的基石。在 Go 的 runtime 包中，接口被分为两种：eface（空接口）和 iface（非空接口）。
type eface struct {
    _type *_type         // 指向具体的类型信息（如 int, string, struct）
    data  unsafe.Pointer // 指向实际数据的指针
}

type iface struct {
    tab  *itab          // 核心：存放接口与类型的映射关系
    data unsafe.Pointer // 指向实际数据的指针
}

itab是接口最精妙的地方，它不仅记录了类型，还缓存了方法表。
type itab struct {
    inter *interfacetype // 接口本身的定义（方法名、包名等）
    _type *_type         // 实际变量的具体类型
    hash  uint32         // 拷贝自 _type.hash，用于快速类型断言
    _     [4]byte
    fun   [1]uintptr     // 函数指针数组，指向具体类型实现的方法地址
}


面试题：
1 interface 是什么？底层结构了解吗？
 -接口在底层是一个双指针结构体。按使用情况分为 eface 或者 iface
 -空接口 interface{}：eface 包含两个指针：一个指向类型信息 _type，一个指向具体数据 data。
 -非空接口底层是 iface：用于定义了方法的接口。包含两个指针：一个指向 itab（包含接口类型、变量类型及方法表），一个指向具体数据 data。

2 空接口 interface{} 是 nil 吗？
 -不是，interface{}是 eface 双指针结构体。
 -包含类型指针和数据指针，只有两者为 nil 时 interface{}才是 nil

3 接口的动态派发（Dynamic Dispatch）是什么？
 -当接口调用方法时，程序会在运行时查看接口 Header 里的 itab。itab 中缓存了该类型实现该接口的所有方法地址。程序根据方法名偏移量找到对应的函数指针并跳转执行。

4 什么是隐式实现（Duck Typing）？有什么好处？
 -Go 接口不需要显式声明 implements。只要一个类型实现了接口要求的所有方法，它就自动实现了该接口。
 -好处：解耦。你可以在不修改原始代码的情况下，为它定义接口（例如给第三方库写接口）。

5 如何在编译期检查是否实现了接口？
 -使用空赋值。这在很多开源库（如 uber-go/zap）中很常见。
 -见 TestInterfaceImp

6 接口调用和直接调用有啥区别？
 -接口调用比直接调用慢
 -逃逸分析：接口动态分配往往会导致对象逃逸到堆上，导致 gc 压力。
 -间接寻址：需要通过 itab 找到方法地址，增加了一次指针跳转。
 -无法内联：编译器很难对接口调用进行内联优化，编译器在编译时不知道接口后面具体是谁，没法把函数直接嵌入进去。

*/

// TestInterfaceNil nil 接口
func TestInterfaceNil(t *testing.T) {
	var p *int = nil
	var i interface{} = p

	fmt.Printf("p 的值是 nil 吗: %v\n", p == nil) // true
	fmt.Printf("i 的值是 nil 吗: %v\n", i == nil) // false! 坑就在这

	// 揭秘：看看 i 的内部
	type eface struct {
		_type uintptr
		data  uintptr
	}
	e := (*eface)(unsafe.Pointer(&i))
	fmt.Printf("i 的类型指针: %x, 数据指针: %x\n", e._type, e.data)
	//对于 p：它是一个纯指针，值为 0x0。
	//对于 i：它是一个 eface 结构体。
	// i._type 指向了 int 的类型信息（不为 nil）。
	// i.data 确实是 0x0 (为 nil)。
	//Go 规定，只有当 eface 的 _type 和 data 同时为 nil 时，接口才等于 nil。
}

// TestInterfaceEmpty 空接口作为万能容器
func TestInterfaceEmpty(t *testing.T) {
	m := make(map[string]interface{})
	m["name"] = "Gemini"
	m["age"] = 1

	// 取值记得要“断言”
	age, ok := m["age"].(int)
	if ok {
		fmt.Printf("断言成功: %d\n", age)
	}
}

// 检查是否实现了接口（编译期保障）
type Worker interface {
	DoWork()
}
type Developer struct{}

func (d *Developer) DoWork() {}

// 检查 Developer 是否实现了 Worker
// 建议在项目中写，体现你的严谨
var _ Worker = (*Developer)(nil)

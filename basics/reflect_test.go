package basics

import (
	"fmt"
	"reflect"
	"testing"
)

/*
Q1: 反射是什么？Go 的反射基于什么？
Q2: reflect.Type 和 reflect.Value 分别是什么？
Q3: 反射三大定律是什么？
Q4: 反射能修改值吗？什么条件下能改？
Q5: 反射的性能开销有多大？
Q6: 反射有哪些实际应用场景？
Q7: interface 的底层结构和反射的关系？
Q8: 反射怎么调用方法？怎么创建对象？

---
Q1: 反射是什么？Go 的反射基于什么？
【理解】
反射 = 程序在运行时检查和操作自身类型信息的能力。

Go 的反射基于 interface 的底层结构。任何值赋给 interface{} 时，会存储两个信息：
  1. 类型信息（_type 指针）：指向类型元数据（大小、方法集、kind 等）
  2. 数据指针（data）：指向实际值

reflect 包就是把 interface{} 拆开，让你能访问这两部分：
  reflect.TypeOf(v)  -> 提取类型信息 -> reflect.Type
  reflect.ValueOf(v) -> 提取值信息   -> reflect.Value

为什么需要反射？
  Go 是静态类型语言，编译期确定类型。但有些场景需要运行时处理未知类型：
  - JSON/XML 序列化（不知道结构体有哪些字段）
  - ORM 框架（根据结构体字段生成 SQL）
  - 依赖注入（根据类型自动装配）
  - fmt.Println（打印任意类型）
【回答】
反射是程序在运行时检查和操作自身类型信息的能力。Go 的反射基于 interface 的底层结构——任何值赋给 interface{} 时存储了类型指针和数据指针，reflect 包就是把这两部分拆开暴露给你。
reflect.TypeOf 提取类型信息，reflect.ValueOf 提取值信息。
需要反射的场景：JSON 序列化、ORM、依赖注入、fmt 打印——这些都需要处理编译期未知的类型。

---
Q2: reflect.Type 和 reflect.Value 分别是什么？
【理解】
reflect.Type（类型信息，只读）：
  - Kind()：底层种类（int、struct、ptr、slice、map...）
  - Name()：类型名（"User"、"int"）
  - NumField() / Field(i)：结构体字段信息
  - NumMethod() / Method(i)：方法信息
  - Elem()：指针/slice/map 等的元素类型

reflect.Value（值信息，可读可写）：
  - Interface()：转回 interface{}
  - Int() / String() / Float()：取具体值
  - Set(v)：设置值（需要可寻址）
  - CanSet()：是否可修改
  - Field(i)：结构体第 i 个字段的 Value
  - Call(args)：调用方法

关系：
  Type 描述"是什么类型"，Value 描述"具体是什么值"。
  Value 可以通过 .Type() 获取对应的 Type。
  Type 是纯类型信息，不持有具体值。
【回答】
reflect.Type 是类型信息（只读）：Kind、Name、字段列表、方法列表等。
reflect.Value 是值信息（可读写）：取值、设值、调用方法、访问字段等。
Type 描述"是什么类型"，Value 描述"具体是什么值"。Value 可以通过 .Type() 获取对应的 Type。

---
Q3: 反射三大定律是什么？
【理解】
Go 官方博客定义的反射三大定律：

定律 1：从 interface 到反射对象
  reflect.TypeOf(v) 和 reflect.ValueOf(v) 把 interface 转成反射对象。
  任何值传入反射函数时，先隐式转成 interface{}，再拆解。

定律 2：从反射对象到 interface
  reflect.Value 可以通过 .Interface() 方法转回 interface{}。
  v.Interface().(int) 可以取回原始值。

定律 3：要修改反射对象的值，它必须是可设置的（settable）
  可设置 = 反射对象持有的是原始变量的地址（可寻址）。
  reflect.ValueOf(x) 传的是 x 的拷贝，不可设置。
  reflect.ValueOf(&x).Elem() 传的是 x 的地址再解引用，可设置。

  x := 42
  v := reflect.ValueOf(x)    // v 持有 x 的拷贝，CanSet()=false
  v = reflect.ValueOf(&x).Elem()  // v 持有 x 的地址，CanSet()=true
  v.SetInt(100)              // x 变成 100
【回答】
定律 1：interface 可以转成反射对象（TypeOf/ValueOf）。
定律 2：反射对象可以转回 interface（.Interface()）。
定律 3：要修改值，反射对象必须可设置——必须传指针再 Elem()，直接传值是拷贝不能改。
核心：reflect.ValueOf(&x).Elem() 才能修改 x，reflect.ValueOf(x) 只是拷贝。

---
Q4: 反射能修改值吗？什么条件下能改？
【理解】
能改的条件：CanSet() == true，即反射对象持有原始变量的地址。

规则：
  1. 必须传指针：reflect.ValueOf(&x).Elem()
  2. 结构体字段必须导出（大写开头）：未导出字段 CanSet()=false
  3. map 的 value 不可直接 Set（需要用 MapIndex + SetMapIndex）

修改结构体字段：
  type User struct { Name string; age int }
  u := User{Name: "Go", age: 10}
  v := reflect.ValueOf(&u).Elem()
  v.FieldByName("Name").SetString("Rust")  // OK，导出字段
  v.FieldByName("age").SetInt(20)           // panic！未导出字段不可 Set

修改 slice 元素：
  s := []int{1, 2, 3}
  v := reflect.ValueOf(s)  // slice 本身是引用类型，元素可修改
  v.Index(0).SetInt(100)   // OK
【回答】
能改的条件：传指针再 Elem()，且字段必须是导出的（大写开头）。
reflect.ValueOf(&x).Elem() 可设置；reflect.ValueOf(x) 是拷贝不可设置。
结构体未导出字段即使传了指针也不能 Set（CanSet()=false）。
slice 元素可以直接修改（slice 本身是引用类型）。

---
Q5: 反射的性能开销有多大？
【理解】
反射比直接操作慢 1~2 个数量级：

操作对比（大致数量级）：
  直接字段访问：~1ns
  反射字段访问：~100ns（慢 100 倍）
  直接方法调用：~1ns
  反射方法调用：~1μs（慢 1000 倍）

开销来源：
  1. 类型检查：每次操作都要检查 Kind、CanSet 等
  2. 内存分配：Value 结构体、interface 装箱/拆箱
  3. 间接寻址：通过指针链访问数据
  4. 无法内联：编译器无法优化反射调用

优化建议：
  - 热路径避免反射，用代码生成替代（如 easyjson 替代 encoding/json）
  - 缓存 reflect.Type（TypeOf 结果可复用）
  - 用 unsafe 替代反射（极端性能场景，牺牲安全性）
  - 泛型（Go 1.18+）能替代部分反射场景
【回答】
反射比直接操作慢 1~2 个数量级：字段访问慢约 100 倍，方法调用慢约 1000 倍。
开销来源：运行时类型检查、interface 装箱拆箱、间接寻址、无法内联优化。
优化：热路径避免反射、缓存 Type、用代码生成替代（easyjson）、泛型替代部分场景。
大多数业务代码反射不是瓶颈，只有高频热路径才需要优化。

---
Q6: 反射有哪些实际应用场景？
【理解】
1. 序列化/反序列化：encoding/json、encoding/xml
   通过反射遍历结构体字段，读取 tag，生成/解析 JSON

2. ORM 框架：gorm、xorm
   根据结构体字段名和 tag 生成 SQL，把查询结果映射回结构体

3. 依赖注入：wire、dig
   根据类型信息自动装配依赖

4. RPC 框架：gRPC、net/rpc
   根据方法签名自动生成调用代码

5. 测试框架：testify
   assert.Equal 需要比较任意类型的值

6. fmt 包：
   fmt.Println 需要打印任意类型，内部大量使用反射

7. 配置解析：viper
   把配置文件映射到结构体字段

8. 验证框架：validator
   根据 struct tag 做字段校验
【回答】
主要应用：
JSON/XML 序列化（encoding/json 遍历字段读 tag）。
ORM（gorm 根据字段生成 SQL）。
依赖注入（根据类型自动装配）。
RPC 框架（根据方法签名生成调用）。
fmt.Println（打印任意类型）。
配置解析和验证框架（struct tag 驱动）。
共同点：都需要在运行时处理编译期未知的类型结构。

---
Q7: interface 的底层结构和反射的关系？
【理解】
interface 底层有两种结构：

eface（空接口 interface{}）：
  type eface struct {
      _type *_type    // 类型信息
      data  unsafe.Pointer  // 数据指针
  }

iface（非空接口，有方法）：
  type iface struct {
      tab  *itab           // 类型+方法表
      data unsafe.Pointer  // 数据指针
  }

反射的本质：
  reflect.TypeOf(v)  -> 读取 eface._type 或 iface.tab._type
  reflect.ValueOf(v) -> 读取 eface.data 或 iface.data + 类型信息

所以反射的前提是值被装箱成 interface{}：
  var x int = 42
  reflect.ValueOf(x)  // x 先隐式转成 interface{} -> eface{_type: intType, data: &42}
                      // 然后 reflect 拆开这个 eface

这也解释了为什么反射有性能开销——每次都要经过 interface 装箱。
【回答】
interface{} 底层是 eface 结构体，存了类型指针（_type）和数据指针（data）。
反射的本质就是拆开 interface：TypeOf 读 _type，ValueOf 读 data + 类型信息。
任何值传入反射函数时先隐式装箱成 interface{}，这也是反射有性能开销的原因——每次都要经过装箱。

---
Q8: 反射怎么调用方法？怎么创建对象？
【理解】
调用方法：
  v := reflect.ValueOf(obj)
  method := v.MethodByName("Hello")
  args := []reflect.Value{reflect.ValueOf("world")}
  result := method.Call(args)

创建对象：
  // 创建指针
  t := reflect.TypeOf(User{})
  ptr := reflect.New(t)  // 等价于 new(User)，返回 *User 的 Value
  ptr.Elem().FieldByName("Name").SetString("Go")

  // 创建 slice
  sliceType := reflect.SliceOf(reflect.TypeOf(0))
  s := reflect.MakeSlice(sliceType, 0, 10)

  // 创建 map
  mapType := reflect.MapOf(reflect.TypeOf(""), reflect.TypeOf(0))
  m := reflect.MakeMap(mapType)

注意：
  Call 的参数和返回值都是 []reflect.Value
  方法必须是导出的才能通过 MethodByName 找到
  New 返回的是指针类型的 Value
【回答】
调用方法：v.MethodByName("Name").Call([]reflect.Value{...})，参数和返回值都是 []reflect.Value。
创建对象：reflect.New(type) 创建指针（等价于 new），reflect.MakeSlice/MakeMap 创建集合类型。
注意：方法必须导出才能 MethodByName 找到，New 返回的是指针类型需要 Elem() 才能操作字段。

*/

// TestReflectBasic Type 和 Value 基础
func TestReflectBasic(t *testing.T) {
	x := 42
	fmt.Println("Type:", reflect.TypeOf(x))   // int
	fmt.Println("Value:", reflect.ValueOf(x))  // 42
	fmt.Println("Kind:", reflect.TypeOf(x).Kind()) // int
}

// TestReflectStruct 结构体反射
func TestReflectStruct(t *testing.T) {
	type User struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	u := User{Name: "Go", Age: 15}
	t2 := reflect.TypeOf(u)
	v := reflect.ValueOf(u)

	for i := 0; i < t2.NumField(); i++ {
		field := t2.Field(i)
		value := v.Field(i)
		tag := field.Tag.Get("json")
		fmt.Printf("字段: %s, 值: %v, tag: %s\n", field.Name, value, tag)
	}
}

// TestReflectModify 反射修改值
func TestReflectModify(t *testing.T) {
	x := 42
	v := reflect.ValueOf(&x).Elem() // 必须传指针再 Elem
	fmt.Println("CanSet:", v.CanSet()) // true
	v.SetInt(100)
	fmt.Println("修改后:", x) // 100
}

// TestReflectCall 反射调用方法
func TestReflectCall(t *testing.T) {
	type Greeter struct{}
	// 方法定义在下面

	g := greeter{}
	v := reflect.ValueOf(g)
	method := v.MethodByName("Hello")
	result := method.Call([]reflect.Value{reflect.ValueOf("反射")})
	fmt.Println(result[0]) // "Hello, 反射"
}

type greeter struct{}

func (g greeter) Hello(name string) string {
	return fmt.Sprintf("Hello, %s", name)
}

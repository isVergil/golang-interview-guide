package basics

import (
	"fmt"
	"testing"
)

/*
Q1: Go 有哪些常用的语法糖？
Q2: := 短变量声明有什么坑？
Q3: range 的底层机制？有什么注意事项？
Q4: ... 可变参数和展开语法是什么？
Q5: 方法接收器的自动转换规则？
Q6: 类型推断在哪些场景下生效？
Q7: Go 的多返回值和命名返回值有什么用？
Q8: init 函数的执行顺序和使用场景？

---
Q1: Go 有哪些常用的语法糖？
【理解】
Go 语法糖不多，但每个都很实用：

1. := 短变量声明：自动推导类型，只能在函数内用
2. [...] 数组长度推导：arr := [...]int{1,2,3} 编译器自动算长度
3. range 迭代：统一遍历 slice/map/channel/string
4. ... 可变参数：func f(args ...int)，调用时 f(1,2,3)
5. ... 展开切片：append(s1, s2...)
6. 方法接收器自动转换：值调指针方法自动取地址，指针调值方法自动解引用
7. 结构体部分初始化：未指定字段自动零值
8. 多返回值 + 空白标识符：_, err := f()
9. defer/go 后面直接跟匿名函数

Go 的设计哲学：语法糖要少而精，不搞隐式魔法。
每个语法糖都有明确的语义，不会让人困惑"这背后到底发生了什么"。
【回答】
Go 常用语法糖：:= 短变量声明、[...] 数组长度推导、range 统一迭代、... 可变参数和展开、方法接收器自动转换、结构体部分初始化（未指定字段零值）。
Go 的设计哲学是语法糖少而精，不搞隐式魔法，每个糖都有明确语义。

---
Q2: := 短变量声明有什么坑？
【理解】
:= 的规则：
  1. 只能在函数内使用（包级别必须用 var）
  2. 左边至少有一个新变量时才能用（否则编译错误）
  3. 会创建新变量，可能遮蔽外层变量（shadowing）

经典坑——变量遮蔽：
  var err error
  if true {
      val, err := doSomething()  // 这里的 err 是新变量！不是外层的 err
      _ = val
  }
  // 外层 err 仍然是 nil，bug！

  修复：
  var err error
  var val int
  if true {
      val, err = doSomething()  // 用 = 赋值给外层变量
  }

多变量 := 的规则：
  x := 1
  x, y := 2, 3  // 合法！x 是已有变量，y 是新变量，至少一个新的就行
  // x 被重新赋值为 2，y 是新声明

for 循环中的 :=：
  for i := 0; i < 10; i++ { ... }
  // i 的作用域只在 for 块内，外面访问不到
【回答】
主要坑是变量遮蔽（shadowing）：if/for 块内 := 创建的是新变量，不会修改外层同名变量。
经典 bug：if 块内 err := f() 创建了新 err，外层 err 仍然是 nil。
修复：先声明变量，块内用 = 赋值。
规则：:= 左边至少有一个新变量才能用；只能在函数内用；多变量时已有变量会被重新赋值。

---
Q3: range 的底层机制？有什么注意事项？
【理解】
range 会在迭代前对被遍历对象做一次拷贝（Go 1.22 之前）：
  s := []int{1, 2, 3}
  for i, v := range s {  // s 被拷贝了一份（slice header 拷贝，底层数组共享）
      ...
  }

注意事项：
  1. range slice 时 v 是值拷贝，修改 v 不影响原 slice
     for _, v := range s { v = 100 }  // s 不变！v 是拷贝

  2. range 大结构体 slice 有性能问题（每次迭代拷贝整个结构体）
     修复：用索引 for i := range s { s[i].Field = ... }

  3. range map 的遍历顺序是随机的（Go 故意打乱，防止依赖顺序）

  4. range string 遍历的是 rune（Unicode 码点），不是 byte
     for i, r := range "你好" { ... }  // r 是 rune 类型

  5. Go 1.22+ 循环变量语义变了：每次迭代创建新变量，不再共享地址
     （解决了循环闭包陷阱）

range channel：
  for v := range ch { ... }  // 一直读直到 ch 被 close
【回答】
range 注意事项：
v 是值拷贝，修改 v 不影响原数据——大结构体用索引访问避免拷贝开销。
range map 顺序随机（Go 故意的）。
range string 遍历的是 rune 不是 byte。
Go 1.22+ 循环变量每次迭代是新变量，解决了闭包陷阱。
range channel 会一直读直到 close。

---
Q4: ... 可变参数和展开语法是什么？
【理解】
■ 可变参数（函数定义）：
  func sum(nums ...int) int {  // nums 的类型是 []int
      total := 0
      for _, n := range nums { total += n }
      return total
  }
  sum(1, 2, 3)  // nums = []int{1, 2, 3}

■ 展开切片（函数调用）：
  s := []int{1, 2, 3}
  sum(s...)  // 等价于 sum(1, 2, 3)，把 slice 展开成多个参数

■ append 的展开：
  s1 := []int{1, 2}
  s2 := []int{3, 4}
  s1 = append(s1, s2...)  // 把 s2 展开追加到 s1

底层：可变参数本质就是 slice，编译器帮你把多个参数打包成 slice。
展开（...）是反过来，把 slice 拆开传入。

注意：
  - 可变参数必须是函数最后一个参数
  - 展开只能用于 slice，不能用于数组（需要先 arr[:]）
  - interface{} 可变参数：func Printf(format string, args ...interface{})
【回答】
可变参数：func f(args ...int)，args 本质是 []int，调用时 f(1,2,3) 自动打包。
展开：f(slice...) 把 slice 拆开传入，常见于 append(s1, s2...)。
规则：可变参数必须是最后一个参数；展开只能用于 slice 不能用于数组。
底层就是 slice 的打包和拆包，没有额外开销。

---
Q5: 方法接收器的自动转换规则？
【理解】
Go 会自动在值和指针之间转换来调用方法：

规则：
  值类型变量可以调用指针接收器方法：编译器自动取地址 (&v).Method()
  指针类型变量可以调用值接收器方法：编译器自动解引用 (*p).Method()

  type User struct { Name string }
  func (u *User) SetName(name string) { u.Name = name }  // 指针接收器
  func (u User) GetName() string { return u.Name }        // 值接收器

  u := User{Name: "Go"}
  u.SetName("Rust")   // 编译器转为 (&u).SetName("Rust")
  p := &u
  p.GetName()         // 编译器转为 (*p).GetName()

限制（不能自动转换的场景）：
  接口满足性检查时不会自动转换：
    type Namer interface { SetName(string) }
    var n Namer = u   // 编译错误！User 没有实现 SetName（需要 *User）
    var n Namer = &u  // OK，*User 实现了 SetName

  map 的 value 不可寻址：
    m := map[string]User{"a": {}}
    m["a"].SetName("x")  // 编译错误！map value 不可取地址
【回答】
值可以调指针方法（自动取地址），指针可以调值方法（自动解引用）。
但接口满足性检查不会自动转换：值类型不满足指针接收器的接口，必须用指针。
map 的 value 不可寻址，不能直接调用指针接收器方法。

---
Q6: 类型推断在哪些场景下生效？
【理解】
Go 的类型推断场景：

1. := 短变量声明：
   x := 42          // int
   s := "hello"     // string
   m := map[string]int{}  // map[string]int

2. range 循环变量：
   for i, v := range slice { ... }  // i 是 int，v 是元素类型

3. 函数返回值接收：
   f, err := os.Open("file")  // f 是 *os.File，err 是 error

4. 泛型类型推断（Go 1.18+）：
   Min(1, 2)  // 推断 T=int

5. 常量的默认类型：
   const x = 42    // 无类型常量，使用时根据上下文确定类型
   var y float64 = x  // x 在这里是 float64

不能推断的场景：
  - 包级别变量必须用 var 显式声明
  - 函数参数类型必须显式写
  - 接口类型不能从实现推断
【回答】
Go 类型推断场景：:= 声明、range 循环变量、函数返回值接收、泛型类型参数、常量使用时的类型确定。
不能推断：包级别变量、函数参数类型、接口类型。
Go 的推断是局部的、简单的，不像 Haskell 那样全局推断——保持代码可读性。

---
Q7: Go 的多返回值和命名返回值有什么用？
【理解】
■ 多返回值：
  Go 函数可以返回多个值，最常见的是 (result, error) 模式。
  func Open(name string) (*File, error) { ... }
  f, err := os.Open("file")

  底层：多返回值通过栈传递，编译器在调用者栈帧上预留空间。

■ 命名返回值：
  func divide(a, b int) (result int, err error) {
      if b == 0 { err = errors.New("divide by zero"); return }
      result = a / b
      return  // 裸 return，返回当前 result 和 err 的值
  }

  用途：
    1. 文档作用：函数签名就能看出返回值含义
    2. defer 修改返回值：defer 闭包可以直接操作命名返回值
    3. 裸 return：直接 return 不写值（不推荐，可读性差）

  注意：
    命名返回值会被初始化为零值
    裸 return 在长函数中可读性差，建议显式写返回值
【回答】
多返回值：Go 的核心模式是 (result, error)，强迫调用者处理错误。底层通过栈传递。
命名返回值：函数签名即文档、defer 可以修改返回值、支持裸 return。
建议：命名返回值用于 defer 修改场景和文档目的，裸 return 在长函数中避免使用（可读性差）。

---
Q8: init 函数的执行顺序和使用场景？
【理解】
init 函数的特殊性：
  - 没有参数没有返回值
  - 不能被手动调用
  - 一个文件可以有多个 init（按出现顺序执行）
  - 一个包可以有多个文件各自有 init

执行顺序：
  1. 先递归初始化所有导入的包（依赖包先初始化）
  2. 包级别变量按声明顺序初始化
  3. init 函数按文件名字母序、文件内出现顺序执行
  4. main 包的 init 最后执行
  5. 最后执行 main()

  import 顺序：A 导入 B，B 导入 C
  执行顺序：C.init -> B.init -> A.init -> A.main

使用场景：
  - 注册驱动：import _ "github.com/go-sql-driver/mysql"（副作用导入）
  - 初始化全局配置/连接池
  - 注册编解码器

注意：
  - init 里不要做耗时操作（会拖慢启动）
  - init 里的错误只能 panic（没有返回值）
  - 过度使用 init 会让初始化顺序难以追踪
【回答】
init 函数自动执行，不能手动调用。执行顺序：依赖包先初始化 -> 包级变量 -> init -> main。
一个文件可以有多个 init，按出现顺序执行。
使用场景：注册数据库驱动（副作用导入）、初始化全局配置。
注意：不要做耗时操作、错误只能 panic、过度使用会让初始化顺序难追踪。

*/

// TestShortVarShadow := 变量遮蔽陷阱
func TestShortVarShadow(t *testing.T) {
	x := "外层"
	if true {
		x := "内层" // 新变量，遮蔽了外层 x
		fmt.Println("块内:", x)
	}
	fmt.Println("块外:", x) // 仍然是"外层"
}

// TestRangeValue range 值拷贝陷阱
func TestRangeValue(t *testing.T) {
	type Item struct{ Val int }
	items := []Item{{1}, {2}, {3}}

	// 错误：v 是拷贝，修改无效
	for _, v := range items {
		v.Val = 100
	}
	fmt.Println("range 值拷贝:", items) // [{1} {2} {3}] 没变

	// 正确：用索引直接修改
	for i := range items {
		items[i].Val = 100
	}
	fmt.Println("索引修改:", items) // [{100} {100} {100}]
}

// TestVariadic 可变参数和展开
func TestVariadic(t *testing.T) {
	fmt.Println("sum:", sum(1, 2, 3, 4, 5))

	s := []int{10, 20, 30}
	fmt.Println("sum slice:", sum(s...)) // 展开 slice
}

func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// TestMethodAutoConvert 方法接收器自动转换
func TestMethodAutoConvert(t *testing.T) {
	type User struct{ Name string }

	u := struct{ Name string }{"Go"}
	p := &u

	// 值调值方法、指针调值方法都可以
	fmt.Println(u.Name)
	fmt.Println(p.Name) // 自动解引用 (*p).Name
}

// BenchmarkRangePerformance range 大结构体的性能对比
func BenchmarkRangePerformance(b *testing.B) {
	type BigStruct struct {
		Data [1024]int
	}
	items := make([]BigStruct, 1000)

	b.Run("ValueCopy", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, item := range items {
				_ = item.Data[0]
			}
		}
	})

	b.Run("IndexAccess", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for i := range items {
				_ = items[i].Data[0]
			}
		}
	})
}

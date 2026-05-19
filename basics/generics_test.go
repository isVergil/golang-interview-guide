package basics

import (
	"fmt"
	"testing"
)

/*
Q1: Go 泛型是什么？解决了什么问题？
Q2: 类型参数和类型约束分别是什么？
Q3: 怎么定义类型约束？interface 在泛型中的新用法？
Q4: ~ 符号是什么意思？为什么需要它？
Q5: comparable 约束是什么？哪些类型满足？
Q6: any 和 interface{} 有什么区别？
Q7: 泛型函数和泛型类型分别怎么用？
Q8: 泛型的类型推断是怎么工作的？什么时候必须显式指定？
Q9: 泛型有哪些限制？不能做什么？
Q10: 泛型的性能开销？底层实现原理？

---
Q1: Go 泛型是什么？解决了什么问题？
【理解】
Go 1.18 引入泛型（Generics），核心是类型参数化——让函数和类型可以适用于多种类型，而不需要为每种类型写重复代码。

没有泛型之前的痛点：
  1. 用 interface{} 做通用容器：丢失类型安全，需要类型断言，运行时才能发现错误
  2. 为每种类型写重复函数：MinInt、MinFloat64、MinString... 逻辑完全一样
  3. 用代码生成（go generate）：维护成本高，可读性差

有泛型之后：
  func Min[T constraints.Ordered](a, b T) T {
      if a < b { return a }
      return b
  }
  // Min(1, 2)、Min(1.5, 2.5)、Min("a", "b") 都能用，编译期类型安全

Go 泛型设计哲学是"够用就好"：
  - 没有泛型方法（方法不能有自己的类型参数）
  - 没有特化（specialization）
  - 没有运算符重载
【回答】
Go 1.18 引入泛型，核心是类型参数化，让函数和类型可以适用于多种类型而不丢失类型安全。
解决的问题：之前要么用 interface{} 丢失类型安全，要么为每种类型写重复代码。泛型让你写一次逻辑，编译期保证类型正确。
Go 泛型设计偏保守——没有泛型方法、没有特化，追求简单够用而不是功能完备。

---
Q2: 类型参数和类型约束分别是什么？
【理解】
类型参数（Type Parameter）：函数或类型定义中的占位符类型，用方括号声明。
  func Print[T any](v T) { fmt.Println(v) }
  //         ^ T 就是类型参数

类型约束（Type Constraint）：限制类型参数可以是哪些类型，本质是一个 interface。
  func Min[T constraints.Ordered](a, b T) T { ... }
  //         ^ constraints.Ordered 就是约束，限制 T 必须支持 < 比较

约束决定了你能对类型参数做什么操作：
  any 约束：什么操作都不能做（只能赋值、传参）
  comparable 约束：可以用 == 和 != 比较
  constraints.Ordered 约束：可以用 < > <= >= 比较
  自定义约束：可以调用约束中声明的方法

语法：func FuncName[T Constraint](params) returnType
  多个类型参数：func Map[T any, R any](s []T, f func(T) R) []R
【回答】
类型参数是函数或类型定义中的占位符类型，用方括号 [T ...] 声明。
类型约束是限制类型参数可以是哪些类型的 interface，决定了你能对类型参数做什么操作。
比如 any 约束什么都不能做，comparable 可以 == 比较，constraints.Ordered 可以 < > 比较。
约束越严格，能做的操作越多，但适用的类型越少。

---
Q3: 怎么定义类型约束？interface 在泛型中的新用法？
【理解】
Go 1.18 扩展了 interface 的语义，interface 现在可以包含：
  1. 方法集（传统用法）：type Stringer interface { String() string }
  2. 类型集（新用法）：直接列出允许的类型

类型集语法：
  // 联合类型：T 可以是 int 或 float64 或 string
  type Number interface {
      int | float64 | string
  }

  // 方法 + 类型混合
  type Addable interface {
      ~int | ~float64
      String() string
  }

  // 用 ~ 表示底层类型
  type Signed interface {
      ~int | ~int8 | ~int16 | ~int32 | ~int64
  }

注意：包含类型集的 interface 只能用作约束，不能用作普通变量类型。
  var x Number  // 编译错误！
  func Add[T Number](a, b T) T  // 正确，只能当约束

为什么？类型集 interface 描述的是"类型的集合"而不是"行为的集合"，
作为变量类型没有意义（不知道具体是哪个类型，也不知道有什么方法）。
【回答】
Go 1.18 扩展了 interface，除了传统的方法集，还能直接列出允许的类型（类型集）。
用 | 表示联合：int | float64 | string 表示 T 可以是这三种之一。
方法和类型可以混合：既要满足类型集又要有某些方法。
重要限制：包含类型集的 interface 只能用作泛型约束，不能当普通变量类型用。因为它描述的是"类型的集合"而不是"行为的集合"。

---
Q4: ~ 符号是什么意思？为什么需要它？
【理解】
~ 表示"底层类型是..."（underlying type），包含所有以该类型为底层类型的自定义类型。

没有 ~：
  type Integer interface { int | int64 }
  type MyInt int  // 底层类型是 int
  // MyInt 不满足 Integer 约束！因为 MyInt != int

有 ~：
  type Integer interface { ~int | ~int64 }
  type MyInt int
  // MyInt 满足 Integer 约束！因为 MyInt 的底层类型是 int

为什么需要？
  Go 里自定义类型很常见：type UserID int、type Duration int64
  如果约束只写 int，这些自定义类型全部不能用，泛型就太受限了。
  ~ 让约束能覆盖所有"本质上是某类型"的自定义类型。

限制：~ 后面只能跟底层类型（非接口类型），不能写 ~MyInterface。

底层类型规则：
  type A int       // A 的底层类型是 int
  type B A         // B 的底层类型还是 int（递归到最底层）
  type C []int     // C 的底层类型是 []int
【回答】
~ 表示"底层类型是..."，包含所有以该类型为底层类型的自定义类型。
比如 ~int 不仅匹配 int，还匹配 type MyInt int、type UserID int 等所有底层是 int 的类型。
没有 ~ 的话，自定义类型全部不能用泛型函数，太受限了。Go 里 type XXX int 这种自定义类型很常见，~ 让泛型能覆盖这些场景。
限制：~ 后面只能跟底层类型，不能跟接口。

---
Q5: comparable 约束是什么？哪些类型满足？
【理解】
comparable 是 Go 内置的约束，表示类型支持 == 和 != 操作。

满足 comparable 的类型：
  - 基本类型：int、float64、string、bool 等
  - 指针类型
  - 数组（元素类型是 comparable 的）
  - 结构体（所有字段都是 comparable 的）
  - interface（运行时比较，可能 panic）
  - channel

不满足 comparable 的类型：
  - slice（不能 ==）
  - map（不能 ==）
  - 函数（不能 ==）
  - 包含上述类型的结构体

典型使用场景：泛型 map 的 key、去重、查找。
  func Contains[T comparable](s []T, target T) bool {
      for _, v := range s {
          if v == target { return true }
      }
      return false
  }

注意：Go 1.20 放宽了规则，任何接口类型都满足 comparable，但运行时比较不可比较的值仍然 panic。
【回答】
comparable 是内置约束，表示类型支持 == 和 != 操作。
满足的：基本类型、指针、数组、结构体（字段都可比较）、channel、interface。
不满足的：slice、map、函数——这些不能用 == 比较。
典型场景：泛型 map 的 key 约束、通用的 Contains/IndexOf 函数。
注意：interface 编译期满足 comparable，但运行时底层如果是 slice 等类型，== 会 panic。

---
Q6: any 和 interface{} 有什么区别？
【理解】
Go 1.18 引入 any 作为 interface{} 的类型别名：
  type any = interface{}

两者完全等价，没有任何运行时区别。any 只是语法糖。

但在泛型中 any 约束和 interface{} 参数有本质区别：
  func PrintGeneric[T any](v T)     // 编译期确定 T 的具体类型，类型安全
  func PrintInterface(v interface{}) // 运行时才知道类型，需要类型断言

  PrintGeneric(42)  // T=int，编译期确定，无装箱
  PrintInterface(42) // v 是 interface{}，装箱了
【回答】
any 就是 interface{} 的类型别名，完全等价，没有运行时区别。Go 1.18 引入只是为了代码更简洁。
但在泛型中 [T any] 和参数 interface{} 有本质区别：泛型是编译期确定具体类型，保留类型安全无装箱；interface{} 是运行时装箱，需要类型断言。
建议统一用 any 替代 interface{}，标准库已经这么做了。

---
Q7: 泛型函数和泛型类型分别怎么用？
【理解】
■ 泛型函数：
  func Map[T any, R any](s []T, f func(T) R) []R {
      result := make([]R, len(s))
      for i, v := range s { result[i] = f(v) }
      return result
  }

■ 泛型类型（泛型结构体）：
  type Stack[T any] struct { items []T }
  func (s *Stack[T]) Push(v T) { s.items = append(s.items, v) }
  func (s *Stack[T]) Pop() (T, bool) { ... }

■ 泛型接口：
  type Container[T any] interface { Get() T; Put(T) }

注意：方法不能有自己的类型参数（Go 的限制）。
  func (s *Stack[T]) Map[R any](f func(T) R) *Stack[R]  // 编译错误！
  // 方法的类型参数只能来自接收者
【回答】
泛型函数：func Name[T Constraint](params) returnType，调用时类型通常可以推断。
泛型类型：type Name[T Constraint] struct{...}，方法里用 func (x *Name[T]) Method()。
泛型接口：type Name[T any] interface{...}。
重要限制：方法不能有自己的类型参数，只能用接收者上已声明的类型参数。想要额外的类型参数只能用顶层函数。

---
Q8: 泛型的类型推断是怎么工作的？什么时候必须显式指定？
【理解】
Go 编译器会根据传入的实参推断类型参数，大多数情况不需要显式指定。

能推断的场景：
  Min(1, 2)        // 推断 T=int
  Min("a", "b")   // 推断 T=string

不能推断、必须显式指定的场景：
  1. 类型参数只出现在返回值，不在参数里：
     func New[T any]() *T { return new(T) }
     New[int]()  // 必须显式指定

  2. 泛型类型实例化：
     s := Stack[int]{}  // 必须显式写

  3. 参数是 interface 类型，编译器无法确定具体类型：
     func Do[T any](v interface{}) T { ... }
     Do[string](42)  // 必须显式指定
【回答】
编译器根据实参类型推断类型参数，大多数情况不需要显式写。
必须显式指定的场景：类型参数只出现在返回值不在参数里（编译器没线索）、泛型类型实例化、参数是 interface 无法确定具体类型。
规则：编译器能从实参推出来就不用写，推不出来就必须写。

---
Q9: 泛型有哪些限制？不能做什么？
【理解】
Go 泛型的主要限制：

1. 方法不能有自己的类型参数：
   func (s *Stack[T]) Map[R any](...) // 编译错误
   只能用顶层函数替代

2. 不能对类型参数直接做类型断言/switch：
   func f[T any](v T) { switch v.(type) { } }  // 编译错误
   必须先转 any：switch any(v).(type) { ... }

3. 没有特化（specialization）：
   不能为特定类型提供不同实现

4. 类型集 interface 不能当变量类型：
   type Number interface { ~int | ~float64 }
   var x Number  // 编译错误

5. 不能用 . 访问类型参数的字段：
   func GetName[T any](v T) string { return v.Name }  // 编译错误
   必须通过约束中的方法来访问

6. 不支持泛型类型别名（Go 1.23 之前）
【回答】
主要限制：
方法不能有自己的类型参数——只能用接收者上的，额外参数要用顶层函数。
不能直接对类型参数做 type switch——要先转 any(v) 再 switch。
没有特化——不能为特定类型提供不同实现。
类型集 interface 不能当变量类型——只能做约束。
不能用 . 访问字段——只能通过约束中声明的方法。
总体来说 Go 泛型偏保守，追求简单而非功能完备。

---
Q10: 泛型的性能开销？底层实现原理？
【理解】
Go 泛型采用 GCShape stenciling（模具+字典）混合方案：

■ GCShape = GC 形状，指类型在 GC 视角下的内存布局。
  相同 GCShape 的类型共享同一份生成代码。
  规则：
    - 所有指针类型共享一个 GCShape（都是 8 字节指针）
    - 值类型按大小和是否含指针分组

■ 实现方式：
  编译器不会为每个具体类型都生成一份代码（那样二进制太大）。
  而是按 GCShape 分组，相同 shape 的类型共用一份函数代码。
  通过"字典"（dictionaries）传递具体类型信息（如大小、方法表）。

■ 性能影响：
  - 比 interface{} 好：没有装箱/拆箱，没有运行时类型断言
  - 比手写具体类型稍差：通过字典间接调用有一点开销
  - 指针类型的泛型函数可能无法内联（通过字典调用）
  - 值类型的小函数通常能被优化到接近手写

■ 对比其他语言：
  C++ 模板：完全单态化，性能最好但二进制膨胀
  Java 泛型：类型擦除，运行时全是 Object，有装箱开销
  Go 泛型：折中方案，按 GCShape 分组 + 字典传递类型信息
【回答】
Go 泛型用 GCShape stenciling 方案：按 GC 内存布局分组，相同 shape 的类型共享一份代码，通过字典传递具体类型信息。
性能比 interface{} 好（无装箱），比手写具体类型稍差（字典间接调用有开销，可能阻碍内联）。
所有指针类型共享一份代码，值类型按大小分组。是 C++ 完全单态化和 Java 类型擦除之间的折中。
实际使用中性能差异很小，不需要为此避免泛型。

*/

// TestGenericFunc 泛型函数基础用法
func TestGenericFunc(t *testing.T) {
	fmt.Println("Min int:", Min(3, 7))
	fmt.Println("Min float:", Min(3.14, 2.71))
	fmt.Println("Min string:", Min("apple", "banana"))

	nums := []int{1, 2, 3, 4}
	strs := Map(nums, func(n int) string { return fmt.Sprintf("(%d)", n) })
	fmt.Println("Map result:", strs)
}

// Ordered 约束：支持 < > 比较的类型
type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64 | ~string
}

func Min[T Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func Map[T any, R any](s []T, f func(T) R) []R {
	result := make([]R, len(s))
	for i, v := range s {
		result[i] = f(v)
	}
	return result
}

// TestGenericType 泛型类型（Stack）
func TestGenericType(t *testing.T) {
	s := &Stack[int]{}
	s.Push(1)
	s.Push(2)
	s.Push(3)

	for {
		v, ok := s.Pop()
		if !ok {
			break
		}
		fmt.Println("Pop:", v)
	}
}

type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(v T) {
	s.items = append(s.items, v)
}

func (s *Stack[T]) Pop() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	v := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return v, true
}

// TestComparable comparable 约束演示
func TestComparable(t *testing.T) {
	fmt.Println("Contains 3:", Contains([]int{1, 2, 3, 4}, 3))
	fmt.Println("Contains 5:", Contains([]int{1, 2, 3, 4}, 5))
	fmt.Println("Contains go:", Contains([]string{"go", "rust"}, "go"))
}

func Contains[T comparable](s []T, target T) bool {
	for _, v := range s {
		if v == target {
			return true
		}
	}
	return false
}

// TestTildeConstraint ~ 底层类型约束演示
func TestTildeConstraint(t *testing.T) {
	type MyInt int
	type Score int

	a, b := MyInt(10), MyInt(20)
	fmt.Println("Min MyInt:", Min(a, b))

	x, y := Score(95), Score(88)
	fmt.Println("Min Score:", Min(x, y))
}

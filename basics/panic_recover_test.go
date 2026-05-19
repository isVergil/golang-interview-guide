package basics

import (
	"fmt"
	"testing"
)

/*
Q1: panic 是什么？什么时候会触发？
Q2: recover 是什么？怎么用？
Q3: panic 的传播机制？跨 goroutine 能 recover 吗？
Q4: panic 和 os.Exit 有什么区别？
Q5: panic 的底层实现？_panic 链表是什么？
Q6: 什么时候该用 panic？什么时候该返回 error？
Q7: recover 的返回值是什么？怎么区分不同的 panic？
Q8: panic 后的 defer 执行顺序？

---
Q1: panic 是什么？什么时候会触发？
【理解】
panic 是 Go 的运行时恐慌机制，表示程序遇到了不可恢复的错误。

自动触发（runtime panic）：
  - 数组/slice 越界
  - nil 指针解引用
  - 向已关闭的 channel 发送数据
  - 类型断言失败（不带 ok）
  - 除零（整数除法）
  - map 并发读写（fatal，不可 recover）

手动触发：
  panic("something went wrong")
  panic(fmt.Errorf("error: %w", err))
  panic(任意值)

panic 发生后的流程：
  1. 当前函数立即停止执行
  2. 按 LIFO 执行当前函数的所有 defer
  3. 返回到调用者，重复上述过程（逐层展开调用栈）
  4. 到达 goroutine 顶层，程序崩溃并打印堆栈

注意：map 并发读写触发的是 fatal（throw），不是 panic，不能被 recover。
【回答】
panic 是运行时恐慌，表示不可恢复的错误。
自动触发：越界、nil 解引用、向关闭 channel 发送、类型断言失败、除零等。
触发后逐层展开调用栈，每层执行 defer，到顶层程序崩溃。
注意：map 并发读写是 fatal throw，不是 panic，不能被 recover 捕获。

---
Q2: recover 是什么？怎么用？
【理解】
recover 用于捕获 panic，让程序从恐慌中恢复正常执行。

使用规则：
  1. 只能在 defer 函数内部调用才有效
  2. 直接 defer recover() 无效（返回值没被处理）
  3. 只能捕获当前 goroutine 的 panic

正确写法：
  defer func() {
      if r := recover(); r != nil {
          fmt.Println("recovered:", r)
      }
  }()

错误写法：
  recover()           // 不在 defer 里，无效
  defer recover()     // 虽然在 defer 里，但返回值被丢弃，无法做后续处理

recover 成功后：
  - panic 被捕获，程序不再崩溃
  - 当前函数的剩余 defer 继续执行
  - 函数正常返回（返回零值，除非 defer 修改了命名返回值）
  - 调用者不知道发生过 panic
【回答】
recover 捕获 panic 让程序恢复正常。只能在 defer 函数内部调用才有效。
正确写法：defer func() { if r := recover(); r != nil { 处理 } }()
recover 成功后当前函数正常返回，调用者不知道发生过 panic。
直接 defer recover() 无效，因为返回值被丢弃无法做后续处理。

---
Q3: panic 的传播机制？跨 goroutine 能 recover 吗？
【理解】
panic 只在当前 goroutine 的调用栈内传播，不会跨 goroutine。

传播路径：
  func A() { B() }
  func B() { C() }
  func C() { panic("boom") }

  C panic -> C 的 defer 执行 -> 返回 B -> B 的 defer 执行 -> 返回 A -> A 的 defer 执行
  如果任何一层的 defer 里 recover 了，传播停止。

跨 goroutine：
  go func() { panic("boom") }()
  // 主 goroutine 的 recover 捕获不到子 goroutine 的 panic！
  // 子 goroutine panic 会直接导致整个程序崩溃

所以每个 goroutine 如果可能 panic，必须自己 defer recover：
  go func() {
      defer func() { recover() }()
      // 可能 panic 的代码
  }()

这也是为什么 HTTP server 的每个请求处理都有 recover 中间件——
防止一个请求的 panic 打挂整个进程。
【回答】
panic 只在当前 goroutine 的调用栈内传播，逐层展开直到被 recover 或到达顶层崩溃。
跨 goroutine 不能 recover：子 goroutine 的 panic 主 goroutine 捕获不到，会直接崩溃整个程序。
所以每个可能 panic 的 goroutine 必须自己 defer recover。HTTP server 的 recover 中间件就是这个原理。

---
Q4: panic 和 os.Exit 有什么区别？
【理解】
panic：
  - 触发 defer 执行
  - 可以被 recover 捕获
  - 打印堆栈信息
  - 逐层展开调用栈

os.Exit(code)：
  - 不触发 defer（直接终止进程）
  - 不能被 recover（不是 panic）
  - 不打印堆栈
  - 立即退出，返回状态码

log.Fatal：
  底层就是 log.Print + os.Exit(1)，也不触发 defer。

使用场景：
  panic：程序逻辑错误，需要堆栈信息排查
  os.Exit：正常退出或致命错误需要立即终止（如配置加载失败）
【回答】
panic 会触发 defer、可以被 recover、打印堆栈、逐层展开。
os.Exit 直接终止进程，不触发 defer、不能 recover、不打印堆栈。
log.Fatal 底层就是 os.Exit(1)，也不触发 defer。
选择：需要堆栈排查用 panic；需要立即干净退出用 os.Exit。

---
Q5: panic 的底层实现？_panic 链表是什么？
【理解】
每个 goroutine 的 g 结构体有一个 _panic 字段，指向 panic 链表头。

type _panic struct {
    argp      unsafe.Pointer // defer 的参数指针
    arg       any            // panic 的参数（传给 panic() 的值）
    link      *_panic        // 链表，指向上一个 panic（嵌套 panic）
    recovered bool           // 是否被 recover 了
    aborted   bool           // 是否被中止
}

嵌套 panic 场景：
  defer 里又 panic 了，会创建新的 _panic 节点挂到链表头。
  recover 只恢复链表头的 panic（最近的那个）。

gopanic 函数流程：
  1. 创建 _panic 结构体，挂到 g._panic 链表头
  2. 遍历 g._defer 链表，逐个执行 defer 函数
  3. 如果某个 defer 调用了 recover：标记 recovered=true，跳转到 recovery
  4. 如果所有 defer 执行完还没 recover：fatalpanic -> 打印堆栈 -> exit(2)

gorecover 函数：
  检查 g._panic 链表头是否存在且未 recovered
  如果是，标记 recovered=true，返回 panic 的参数
  否则返回 nil
【回答】
每个 goroutine 有一个 _panic 链表，panic 时创建节点挂到链表头。
gopanic 流程：创建 _panic 节点 -> 遍历执行 defer -> 某个 defer 里 recover 则标记 recovered 恢复 -> 都没 recover 则打印堆栈崩溃。
嵌套 panic（defer 里又 panic）会创建新节点挂到链表头，recover 只恢复最近的那个。

---
Q6: 什么时候该用 panic？什么时候该返回 error？
【理解】
error（99% 的场景）：
  - 所有"预期内可能发生"的错误
  - 文件不存在、网络超时、参数校验失败
  - 库代码必须返回 error，让调用者决定

panic（极少数场景）：
  - 程序初始化失败，无法继续运行（必要配置缺失、必要依赖连不上）
  - 逻辑上不可能到达的代码路径（说明有 bug）
  - 违反了函数的前置条件（如 regexp.MustCompile）

Must 模式：
  标准库的 MustXxx 函数在失败时 panic，用于初始化阶段：
  var re = regexp.MustCompile(`\d+`)  // 编译失败说明正则写错了，应该 panic
  template.Must(template.New("").Parse(tmpl))

库代码原则：
  内部可以 panic，但必须在包边界 recover 转成 error 返回。
  不要让 panic 泄漏到调用者。
【回答】
error 用于所有预期内可能发生的错误（网络超时、文件不存在等），占 99%。
panic 只用于：初始化失败无法继续运行、逻辑上不可能的代码路径（说明有 bug）。
Must 模式（MustCompile 等）用于初始化阶段，失败说明代码写错了。
库代码原则：内部可以 panic，但必须在包边界 recover 转 error，不让 panic 泄漏给调用者。

---
Q7: recover 的返回值是什么？怎么区分不同的 panic？
【理解】
recover() 返回传给 panic() 的参数（any 类型）。

  panic("string error")     -> recover() 返回 "string error"（string）
  panic(42)                 -> recover() 返回 42（int）
  panic(errors.New("err"))  -> recover() 返回 *errors.errorString（error）
  panic(nil)                -> recover() 返回 nil（Go 1.21+ 会包装成 *runtime.PanicNilError）

区分不同 panic：
  defer func() {
      switch r := recover(); r.(type) {
      case nil:
          // 没有 panic
      case string:
          fmt.Println("string panic:", r)
      case error:
          fmt.Println("error panic:", r)
      default:
          fmt.Println("unknown panic:", r)
      }
  }()

Go 1.21 变化：
  panic(nil) 以前 recover 返回 nil，无法区分"没 panic"和"panic(nil)"。
  Go 1.21+ panic(nil) 会被包装成 *runtime.PanicNilError，recover 不再返回 nil。
【回答】
recover() 返回传给 panic() 的参数，类型是 any。可以用 type switch 区分不同类型的 panic。
Go 1.21+ 的变化：panic(nil) 会被包装成 *runtime.PanicNilError，recover 不再返回 nil，解决了无法区分"没 panic"和"panic(nil)"的问题。

---
Q8: panic 后的 defer 执行顺序？
【理解】
panic 后，当前函数的 defer 按 LIFO 顺序全部执行，然后返回调用者继续展开。

  func f() {
      defer fmt.Println("defer 1")  // 最后执行
      defer fmt.Println("defer 2")  // 中间执行
      defer fmt.Println("defer 3")  // 最先执行
      panic("boom")
      fmt.Println("不会执行")
  }
  // 输出：defer 3 -> defer 2 -> defer 1 -> panic 堆栈

如果 defer 里又 panic：
  原来的 panic 被"覆盖"，新 panic 继续传播。
  但原 panic 信息不会丢失（_panic 链表保留了）。

如果 defer 里 recover 后又 panic：
  recover 成功恢复了第一个 panic
  新 panic 开始新的传播流程
【回答】
panic 后当前函数的所有 defer 按 LIFO 顺序全部执行，然后返回调用者继续展开。
panic 之后注册的 defer 不会执行（因为代码已经停了），只有 panic 之前注册的 defer 会执行。
defer 里又 panic 会创建新的 panic 节点继续传播，原 panic 信息保留在链表中。

*/

// TestPanicRecover 基础 panic + recover
func TestPanicRecover(t *testing.T) {
	fmt.Println("开始")
	safeCall()
	fmt.Println("正常继续（panic 已被 recover）")
}

func safeCall() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recovered:", r)
		}
	}()
	panic("something went wrong")
}

// TestPanicDefer panic 后 defer 的执行顺序
func TestPanicDefer(t *testing.T) {
	defer func() { recover() }() // 防止测试崩溃

	defer fmt.Println("defer 1（最后注册，最先执行）")
	defer fmt.Println("defer 2")
	defer fmt.Println("defer 3（最先注册，最后执行）")
	panic("boom")
}

// TestPanicType 区分不同类型的 panic
func TestPanicType(t *testing.T) {
	handlePanic(func() { panic("string panic") })
	handlePanic(func() { panic(42) })
	handlePanic(func() { panic(fmt.Errorf("error panic")) })
}

func handlePanic(f func()) {
	defer func() {
		switch r := recover().(type) {
		case string:
			fmt.Println("string:", r)
		case error:
			fmt.Println("error:", r)
		default:
			fmt.Printf("other(%T): %v\n", r, r)
		}
	}()
	f()
}

// TestPanicCrossGoroutine 跨 goroutine 不能 recover
func TestPanicCrossGoroutine(t *testing.T) {
	done := make(chan struct{})

	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("子 goroutine 自己 recover:", r)
			}
			close(done)
		}()
		panic("子 goroutine panic")
	}()

	<-done
	fmt.Println("主 goroutine 正常继续")
}

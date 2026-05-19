package basics

import (
	"errors"
	"fmt"
	"testing"
)

/*
Q1: Go 的错误处理机制是什么？为什么不用 try-catch？
Q2: error 接口的底层是什么？
Q3: 自定义错误类型怎么实现？
Q4: errors.Is 和 errors.As 有什么区别？分别什么场景用？
Q5: fmt.Errorf 的 %w 和 %v 有什么区别？
Q6: 错误包装（error wrapping）是什么？为什么需要？
Q7: panic 和 error 怎么选？什么时候该 panic？
Q8: Go 错误处理的最佳实践？

---
Q1: Go 的错误处理机制是什么？为什么不用 try-catch？
【理解】
Go 用返回值显式处理错误，不用异常机制（try-catch）。
设计哲学：错误是正常的控制流，不是"异常"。调用者必须显式处理每个可能的错误。

对比：
  Java/Python：错误是异常，可以被跳过（不 catch 就向上抛），容易忽略。
  Go：错误是返回值，不处理就编译警告（unused variable），强迫你面对它。

Go 的选择理由：
  1. 显式 > 隐式：看代码就知道哪里会出错，不会有隐藏的控制流跳转
  2. 简单可预测：没有 try-catch 的作用域和 stack unwinding 复杂性
  3. 性能好：不走异常表，没有 stack trace 开销（除非 panic）

缺点：
  代码里大量 if err != nil { return err }，比较冗长（但清晰）。
【回答】
Go 用返回值显式处理错误，不用 try-catch 异常机制。
设计哲学是：错误是正常的控制流，调用者必须显式处理。好处是代码清晰可预测，看一眼就知道哪里会出错；缺点是 if err != nil 比较冗长。
Go 认为隐式的异常传播容易让开发者忽略错误，显式返回值强迫你面对每一个可能的失败。

---
Q2: error 接口的底层是什么？
【理解】
error 是 Go 内置的接口，只有一个方法：
  type error interface {
      Error() string
  }

任何实现了 Error() string 方法的类型都是 error。

标准库最常用的实现：
  // errors 包
  type errorString struct {
      s string
  }
  func (e *errorString) Error() string { return e.s }

  // errors.New("something went wrong") 返回 *errorString

注意：errors.New 返回的是指针 *errorString，所以两次 errors.New("same") 是不同的 error（地址不同）。
这是故意的设计——相同文本的错误不一定是同一种错误。
【回答】
error 是 Go 内置接口，只有一个方法：Error() string。任何实现了这个方法的类型都满足 error 接口。
标准库 errors.New 底层是一个 errorString 结构体，只存一个字符串。返回的是指针，所以两次 errors.New 相同文本也不相等（地址不同），这是故意的——相同文本不代表同一种错误。

---
Q3: 自定义错误类型怎么实现？
【理解】
当需要携带更多上下文时，自定义错误类型：
  type NotFoundError struct {
      Resource string
      ID       int
  }
  func (e *NotFoundError) Error() string {
      return fmt.Sprintf("%s(id=%d) not found", e.Resource, e.ID)
  }

使用场景：
  - 需要携带额外字段（错误码、资源名、请求 ID 等）
  - 调用者需要根据错误类型做不同处理（用 errors.As 提取）
  - 需要实现 Unwrap() 方法支持错误链

实现 Unwrap 支持错误链：
  type WrapError struct {
      msg string
      err error  // 被包装的原始错误
  }
  func (e *WrapError) Error() string { return e.msg + ": " + e.err.Error() }
  func (e *WrapError) Unwrap() error { return e.err }
【回答】
实现 Error() string 方法就是自定义错误类型。当需要携带额外上下文（错误码、资源名等）时使用。
如果要支持错误链，再实现 Unwrap() error 方法，这样 errors.Is 和 errors.As 就能沿链查找。
典型场景：调用者需要用 errors.As 提取特定类型的错误并读取额外字段做分支处理。

---
Q4: errors.Is 和 errors.As 有什么区别？分别什么场景用？
【理解】
errors.Is(err, target) -> 判断错误链中是否包含某个特定值（值比较）
errors.As(err, &target) -> 从错误链中提取某个特定类型的错误（类型断言）

区别：
  Is：比较的是"值"，适用于 sentinel error（预定义的错误值）
      errors.Is(err, sql.ErrNoRows)  // err 链里有没有 sql.ErrNoRows？
  As：提取的是"类型"，适用于自定义错误类型（需要读取额外字段）
      var notFound *NotFoundError
      errors.As(err, &notFound)  // err 链里有没有 *NotFoundError 类型？

两者都会沿 Unwrap() 链递归查找，不只看最外层。

为什么不能直接用 == 和类型断言？
  err == sql.ErrNoRows  // 如果 err 被 fmt.Errorf("%w", ...) 包装了一层，== 就失败
  err.(*NotFoundError)  // 同理，包装后类型断言失败
  errors.Is / errors.As 会沿错误链逐层剥开查找。
【回答】
errors.Is 判断错误链中是否包含某个特定值——用于 sentinel error，比如 errors.Is(err, sql.ErrNoRows)。
errors.As 从错误链中提取某个特定类型的错误——用于自定义错误类型，比如提取 *NotFoundError 读取额外字段。
两者都会沿 Unwrap 链递归查找。不能直接用 == 或类型断言，因为错误可能被包装了多层。

---
Q5: fmt.Errorf 的 %w 和 %v 有什么区别？
【理解】
%w（wrap）：包装错误，保留错误链，后续可以用 errors.Is/As 查找原始错误。
%v（value）：只拼接文本，丢失错误链，后续 errors.Is/As 找不到原始错误。

  original := errors.New("file not found")

  wrapped := fmt.Errorf("open config: %w", original)
  // errors.Is(wrapped, original) == true ✓

  formatted := fmt.Errorf("open config: %v", original)
  // errors.Is(formatted, original) == false ✗ 错误链断了

底层：%w 会让返回的 error 实现 Unwrap() 方法；%v 只是 string 拼接。

Go 1.20+ 支持多个 %w：
  err := fmt.Errorf("failed: %w and %w", err1, err2)
  // errors.Is(err, err1) == true
  // errors.Is(err, err2) == true
【回答】
%w 包装错误并保留错误链，后续 errors.Is/As 能找到原始错误。
%v 只是文本拼接，错误链断了，后续找不到原始错误。
规则：想让调用者能判断根因用 %w；想隐藏实现细节（不暴露内部错误类型）用 %v。

---
Q6: 错误包装（error wrapping）是什么？为什么需要？
【理解】
错误包装 = 给原始错误加一层上下文描述，同时保留原始错误的可追溯性。

没有包装：
  return err  // 调用者只知道出错了，不知道是在哪一步出的
  // "file not found" -> 哪个文件？哪个环节？

有包装：
  return fmt.Errorf("load user config: %w", err)
  // "load user config: file not found" -> 清楚知道是加载用户配置时找不到文件

好处：
  1. 调用者看到完整的错误路径（像 stack trace 但更轻量）
  2. 同时能用 errors.Is/As 判断根因类型
  3. 不需要打印完整 stack trace，错误信息本身就是路径

原则：
  每一层加自己的上下文，不要重复底层的信息。
  "open db connection: dial tcp: timeout" 而不是 "open db connection: open db connection failed: timeout"
【回答】
错误包装就是用 fmt.Errorf("%w", err) 给原始错误加一层上下文，同时保留错误链。
好处：调用者既能看到完整的错误路径（哪一层出的问题），又能用 errors.Is/As 判断根因做不同处理。
原则：每一层加自己的上下文描述，不要重复底层信息。

---
Q7: panic 和 error 怎么选？什么时候该 panic？
【理解】
核心原则：可恢复的用 error，不可恢复的用 panic。

该用 error（99% 的场景）：
  - 文件不存在、网络超时、参数校验失败、数据库查询失败
  - 任何"预期内可能发生"的错误

该用 panic（极少数场景）：
  - 程序初始化失败（配置加载失败、必要依赖连不上）——没法继续运行
  - 逻辑上不可能到达的代码路径（如 switch default）——说明有 bug
  - 标准库中表示编程错误：slice 越界、nil 指针、类型断言失败

注意：
  - 库代码（给别人用的包）不应该 panic，应该返回 error 让调用者决定
  - 如果必须 panic，在包内部 recover 并转成 error 返回（如 encoding/json）
【回答】
可恢复的用 error，不可恢复的用 panic。
error 用于所有"预期内可能发生"的错误（网络超时、文件不存在等）。
panic 只用于：程序初始化失败无法继续运行、逻辑上不可能的代码路径（说明有 bug）。
库代码不应该 panic，应该返回 error。如果内部必须 panic，要在包内 recover 转成 error。

---
Q8: Go 错误处理的最佳实践？
【理解】
1. 只处理一次：要么处理（记日志/降级/重试），要么包装后向上传递，不要既处理又传递
2. 错误包装加上下文：return fmt.Errorf("doSomething: %w", err)
3. sentinel error 用 var 定义：var ErrNotFound = errors.New("not found")
4. 自定义错误类型用指针 receiver：func (e *MyError) Error() string
5. 判断错误用 errors.Is/As，不要用 == 或字符串匹配
6. 不要忽略错误：_ = doSomething() 是代码坏味道
7. 错误信息小写开头不加标点：遵循 Go 惯例，方便链式拼接
8. 包内部 panic 用 recover 兜底转 error

反模式：
  - if err.Error() == "not found" // 字符串比较，太脆弱
  - log.Error(err); return err    // 既打了日志又向上传，重复处理
  - panic("something wrong")      // 库代码不该 panic
【回答】
核心实践：
只处理一次——要么处理要么传递，不要两者都做。
包装加上下文——fmt.Errorf("xxx: %w", err)，让错误路径清晰。
判断错误用 errors.Is/As——不用 == 或字符串匹配。
sentinel error 用包级 var 定义——var ErrNotFound = errors.New("not found")。
错误信息小写开头不加标点——方便链式拼接。
不要忽略错误——_ = f() 是坏味道。

*/

// TestErrorBasic 基础错误处理
func TestErrorBasic(t *testing.T) {
	err := doSomething(0)
	if err != nil {
		fmt.Println("错误:", err)
	}

	err = doSomething(1)
	if err != nil {
		fmt.Println("错误:", err)
	} else {
		fmt.Println("成功")
	}
}

func doSomething(n int) error {
	if n == 0 {
		return errors.New("n cannot be zero")
	}
	return nil
}

// TestErrorIs errors.Is 沿错误链查找值
func TestErrorIs(t *testing.T) {
	// sentinel error
	var ErrNotFound = errors.New("not found")

	// 包装一层
	wrapped := fmt.Errorf("get user: %w", ErrNotFound)
	// 再包装一层
	doubleWrapped := fmt.Errorf("handle request: %w", wrapped)

	// errors.Is 能沿链找到原始错误
	fmt.Println(errors.Is(doubleWrapped, ErrNotFound)) // true
	fmt.Println(doubleWrapped == ErrNotFound)          // false（直接比较失败）
}

// TestErrorAs errors.As 沿错误链提取类型
func TestErrorAs(t *testing.T) {
	err := &NotFoundError{Resource: "user", ID: 42}
	wrapped := fmt.Errorf("query failed: %w", err)

	var target *NotFoundError
	if errors.As(wrapped, &target) {
		fmt.Printf("资源: %s, ID: %d\n", target.Resource, target.ID)
	}
}

// NotFoundError 自定义错误类型
type NotFoundError struct {
	Resource string
	ID       int
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s(id=%d) not found", e.Resource, e.ID)
}

// TestErrorWrap %w vs %v 对比
func TestErrorWrap(t *testing.T) {
	original := errors.New("connection refused")

	// %w 保留错误链
	withW := fmt.Errorf("connect db: %w", original)
	fmt.Printf("%%w - errors.Is: %v\n", errors.Is(withW, original)) // true

	// %v 断开错误链
	withV := fmt.Errorf("connect db: %v", original)
	fmt.Printf("%%v - errors.Is: %v\n", errors.Is(withV, original)) // false
}

// TestPanicVsError panic 和 error 的使用边界
func TestPanicVsError(t *testing.T) {
	// 库代码内部 panic + recover 转 error 的模式
	result, err := safeDiv(10, 0)
	if err != nil {
		fmt.Println("安全除法错误:", err)
	} else {
		fmt.Println("结果:", result)
	}
}

func safeDiv(a, b int) (result int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic recovered: %v", r)
		}
	}()
	return a / b, nil
}

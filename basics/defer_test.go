package basics

import (
	"fmt"
	"testing"
)

/*
Q1: defer 的执行顺序是什么？
Q2: defer、return、返回值到底谁先谁后？
Q3: defer 的参数预计算是什么意思？
Q4: defer 闭包和 defer 传参有什么区别？
Q5: defer 能修改函数返回值吗？什么情况能，什么情况不能？
Q6: defer 遇到 panic 会怎样？recover 怎么配合？
Q7: defer 的底层实现？性能开销大吗？
Q8: defer 有哪些常见的使用场景和坑？

---
Q1: defer 的执行顺序是什么？
【理解】
defer 是栈结构（LIFO），后注册的先执行。
多个 defer 按注册的逆序执行，像"叠盘子"：最后放上去的最先拿走。
defer 只在函数返回时（return/panic/函数体结束）触发，不是在代码块结束时。
【回答】
defer 是栈结构，后进先出（LIFO）。多个 defer 按注册的逆序执行。
注意 defer 是函数级别的，不是代码块级别的——for 循环里写 defer 不会在每次迭代结束时执行，而是在整个函数返回时才一起执行（容易导致资源泄漏）。

---
Q2: defer、return、返回值到底谁先谁后？
【理解】
return 不是原子操作，编译器拆成三步：
  Step1: 设置返回值（把表达式结果赋给返回值变量）
  Step2: 执行 defer 函数
  Step3: RET 指令，函数彻底退出

伪代码展示：
  func f() int {
      x := 5
      return x
  }
  等价于：
  func f() int {
      x := 5
      returnValue = x   // Step1: 设置返回值
      执行所有 defer     // Step2: 执行 defer
      return            // Step3: RET 退出
  }

关键推论：defer 执行时，返回值已经被设置了，但还没有真正返回给调用者。
所以 defer 有机会修改返回值（前提是命名返回值）。
【回答】
return 不是原子操作，编译器拆成三步：先设置返回值，再执行 defer，最后 RET 退出。
所以 defer 是在返回值被赋值之后、函数真正退出之前执行的。这意味着 defer 有机会读取甚至修改返回值。

---
Q3: defer 的参数预计算是什么意思？
【理解】
defer 语句注册时，参数的值就已经确定了（被"快照"）：
  i := 1
  defer fmt.Println(i)  // 此刻 i=1 就被拷贝进去了
  i++
  // 输出 1，不是 2

原因：defer 注册时，会对参数进行求值并拷贝。后续对变量的修改不会影响已经拷贝的值。
这和函数调用传参是同一个道理——Go 是值传递，defer 注册就相当于"预约了一次函数调用"，参数在预约时就定了。
【回答】
defer 注册时，参数的值就已经被求值并拷贝了，后续修改不影响。
这就是"参数预计算"——defer fmt.Println(i) 在 defer 那一行就把 i 的当前值快照下来了。
本质和函数值传递一样，defer 注册相当于预约了一次函数调用，参数在预约时就确定了。

---
Q4: defer 闭包和 defer 传参有什么区别？
【理解】
两种模式对比：
  x := "初始值"

  // 模式 A：传参（快照）
  defer fmt.Println(x)              // 注册时拷贝 x 的值，输出"初始值"

  // 模式 B：闭包（引用）
  defer func() { fmt.Println(x) }() // 闭包捕获 x 的地址，执行时才去读，输出"修改值"

  // 模式 C：闭包+传参（伪装的快照）
  defer func(val string) { fmt.Println(val) }(x) // x 作为参数被拷贝，输出"初始值"

  x = "修改值"

核心区别：
  传参 -> 值在 defer 注册时确定（快照）
  闭包 -> 值在 defer 执行时才去读（引用，用的是最新值）
【回答】
defer 传参是快照：注册时就把值拷贝进去了，后续修改不影响。
defer 闭包是引用：捕获的是变量地址，执行时才去读，看到的是最新值。
如果在闭包参数里传入变量，那又变成了快照（本质还是值拷贝）。

---
Q5: defer 能修改函数返回值吗？什么情况能，什么情况不能？
【理解】
命名返回值 -> 能改。defer 闭包直接操作的就是返回值变量本身。
匿名返回值 -> 不能改。返回值已经被拷贝到一个调用者看不到的临时位置。

  // 匿名返回值：defer 改不了
  func f1() int {
      i := 5
      defer func() { i++ }()  // 改的是局部变量 i，不是返回值
      return i                 // 返回值=5，不变
  }

  // 命名返回值：defer 能改
  func f2() (result int) {
      result = 5
      defer func() { result++ }()  // 直接改 result 变量
      return result                  // 最终返回 6
  }

原因回到 Q2 的三步分解：
  f1: returnValue = i(5) -> defer改i(变成6) -> return(返回的是 returnValue=5)
  f2: result = 5 -> defer改result(变成6) -> return(返回的就是 result=6)
【回答】
命名返回值能改，匿名返回值不能改。
因为命名返回值时，defer 闭包操作的就是返回值变量本身；匿名返回值时，返回值已经被拷贝到调用者看不到的临时位置，defer 改的只是局部变量。
经典面试题：命名返回值 result=5，defer 里 result++，最终返回 6。

---
Q6: defer 遇到 panic 会怎样？recover 怎么配合？
【理解】
panic 发生时的执行流程：
  Step1: 当前函数停止执行
  Step2: 按 LIFO 顺序执行当前函数的所有 defer
  Step3: 如果某个 defer 里调用了 recover()，则捕获 panic，程序恢复正常
  Step4: 如果没有 recover，panic 沿调用栈向上传播，重复上述过程
  Step5: 到达 goroutine 顶层还没被 recover，程序崩溃

recover 规则：
  - 只有在 defer 函数内部调用 recover 才有效
  - 直接在函数体里写 recover() 无效（因为 panic 时函数体已经停了）
  - recover 只能捕获当前 goroutine 的 panic，跨协程无效

  // 正确写法
  defer func() {
      if r := recover(); r != nil {
          fmt.Println("捕获 panic:", r)
      }
  }()

  // 错误写法（无效）
  recover()  // 不在 defer 里，无效
  defer recover()  // 直接 defer recover() 也无效，因为 recover 的返回值没被处理
【回答】
panic 时会按 LIFO 顺序执行当前函数的所有 defer。如果某个 defer 里调用了 recover()，就能捕获 panic，程序恢复正常；否则 panic 向上传播直到崩溃。
recover 只在 defer 函数内部有效，直接写在函数体里无效。而且 recover 只能捕获当前 goroutine 的 panic，跨协程不行。

---
Q7: defer 的底层实现？性能开销大吗？
【理解】
Go 1.12 及之前：堆分配
  每次 defer 都在堆上分配一个 _defer 结构体，挂到 goroutine 的 defer 链表上。
  开销：约 50ns/次。

Go 1.13：栈分配优化
  大部分 defer 的 _defer 结构体直接在栈上分配，避免堆分配。
  开销：约 35ns/次。

Go 1.14+：开放编码（open-coded defer）
  编译器直接把 defer 内联展开到函数末尾，不再走链表。
  条件：defer 数量 <= 8 且没有在循环里。
  开销：接近直接调用函数，约 6ns/次。

_defer 结构体（简化）：
  type _defer struct {
      siz     int32    // 参数区域大小
      started bool     // 是否已开始执行
      sp      uintptr  // 栈指针
      pc      uintptr  // 程序计数器
      fn      *funcval // defer 的函数
      link    *_defer  // 链表指向下一个 defer
  }

【回答】
Go 1.14 之后引入了开放编码优化，编译器直接把 defer 内联到函数末尾，不走链表，开销接近普通函数调用（约 6ns）。
条件是 defer 数量不超过 8 个且不在循环里。不满足条件时退化为栈分配（~35ns）或堆分配（~50ns）。
所以现在 defer 的性能开销已经很小了，不需要为了性能避免使用 defer。

---
Q8: defer 有哪些常见的使用场景和坑？
【理解】
常见场景：
  1. 资源释放：mu.Unlock()、file.Close()、rows.Close()、conn.Close()
  2. recover panic：服务端兜底，防止一个请求的 panic 打挂整个进程
  3. 修改命名返回值：比如统一给 err 加上下文信息
  4. 计时/打点：记录函数耗时

常见坑：
  1. for 循环里 defer：defer 不会在每次迭代结束时执行，会堆积到函数返回
     修复：把循环体提取成独立函数，或手动在每次迭代结束时释放
  2. defer 后面跟方法调用，receiver 是值还是指针要注意
  3. defer 配合 os.Exit 不执行：os.Exit 直接退出进程，不走 defer
  4. defer nil 函数会 panic：var f func(); defer f() -> panic

【回答】
常见场景：资源释放（锁、文件、连接）、recover panic 兜底、修改命名返回值加上下文、计时打点。
常见坑：
for 循环里 defer 会堆积到函数返回才执行，可能导致文件句柄耗尽——要提取成独立函数。
os.Exit 不会触发 defer。
defer nil 函数会 panic。
参数预计算陷阱——以为用了最新值其实用的是注册时的快照。

*/

// TestDeferOrder LIFO 顺序展示
func TestDeferOrder(t *testing.T) {
	defer fmt.Println("第一个注册（最后执行）")
	defer fmt.Println("第二个注册")
	defer fmt.Println("第三个注册（最先执行）")
	// 输出：三 -> 二 -> 一
}

// TestDeferPrecompute 参数预计算 vs 闭包引用
func TestDeferPrecompute(t *testing.T) {
	x := "初始值"

	// 传参：快照，输出"初始值"
	defer fmt.Println("1.快照:", x)

	// 闭包：引用，输出"修改值"
	defer func() {
		fmt.Println("2.引用:", x)
	}()

	// 闭包+传参：还是快照，输出"初始值"
	defer func(val string) {
		fmt.Println("3.伪装快照:", val)
	}(x)

	x = "修改值"
}

// TestDeferReturn defer 修改返回值
func TestDeferReturn(t *testing.T) {
	fmt.Println("匿名返回值:", returnAnonymous()) // 5
	fmt.Println("命名返回值:", returnNamed())     // 6
}

// 匿名返回值：defer 改不了
func returnAnonymous() int {
	i := 5
	defer func() { i++ }() // 改的是局部变量，不是返回值
	return i               // 返回值=5，定死了
}

// 命名返回值：defer 能改
func returnNamed() (result int) {
	result = 5
	defer func() { result++ }() // 直接改 result
	return result               // 最终返回 6
}

// TestDeferRecover 配合 recover 捕获 panic
func TestDeferRecover(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("捕获 panic:", r)
		}
	}()

	a, b := 1, 0
	fmt.Println("result:", a/b) // 触发 panic: integer divide by zero
}

// TestDeferLoop 循环里 defer 的坑
func TestDeferLoop(t *testing.T) {
	// 错误写法：10 个 defer 堆积到函数返回才执行
	// for i := 0; i < 10; i++ {
	//     f, _ := os.Open(fmt.Sprintf("file_%d.txt", i))
	//     defer f.Close()  // 不会在每次迭代关闭！
	// }

	// 正确写法：提取成独立函数
	// for i := 0; i < 10; i++ {
	//     processFile(fmt.Sprintf("file_%d.txt", i))
	// }

	fmt.Println("(循环 defer 示例已注释，避免实际打开文件)")
}

package basics

import (
	"fmt"
	"testing"
)

/*
Q1: slice 底层结构是怎样的？
Q2: slice 作为参数是值传递还是引用传递？
Q3: append 未触发扩容时，外部能看到新元素吗？
Q4: slice 什么时候触发扩容？怎么扩容的？
Q5: 扩容后容量不是严格 2 倍是为什么？
Q6: 频繁扩容有什么影响？如何避免？
Q7: 切片的内存泄漏？怎么解决？
Q8: slice 地址会变化吗？

---
Q1: slice 底层结构是怎样的？
【理解】
切片本身不是数组，是一个 SliceHeader 结构体：
  type SliceHeader struct {
      Data uintptr  // 指向底层数组第一个元素的指针
      Len  int      // 当前切片的元素个数
      Cap  int      // 底层数组从 Data 开始最多能装多少元素
  }

关键特性：
  1. 共享数组：多个切片可以指向同一个底层数组，修改一个会影响另一个
  2. append 扩容：容量不够时分配新数组，旧数组被 GC（如果无其他引用）
  3. 值语义：切片变量本身是 24 字节的结构体（指针8+长度8+容量8）
【回答】
slice 底层是 SliceHeader 结构体，包含三个字段：Data（指向底层数组的指针）、Len（长度）、Cap（容量）。
多个切片可以共享同一个底层数组，修改会互相影响。append 触发扩容后会分配新数组，旧切片不受影响。

---
Q2: slice 作为参数是值传递还是引用传递？
【理解】
Go 只有值传递，slice 传递的是 SliceHeader 的副本（24 字节）。

两种情况：
  1. 函数内修改已有元素：内外共享底层数组，修改对外部可见
  2. 函数内 append 触发扩容：内部指向新数组，外部不受影响

  func modify(s []int) {
      s[0] = 99    // 外部能看到（共享底层数组）
      s = append(s, 100)  // 如果扩容了，外部看不到
  }
【回答】
Go 只有值传递，slice 传的是 header 副本（指针+长度+容量）。
函数内修改已有元素外部能看到（共享底层数组）；但 append 触发扩容后内部指向新数组，外部完全不受影响。

---
Q3: append 未触发扩容时，外部能看到新元素吗？
【理解】
看不到。即使底层数组有空位，append 确实把数据写进去了（Data 指针没变），
但外部切片的 Len 字段还是旧值，访问不到新写入的位置。

  s := make([]int, 2, 5)  // len=2, cap=5
  func add(s []int) {
      s = append(s, 100)  // 没扩容，底层数组[2]确实写了100
      // 但外部 s 的 Len 还是 2，s[2] 会 panic
  }

本质：外部拿到的是 header 的旧副本，Len 没被更新。
【回答】
看不到。虽然底层数组确实写入了数据（没扩容时 Data 指针不变），但外部切片的 Len 字段是旧值，访问不到新位置。
本质是值传递——外部的 header 副本没有被更新。

---
Q4: slice 什么时候触发扩容？怎么扩容的？
【理解】
当 append 时 len 超过 cap 就触发扩容，整体搬迁（不同于 map 的渐进式）。

Go 1.18 之前（1024 分界线）：
  旧容量 < 1024：直接翻倍（2倍）
  旧容量 ≥ 1024：每次增加 25%（1.25倍），直到满足需求

Go 1.18+（256 分界线，更平滑）：
  期望容量 > 旧容量×2：直接使用期望容量
  旧容量 < 256：直接翻倍
  旧容量 ≥ 256：newcap = oldcap + (oldcap + 3×256) / 4
  随着容量增大，增长比例从 2.0 缓慢下降到 1.25
【回答】
append 时 len 超过 cap 触发扩容。Go 1.18+ 以 256 为分界：小于 256 翻倍，大于等于 256 用公式 newcap = oldcap + (oldcap + 768)/4，增长比例从 2.0 平滑下降到 1.25。
扩容是整体搬迁，分配新数组并拷贝全部数据。

---
Q5: 扩容后容量不是严格 2 倍是为什么？
【理解】
内存对齐机制：CPU 按 4/8 字节对齐读取数据，Go 的内存分配器（基于 TCMalloc）
有固定的 size class（8, 16, 32, 48, 64, 80, 96...字节）。

扩容计算出理论容量后，还要向上取整到最近的 size class。
例如理论算出需要 24 字节，实际分配 32 字节，所以容量比预期大。
【回答】
内存对齐。Go 内存分配器有固定的 size class，扩容后的理论容量会向上取整到最近的规格，所以实际容量比公式算出来的大。

---
Q6: 频繁扩容有什么影响？如何避免？
【理解】
每次扩容 = 一次内存分配 + 一次数据拷贝，频繁扩容极大消耗 CPU 和内存。

避免方法：已知数据量时用 make 预分配容量。
  s := make([]int, 0, 1000)  // 预分配 cap=1000，append 1000 次不扩容
【回答】
频繁扩容导致多次内存分配和数据拷贝，消耗 CPU。
避免方法：已知数据量时 make([]T, 0, cap) 预分配容量，一次分配到位。

---
Q7: 切片的内存泄漏？怎么解决？
【理解】
大切片截取小段后，小段依然引用大切片的底层数组，导致整个大数组无法被 GC。

  huge := make([]int, 1000000)
  small := huge[0:2]  // small 持有百万级数组的引用，GC 不掉

解决：用 copy 断开引用关系。
  small := make([]int, 2)
  copy(small, huge[0:2])  // 独立的底层数组，huge 可以被 GC
【回答】
大切片截取小段后，小段仍引用大切片的底层数组，导致大内存无法释放。
解决方法：用 copy 拷贝到新切片，断开与原数组的引用关系。

---
Q8: slice 地址会变化吗？
【理解】
底层数组地址在两种情况下变化：
  1. 扩容：容量不够触发扩容，分配新数组，地址变了
  2. 不扩容：容量够，地址不变

切片变量本身的地址（&s）不会变，变的是 Data 指针指向的底层数组地址。
【回答】
扩容后底层数组地址会变（分配了新数组），不扩容则不变。
注意区分：切片变量的地址（&s）不变，变的是内部 Data 指针指向的底层数组。

*/

// TestSliceExpandAlgorithm 扩容
func TestSliceExpandAlgorithm(t *testing.T) {
	// 场景：旧容量为 256，再 append 一个元素
	s := make([]int, 256)
	fmt.Printf("扩容前: len=%d, cap=%d\n", len(s), cap(s))

	s = append(s, 1)

	// 根据 1.18+ 公式：newcap = 256 + (256 + 768)/4 = 256 + 256 = 512
	// 但要注意内存对齐，实际结果可能会略有不同
	fmt.Printf("扩容后: len=%d, cap=%d\n", len(s), cap(s))
}

// TestSliceLogic 作为参数传递后扩容
func TestSliceLogic(t *testing.T) {
	outerS := []int{1, 2} // len=2, cap=2
	fmt.Printf("函数前 - 地址: %p, 值: %v\n", outerS, outerS)

	modifySlice(outerS)

	fmt.Printf("函数后 - 地址: %p, 值: %v\n", outerS, outerS)
	// 结果：outerS 依然是 [1, 2]，地址也没变
}

func modifySlice(innerS []int) {
	fmt.Printf("进入函数 - 地址: %p, 长度: %d, 容量: %d\n", innerS, len(innerS), cap(innerS))

	innerS = append(innerS, 100) // 这里触发扩容

	fmt.Printf("append后 - 地址: %p, 长度: %d, 容量: %d\n", innerS, len(innerS), cap(innerS))
	innerS[0] = 999 // 改的是新数组的值
}

// TestSliceCopy 陷阱题：copy 与赋值的区别
func TestSliceCopy(t *testing.T) {
	src := []int{1, 2, 3}
	dst := make([]int, len(src))

	copy(dst, src) // 真正的拷贝，数据独立
	dst[0] = 99
	fmt.Println("src:", src) // src 保持 [1, 2, 3]
}

// TestSliceLeak 性能题：切片导致的内存泄漏
func TestSliceLeak(t *testing.T) {
	// 模拟一个巨大的底层数组
	huge := make([]int, 1000000)

	// 只取其中一小段
	small := huge[0:2]

	// 问：此时 huge 占用的百万级内存能被 GC 吗？
	// 答：不能。因为 small 变量依然持有底层数组的指针。
	// 解决方法：使用 copy(newSmall, huge[0:2])，断开与原大数组的联系。
	_ = small
}

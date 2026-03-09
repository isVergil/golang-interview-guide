package base

import (
	"fmt"
	"testing"
)

/*
切片 Slice
type SliceHeader struct {
    Data uintptr // 指向底层数组第一个元素的指针
    Len  int     // 切片当前的长度
    Cap  int     // 切片底层的容量
}
header 可以理解成 元数据，类似的有字符串、接口的 header，它们在 runtime 包中定义，编译器和运行时（Runtime）全靠它们来管理内存。
共享数组：对切片的修改会影响到共享同一个底层数组的其他切片，除非触发了扩容。
append 扩容：扩容会产生新数组，旧数组会被 GC（如果没有其他引用），地址发生变化。
内存泄漏：大切片截取小段后，由于小段依然引用大切片的底层数组，导致大内存无法释放。

面试题：
1 slice 作为参数是值传递还是引用传递？
 -golang 只有值传递，slice 作为参数传递的是 slice header 的副本
 -扩容：在函数内部扩容不影响外部的 slice
 -改值：共用底层数组，修改会成功

2 slice 作为参数传递时，如果 append 了但是没有触发扩容，外面能看到新元素吗？
 -看不到，虽然底层数组后面可能有空位，append 确实把数据写进去了（此时内外 Data 指针确实还一样），但外部切片的 Len 字段还是旧的值。

3 slice 底层结构是怎样的？
 -切片本身并不是数组，它是一个 SliceHeader（结构体）。
 -切片底层其实是一个结构体，包含三个字段：
  -Data 指针：指向底层数组的起始地址。
  -Len（长度）：当前切片里的元素个数。
  -Cap（容量）：底层数组从指针开始算，最多能装多少元素。

4 讲讲切片的内存泄露？怎么解决？
 -大切片截取小段后，由于小段依然引用大切片的底层数组，导致大内存无法释放，导致内存泄露。
 -使用 copy 函数拷贝切片，断开与原切片的联系

5 slice 地址会变化吗？
 -扩容后地址可能会发生变化，即容量不够触发扩容，会分配新数组并拷贝数据，地址就变了。
 -扩容后如果容量够，地址不变。

6 slice 什么时候触发扩容？怎么扩容的？
 -不同于 map 的渐进式，slice 是整体搬迁的
 -当执行 append 时切片长度 len 超过容量 cap 时，就会触发扩容，申请一块更大的内存，并将旧数据拷贝过去。
 -Go 1.18 之前：1024 分界线，小于直接翻倍，大于库容 1.25 倍
  -如果旧容量 < 1024，则直接翻倍（2倍）。
  -如果旧容量 ≥ 1024，则每次增加 25%（1.25倍），直到大于等于期望容量。
 -Go 1.18 之后：255 分界线，小于直接翻倍，大于则按公式扩容，随着容量变大，增长比例会从 2.0 缓慢下降到 1.25。
  -如果期望容量大于旧容量的两倍，则直接使用期望容量。
  -如果旧容量 < 256，则直接翻倍（2倍）。
  -如果旧容量 ≥ 256，则使用公式：newcap = oldcap + (oldcap + 3*256) / 4。

7 频繁扩容有什么影响？如何避免频繁扩容？
 -频繁扩容会导致多次内存分配和数据拷贝，极大地消耗 CPU 性能。
 -在已知数据量的情况下，使用 make([]T, len, cap) 预分配容量。

8 有的时候扩容后并不是严格的按照 2 倍或者公式来的，是为什么？
 -有内存对齐的机制：cpu 一般是 4 或 8 的倍数读取数据，为了让数据对齐，不至于需要多个分隔的块读取数据，扩容重新分配内存时也会考虑内存规格对齐，分配的比实际理论算出来的多。

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

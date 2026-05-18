package top100liked

import (
	"container/heap"
	"testing"
)

// 295. 数据流的中位数 (Find Median from Data Stream)
//
// 题目描述:
// 设计一个支持以下两种操作的数据结构：
// - addNum(num) 从数据流中添加一个整数到数据结构中
// - findMedian() 返回目前所有元素的中位数
//
// 示例：
// addNum(1), addNum(2), findMedian() -> 1.5
// addNum(3), findMedian() -> 2.0
//
// 提示：用一个大顶堆存较小的一半，一个小顶堆存较大的一半，保持两堆平衡

type MaxHeap []int

func (h MaxHeap) Len() int            { return len(h) }
func (h MaxHeap) Less(i, j int) bool  { return h[i] > h[j] }
func (h MaxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *MaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

type MinHeap []int

func (h MinHeap) Len() int            { return len(h) }
func (h MinHeap) Less(i, j int) bool  { return h[i] < h[j] }
func (h MinHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MinHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

type MedianFinder struct {
	lo *MaxHeap // 左半（较小），堆顶 = 左半最大值
	hi *MinHeap // 右半（较大），堆顶 = 右半最小值
}

func NewMedianFinder() MedianFinder {
	return MedianFinder{lo: &MaxHeap{}, hi: &MinHeap{}}
}

func (mf *MedianFinder) AddNum(num int) {
	heap.Push(mf.lo, num)             // Step 1: 先放 lo
	heap.Push(mf.hi, heap.Pop(mf.lo)) // Step 2: lo顶 → hi（保证顺序）
	if mf.lo.Len() < mf.hi.Len() {    // Step 3: 平衡大小
		heap.Push(mf.lo, heap.Pop(mf.hi))
	}
}

func (mf *MedianFinder) FindMedian() float64 {
	if mf.lo.Len() > mf.hi.Len() {
		return float64((*mf.lo)[0]) // 奇数个：lo 多一个，lo顶就是中位数
	}
	return (float64((*mf.lo)[0]) + float64((*mf.hi)[0])) / 2.0 // 偶数个：两顶平均
}

func TestMedianFinder(t *testing.T) {
	mf := NewMedianFinder()
	mf.AddNum(1)
	mf.AddNum(2)

	if got := mf.FindMedian(); got != 1.5 {
		t.Errorf("FindMedian() = %v, want 1.5", got)
	}

	mf.AddNum(3)
	if got := mf.FindMedian(); got != 2.0 {
		t.Errorf("FindMedian() = %v, want 2.0", got)
	}

	mf.AddNum(4)
	if got := mf.FindMedian(); got != 2.5 {
		t.Errorf("FindMedian() = %v, want 2.5", got)
	}
}

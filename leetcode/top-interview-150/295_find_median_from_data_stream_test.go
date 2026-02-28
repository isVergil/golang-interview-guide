package topinterview150

import (
	"testing"
)

// 295. 数据流的中位数 (Find Median from Data Stream)
//
// 题目描述:
// 中位数是有序整数列表中的中间值。如果列表的大小是偶数，则没有中间值，中位数是两个中间值的平均值。
// 例如 arr = [2,3,4] 的中位数是 3 。
// 例如 arr = [2,3] 的中位数是 (2 + 3) / 2 = 2.5 。
// 实现 MedianFinder 类:
// MedianFinder() 初始化 MedianFinder 对象。
// void addNum(int num) 将数据流中的整数 num 添加到数据结构中。
// double findMedian() 返回到目前为止所有元素的中位数。与实际答案相差 10-5 以内的答案将被接受。
//
// 示例 1：
// 输入：
// ["MedianFinder", "addNum", "addNum", "findMedian", "addNum", "findMedian"]
// [[], [1], [2], [], [3], []]
// 输出：
// [null, null, null, 1.5, null, 2.0]
type MedianFinder struct {
	queMin *hp // 大顶堆，存储较小的一半
	queMax *hp // 小顶堆，存储较大的一半
}

func ConstructorMedianFinder() MedianFinder {
	return MedianFinder{
		queMin: &hp{compare: func(a, b int) bool { return a > b }},
		queMax: &hp{compare: func(a, b int) bool { return a < b }},
	}
}

func (this *MedianFinder) AddNum(num int) {
	minQ, maxQ := this.queMin, this.queMax
	// 1. 插入逻辑：先入大顶堆，再将大顶堆顶弹出送入小顶堆，保证顺序
	if minQ.Len() == 0 || num <= minQ.Ints[0] {
		minQ.push(num)
		// 保持左右平衡：左边最多比右边多一个
		if minQ.Len() > maxQ.Len()+1 {
			maxQ.push(minQ.pop())
		}
	} else {
		maxQ.push(num)
		// 保持左右平衡：右边不能比左边多
		if maxQ.Len() > minQ.Len() {
			minQ.push(maxQ.pop())
		}
	}
}

func (this *MedianFinder) FindMedian() float64 {
	if this.queMin.Len() > this.queMax.Len() {
		return float64(this.queMin.Ints[0])
	}
	return float64(this.queMin.Ints[0]+this.queMax.Ints[0]) / 2.0
}

// 辅助结构：高性能切片堆
type hp struct {
	Ints    []int
	compare func(a, b int) bool
}

func (h *hp) Len() int { return len(h.Ints) }
func (h *hp) push(v int) {
	h.Ints = append(h.Ints, v)
	h.up(h.Len() - 1)
}
func (h *hp) pop() int {
	v := h.Ints[0]
	h.Ints[0] = h.Ints[h.Len()-1]
	h.Ints = h.Ints[:h.Len()-1]
	h.down(0)
	return v
}
func (h *hp) up(i int) {
	for {
		j := (i - 1) / 2
		if j < 0 || !h.compare(h.Ints[i], h.Ints[j]) {
			break
		}
		h.Ints[i], h.Ints[j] = h.Ints[j], h.Ints[i]
		i = j
	}
}
func (h *hp) down(i int) {
	for {
		l := 2*i + 1
		if l >= h.Len() {
			break
		}
		j := l
		if r := l + 1; r < h.Len() && h.compare(h.Ints[r], h.Ints[l]) {
			j = r
		}
		if !h.compare(h.Ints[j], h.Ints[i]) {
			break
		}
		h.Ints[i], h.Ints[j] = h.Ints[j], h.Ints[i]
		i = j
	}
}

func TestMedianFinder(t *testing.T) {
	// 测试中位数查找逻辑
}

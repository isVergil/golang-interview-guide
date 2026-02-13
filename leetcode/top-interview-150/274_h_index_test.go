package topinterview150

import "testing"

// 274. H 指数 (H-Index)
//
// 题目描述:
// 给你一个整数数组 citations ，其中 citations[i] 表示研究者的第 i 篇论文被引用的次数。计算并返回该研究者的 h 指数。
// 根据维基百科上 h 指数的定义：h 代表“高引用次数” ，一名科研人员的 h 指数是指他（她）至少发表了 h 篇论文，
// 并且 至少 有 h 篇论文被引用次数大于等于 h 。如果 h 有多种可能的值，h 指数 是其中最大的那个。
//
// 示例 1：
// 输入：citations = [3,0,6,1,5]
// 输出：3
// 解释：给定数组表示研究者总共有 5 篇论文，每篇论文相应的被引用了 3, 0, 6, 1, 5 次。
//     由于研究者有 3 篇论文每篇 至少 被引用了 3 次，其余 2 篇论文每篇被引用 不多于 3 次，所以她的 h 指数是 3。
//
// 示例 2：
// 输入：citations = [1,3,1]
// 输出：1

func hIndex(citations []int) int {
	// 在这里写入你的代码
	// 提示：可以排序后遍历，或者使用计数排序 (Bucket Sort) 优化。
	panic("请实现该函数")
}

func TestHIndex(t *testing.T) {
	tests := []struct {
		name      string
		citations []int
		expected  int
	}{
		{
			name:      "Example 1",
			citations: []int{3, 0, 6, 1, 5},
			expected:  3,
		},
		{
			name:      "Example 2",
			citations: []int{1, 3, 1},
			expected:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: Uncomment when implementation is ready
			// if got := hIndex(tt.citations); got != tt.expected {
			// 	t.Errorf("hIndex() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

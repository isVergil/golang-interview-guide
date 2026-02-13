package topinterview150

import (
	"testing"
)

// 128. 最长连续序列 (Longest Consecutive Sequence)
//
// 题目描述:
// 给定一个未排序的整数数组 nums ，找出数字连续的最长序列（不要求序列元素在原数组中连续）的长度。
// 请你设计并实现时间复杂度为 O(n) 的算法解决此问题。
//
// 示例 1：
// 输入：nums = [100,4,200,1,3,2]
// 输出：4
// 解释：最长数字连续序列是 [1, 2, 3, 4]。它的长度为 4。
//
// 示例 2：
// 输入：nums = [0,3,7,2,5,8,4,6,0,1]
// 输出：9

func longestConsecutive(nums []int) int {
	numSet := make(map[int]bool)
	for _, num := range nums {
		numSet[num] = true
	}

	res := 0
	for num := range numSet {
		if !numSet[num-1] {
			curNum := num
			curStart := 1
			for numSet[curNum+1] {
				curNum++
				curStart++
			}
			res = max(curStart, res)
		}
	}
	return res
}

func TestLongestConsecutive(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected int
	}{
		{
			name:     "Example 1",
			nums:     []int{100, 4, 200, 1, 3, 2},
			expected: 4,
		},
		{
			name:     "Example 2",
			nums:     []int{0, 3, 7, 2, 5, 8, 4, 6, 0, 1},
			expected: 9,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: Uncomment when implementation is ready
			// if got := longestConsecutive(tt.nums); got != tt.expected {
			// 	t.Errorf("longestConsecutive() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

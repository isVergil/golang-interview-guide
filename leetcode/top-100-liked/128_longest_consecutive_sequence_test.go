package top100liked

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
	// 将所有元素放入哈希集合，用于 O(1) 查找
	set := make(map[int]bool, len(nums))
	for _, n := range nums {
		set[n] = true
	}

	res := 0
	for n := range set {
		// 只从连续序列的起点开始计数，避免重复遍历
		if set[n-1] {
			continue
		} // 从起点向后延伸，统计连续长度
		length := 1
		for set[n+length] {
			length++
		} // 更新最长长度
		if length > res {
			res = length
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
			name:     "示例1",
			nums:     []int{100, 4, 200, 1, 3, 2},
			expected: 4,
		},
		{
			name:     "示例2",
			nums:     []int{0, 3, 7, 2, 5, 8, 4, 6, 0, 1},
			expected: 9,
		},
		{
			name:     "空数组",
			nums:     []int{},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := longestConsecutive(tt.nums)
			if got != tt.expected {
				t.Errorf("longestConsecutive() = %v, want %v", got, tt.expected)
			}
		})
	}
}

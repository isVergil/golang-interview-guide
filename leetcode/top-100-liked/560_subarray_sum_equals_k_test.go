package top100liked

import (
	"testing"
)

// 560. 和为 K 的子数组 (Subarray Sum Equals K)
//
// 题目描述:
// 给你一个整数数组 nums 和一个整数 k ，请你统计并返回该数组中和为 k 的子数组的个数。
// 子数组是数组中元素的连续非空序列。
//
// 示例 1：
// 输入：nums = [1,1,1], k = 2
// 输出：2
//
// 示例 2：
// 输入：nums = [1,2,3], k = 3
// 输出：2

func subarraySum(nums []int, k int) int {
	var prefixCount = map[int]int{0: 1}
	sum, res := 0, 0
	for _, num := range nums {
		sum += num
		res += prefixCount[sum-k]
		prefixCount[sum]++
	}
	return res
}

func TestSubarraySum(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		k        int
		expected int
	}{
		{
			name:     "示例1",
			nums:     []int{1, 1, 1},
			k:        2,
			expected: 2,
		},
		{
			name:     "示例2",
			nums:     []int{1, 2, 3},
			k:        3,
			expected: 2,
		},
		{
			name:     "含负数",
			nums:     []int{1, -1, 0},
			k:        0,
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := subarraySum(tt.nums, tt.k)
			if got != tt.expected {
				t.Errorf("subarraySum() = %v, want %v", got, tt.expected)
			}
		})
	}
}

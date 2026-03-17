package top100liked

import (
	"testing"
)

// 560. 和为 K 的子数组 (Subarray Sum Equals K)
//
// 题目描述:
// 给你一个整数数组 nums 和一个整数 k ，请你统计并返回 该数组中和为 k 的子数组的个数 。
// 子数组是数组中元素的连续序列。
//
// 示例 1：
// 输入：nums = [1,1,1], k = 2
// 输出：2
//
// 示例 2：
// 输入：nums = [1,2,3], k = 3
// 输出：2

func subarraySum(nums []int, k int) int {
	panic("not implemented")
}

func TestSubarraySum(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		k        int
		expected int
	}{
		{"Example 1", []int{1, 1, 1}, 2, 2},
		{"Example 2", []int{1, 2, 3}, 3, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := subarraySum(tt.nums, tt.k); got != tt.expected {
			// 	t.Errorf("subarraySum() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

package top100liked

import (
	"testing"
)

// 53. 最大子数组和 (Maximum Subarray)
//
// 题目描述:
// 给你一个整数数组 nums ，请你找出一个具有最大和的连续子数组（子数组最少包含一个元素），
// 返回其最大和。子数组是数组中的一个连续部分。
//
// 示例 1：
// 输入：nums = [-2,1,-3,4,-1,2,1,-5,4]
// 输出：6
// 解释：连续子数组 [4,-1,2,1] 的和最大，为 6。
//
// 示例 2：
// 输入：nums = [1]
// 输出：1
//
// 示例 3：
// 输入：nums = [5,4,-1,7,8]
// 输出：23

func maxSubArray(nums []int) int {
	res, cur := nums[0], nums[0]
	for i := 1; i < len(nums); i++ {
		if cur > 0 {
			cur += nums[i]
		} else {
			cur = nums[i]
		}
		if cur > res {
			res = cur
		}
	}
	return res
}

func TestMaxSubArray(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected int
	}{
		{
			name:     "示例1",
			nums:     []int{-2, 1, -3, 4, -1, 2, 1, -5, 4},
			expected: 6,
		},
		{
			name:     "示例2",
			nums:     []int{1},
			expected: 1,
		},
		{
			name:     "示例3",
			nums:     []int{5, 4, -1, 7, 8},
			expected: 23,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := maxSubArray(tt.nums)
			if got != tt.expected {
				t.Errorf("maxSubArray() = %v, want %v", got, tt.expected)
			}
		})
	}
}

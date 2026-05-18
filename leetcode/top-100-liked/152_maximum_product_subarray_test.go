package top100liked

import "testing"

// 152. 乘积最大子数组 (Maximum Product Subarray)
//
// 题目描述:
// 给你一个整数数组 nums，请你找出数组中乘积最大的非空连续子数组，并返回该子数组所对应的乘积。
//
// 示例 1：
// 输入：nums = [2,3,-2,4]
// 输出：6
// 解释：子数组 [2,3] 有最大乘积 6
//
// 示例 2：
// 输入：nums = [-2,0,-1]
// 输出：0
//
// 提示：同时维护 curMax 和 curMin，遇到负数时交换（负负得正）

func maxProduct(nums []int) int {
	dpMax, dpMin, res := nums[0], nums[0], nums[0]
	for i := 1; i < len(nums); i++ {
		a := dpMax * nums[i]
		b := dpMin * nums[i]
		c := nums[i]
		dpMax = max(a, max(b, c))
		dpMin = min(a, min(b, c))
		res = max(res, dpMax)
	}
	return res
}

func TestMaxProduct(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected int
	}{
		{name: "示例1", nums: []int{2, 3, -2, 4}, expected: 6},
		{name: "示例2", nums: []int{-2, 0, -1}, expected: 0},
		{name: "全负数", nums: []int{-2, -3, -4}, expected: 12},
		{name: "单个负数", nums: []int{-2}, expected: -2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := maxProduct(tt.nums)
			if got != tt.expected {
				t.Errorf("maxProduct() = %v, want %v", got, tt.expected)
			}
		})
	}
}

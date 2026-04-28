package top100liked

import (
	"testing"
)

// 416. 分割等和子集 (Partition Equal Subset Sum)
//
// 题目描述:
// 给你一个只包含正整数的非空数组 nums 。请你判断是否可以将这个数组分割成两个子集，
// 使得两个子集的元素和相等。
//
// 示例 1：
// 输入：nums = [1,5,11,5]
// 输出：true（数组可以分割成 [1, 5, 5] 和 [11]）
//
// 示例 2：
// 输入：nums = [1,2,3,5]
// 输出：false（数组不能分割成两个元素和相等的子集）

func canPartition(nums []int) bool {
	sum, maxVal := 0, 0
	for _, v := range nums {
		sum += v
		maxVal = max(v, maxVal)
	}
	if sum%2 != 0 {
		return false
	}

	target := sum / 2
	if maxVal > target {
		return false
	}

	// dp[j] 表示是否能恰好凑出和为 j
	dp := make([]bool, target+1)
	dp[0] = true
	for _, num := range nums {
		for j := target; j >= num; j-- {
			dp[j] = dp[j] || dp[j-num]
		}
		if dp[target] {
			return true
		}
	}

	return false

}

func TestCanPartition(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected bool
	}{
		{name: "示例1", nums: []int{1, 5, 11, 5}, expected: true},
		{name: "示例2", nums: []int{1, 2, 3, 5}, expected: false},
		{name: "两个相同", nums: []int{1, 1}, expected: true},
		{name: "单元素", nums: []int{1}, expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := canPartition(tt.nums)
			if got != tt.expected {
				t.Errorf("canPartition() = %v, want %v", got, tt.expected)
			}
		})
	}
}

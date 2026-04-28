package top100liked

import (
	"testing"
)

// 300. 最长递增子序列 (Longest Increasing Subsequence)
//
// 题目描述:
// 给你一个整数数组 nums ，找到其中最长严格递增子序列的长度。
// 子序列是由数组派生而来的序列，删除（或不删除）数组中的元素而不改变其余元素的顺序。
//
// 示例 1：
// 输入：nums = [10,9,2,5,3,7,101,18]
// 输出：4（最长递增子序列是 [2,3,7,101]）
//
// 示例 2：
// 输入：nums = [0,1,0,3,2,3]
// 输出：4
//
// 示例 3：
// 输入：nums = [7,7,7,7,7,7,7]
// 输出：1

// 动态规划
func lengthOfLIS1(nums []int) int {
	n := len(nums)
	dp := make([]int, n)
	res := 1
	for i := range dp {
		dp[i] = 1
	}
	for i := 1; i < n; i++ {
		for j := 0; j < i; j++ {
			if nums[i] > nums[j] {
				dp[i] = max(dp[i], dp[j]+1)
			}
		}
		res = max(res, dp[i])
	}
	return res
}

// 贪心+二分
func lengthOfLIS2(nums []int) int {
	// 维护一个递增序列，每个元素代表能组成当前长度的最小元素
	tail := make([]int, 0)
	tail = append(tail, nums[0])
	for i := 1; i < len(nums); i++ {
		if nums[i] > tail[len(tail)-1] {
			tail = append(tail, nums[i])
		} else {
			// 二分第一个 > nums[i] 并替换
			l, r := 0, len(tail)-1
			for l < r {
				mid := l + (r-l)/2
				if tail[mid] >= nums[i] {
					r = mid
				} else {
					l = mid + 1
				}
			}
			tail[l] = nums[i]
		}
	}
	return len(tail)
}

func TestLengthOfLIS(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected int
	}{
		{name: "示例1", nums: []int{10, 9, 2, 5, 3, 7, 101, 18}, expected: 4},
		{name: "示例2", nums: []int{0, 1, 0, 3, 2, 3}, expected: 4},
		{name: "示例3", nums: []int{7, 7, 7, 7, 7, 7, 7}, expected: 1},
		{name: "单元素", nums: []int{1}, expected: 1},
		{name: "递增", nums: []int{1, 2, 3, 4, 5}, expected: 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := lengthOfLIS1(tt.nums)
			if got != tt.expected {
				t.Errorf("lengthOfLIS1() = %v, want %v", got, tt.expected)
			}
		})

		t.Run(tt.name, func(t *testing.T) {
			got := lengthOfLIS2(tt.nums)
			if got != tt.expected {
				t.Errorf("lengthOfLIS2() = %v, want %v", got, tt.expected)
			}
		})
	}
}

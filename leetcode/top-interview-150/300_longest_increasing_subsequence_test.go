package topinterview150

import (
	"testing"
)

// 300. 最长递增子序列 (Longest Increasing Subsequence)
//
// 题目描述:
// 给你一个整数数组 nums ，找到其中最长严格递增子序列的长度。
// 子序列 是由数组派生而来的序列，删除（或不删除）数组中的元素而不改变其余元素的顺序。例如，[3,6,2,7] 是数组 [0,3,1,6,2,2,7] 的子序列。
//
// 示例 1：
// 输入：nums = [10,9,2,5,3,7,101,18]
// 输出：4
// 解释：最长递增子序列是 [2,3,7,101]，因此长度为 4 。
//
// 示例 2：
// 输入：nums = [0,1,0,3,2,3]
// 输出：4
//
// 示例 3：
// 输入：nums = [7,7,7,7,7,7,7]
// 输出：1

func lengthOfLIS(nums []int) int {
	n := len(nums)
	if n == 0 {
		return 0
	}

	dp := make([]int, n)
	res := 1

	for i := 0; i < n; i++ {
		// 初始化
		dp[i] = 1

		for j := 0; j < i; j++ {
			if nums[i] > nums[j] {
				if dp[j]+1 > dp[i] {
					dp[i] = dp[j] + 1
				}
			}
		}

		if dp[i] > res {
			res = dp[i]
		}
	}

	return res
}

func TestLengthOfLIS(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected int
	}{
		{"Example 1", []int{10, 9, 2, 5, 3, 7, 101, 18}, 4},
		{"Example 2", []int{0, 1, 0, 3, 2, 3}, 4},
		{"Example 3", []int{7, 7, 7, 7, 7, 7, 7}, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := lengthOfLIS(tt.nums); got != tt.expected {
			// 	t.Errorf("lengthOfLIS() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

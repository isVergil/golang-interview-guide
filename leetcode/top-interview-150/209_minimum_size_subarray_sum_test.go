package topinterview150

import (
	"testing"
)

// 209. 长度最小的子数组 (Minimum Size Subarray Sum)
//
// 题目描述:
// 给定一个含有 n 个正整数的数组和一个正整数 target 。
// 找出该数组中满足其和 ≥ target 的长度最小的 连续子数组 [numsl, numsl+1, ..., numsr-1, numsr] ，并返回其长度。如果不存在符合条件的子数组，返回 0 。
//
// 示例 1：
// 输入：target = 7, nums = [2,3,1,2,4,3]
// 输出：2
// 解释：子数组 [4,3] 是该条件下的长度最小的子数组。
//
// 示例 2：
// 输入：target = 4, nums = [1,4,4]
// 输出：1
//
// 示例 3：
// 输入：target = 11, nums = [1,1,1,1,1,1,1,1]
// 输出：0

func minSubArrayLen(target int, nums []int) int {
	panic("not implemented")
}

func TestMinSubArrayLen(t *testing.T) {
	tests := []struct {
		name     string
		target   int
		nums     []int
		expected int
	}{
		{"Example 1", 7, []int{2, 3, 1, 2, 4, 3}, 2},
		{"Example 2", 4, []int{1, 4, 4}, 1},
		{"Example 3", 11, []int{1, 1, 1, 1, 1, 1, 1, 1}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := minSubArrayLen(tt.target, tt.nums); got != tt.expected {
			// 	t.Errorf("minSubArrayLen() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

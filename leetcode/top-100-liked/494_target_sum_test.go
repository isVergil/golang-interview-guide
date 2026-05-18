package top100liked

import (
	"math"
	"testing"
)

// 494. 目标和 (Target Sum)
//
// 题目描述:
// 给你一个非负整数数组 nums 和一个整数 target。
// 向数组中的每个整数前添加 '+' 或 '-'，然后串联起所有整数，可以构造一个表达式。
// 返回可以通过上述方法构造的、运算结果等于 target 的不同表达式的数目。
//
// 示例 1：
// 输入：nums = [1,1,1,1,1], target = 3
// 输出：5
// 解释：
// -1 + 1 + 1 + 1 + 1 = 3
// +1 - 1 + 1 + 1 + 1 = 3
// +1 + 1 - 1 + 1 + 1 = 3
// +1 + 1 + 1 - 1 + 1 = 3
// +1 + 1 + 1 + 1 - 1 = 3
//
// 提示：转化为 01 背包问题，选一部分数使其和为 (sum + target) / 2

func findTargetSumWays(nums []int, target int) int {
	sum := 0
	for _, v := range nums {
		sum += v
	}

	if math.Abs(float64(sum)) > float64(sum) || (target+sum)%2 != 0 {
		return 0
	}
	cap := (target + sum) / 2

	// dp[j] 和为 j 的方案数
	dp := make([]int, cap+1)
	dp[0] = 1
	for _, v := range nums {
		for j := cap; j >= v; j-- {
			dp[j] += dp[j-v]
		}
	}
	return dp[cap]
}

func TestFindTargetSumWays(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		target   int
		expected int
	}{
		{name: "示例1", nums: []int{1, 1, 1, 1, 1}, target: 3, expected: 5},
		{name: "单元素", nums: []int{1}, target: 1, expected: 1},
		{name: "不可能", nums: []int{1}, target: 2, expected: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findTargetSumWays(tt.nums, tt.target)
			if got != tt.expected {
				t.Errorf("findTargetSumWays() = %v, want %v", got, tt.expected)
			}
		})
	}
}

package top100liked

import (
	"testing"
)

// 198. 打家劫舍 (House Robber)
//
// 题目描述:
// 你是一个专业的小偷，计划偷窃沿街的房屋。每间房内都藏有一定的现金，
// 影响你偷窃的唯一制约因素就是相邻的房屋装有相互连通的防盗系统，
// 如果两间相邻的房屋在同一晚上被小偷闯入，系统会自动报警。
// 给定一个代表每个房屋存放金额的非负整数数组，计算你不触动警报装置的情况下，一夜之内能够偷窃到的最高金额。
//
// 示例 1：
// 输入：[1,2,3,1]
// 输出：4（偷窃 1 号房屋(金额=1)和 3 号房屋(金额=3)，偷窃金额 = 1 + 3 = 4）
//
// 示例 2：
// 输入：[2,7,9,3,1]
// 输出：12（偷窃 1、3、5 号房屋，偷窃金额 = 2 + 9 + 1 = 12）

func rob(nums []int) int {
	prev, cur := 0, 0
	for _, num := range nums {
		prev, cur = cur, max(cur, num+prev)
	}
	return cur
}

func TestRob(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected int
	}{
		{name: "示例1", nums: []int{1, 2, 3, 1}, expected: 4},
		{name: "示例2", nums: []int{2, 7, 9, 3, 1}, expected: 12},
		{name: "单个房屋", nums: []int{5}, expected: 5},
		{name: "两个房屋", nums: []int{1, 2}, expected: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rob(tt.nums)
			if got != tt.expected {
				t.Errorf("rob() = %v, want %v", got, tt.expected)
			}
		})
	}
}

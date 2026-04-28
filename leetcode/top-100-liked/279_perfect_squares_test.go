package top100liked

import (
	"testing"
)

// 279. 完全平方数 (Perfect Squares)
//
// 题目描述:
// 给你一个整数 n ，返回和为 n 的完全平方数的最少数量。
// 完全平方数是一个整数，其值等于另一个整数的平方；即其值等于一个整数自乘的积。
// 例如，1、4、9 和 16 都是完全平方数，而 3 和 11 不是。
//
// 示例 1：
// 输入：n = 12
// 输出：3（12 = 4 + 4 + 4）
//
// 示例 2：
// 输入：n = 13
// 输出：2（13 = 4 + 9）

// 动态规划：dp[i] 表示组成 i 的最少完全平方数个数
// 对每个 i，尝试所有 j*j <= i 的完全平方数，取最小值
// 时间 O(n√n)，空间 O(n)
func numSquares(n int) int {
	dp := make([]int, n+1)
	for i := 1; i <= n; i++ {
		dp[i] = i // 最坏情况：全用 1（1+1+...+1）
		for j := 1; j*j <= i; j++ {
			// 用一个 j*j，剩下的是 dp[i-j*j]
			dp[i] = min(dp[i], dp[i-j*j]+1)
		}
	}
	return dp[n]
}

func TestNumSquares(t *testing.T) {
	tests := []struct {
		name     string
		n        int
		expected int
	}{
		{name: "示例1", n: 12, expected: 3},
		{name: "示例2", n: 13, expected: 2},
		{name: "完全平方数", n: 4, expected: 1},
		{name: "n=1", n: 1, expected: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := numSquares(tt.n)
			if got != tt.expected {
				t.Errorf("numSquares(%d) = %v, want %v", tt.n, got, tt.expected)
			}
		})
	}
}

package top100liked

import (
	"testing"
)

// 62. 不同路径 (Unique Paths)
//
// 题目描述:
// 一个机器人位于一个 m x n 网格的左上角。机器人每次只能向下或者向右移动一步。
// 机器人试图达到网格的右下角。问总共有多少条不同的路径？
//
// 示例 1：
// 输入：m = 3, n = 7
// 输出：28
//
// 示例 2：
// 输入：m = 3, n = 2
// 输出：3
// 解释：从左上角开始，总共有 3 条路径可以到达右下角。
// 1. 向右 -> 向下 -> 向下
// 2. 向下 -> 向下 -> 向右
// 3. 向下 -> 向右 -> 向下
// 回溯
func uniquePaths(m int, n int) int {
	res := 0
	var backtrack func(i, j int)
	backtrack = func(i, j int) {
		if i == m-1 && j == n-1 {
			res++
			return
		}
		if i >= m || j >= n {
			return
		}
		backtrack(i, j+1) // 向右
		backtrack(i+1, j) // 向下
	}
	backtrack(0, 0)
	return res
}

// DP 到达当前位置的步数 = 左侧 + 上侧
// 从左到右更新时，还没更新的是上一行的值（上方），已更新的是当前行的值（左方），刚好就是二维 DP 的两个依赖。
func uniquePaths1(m int, n int) int {
	dp := make([]int, n)
	for i := 0; i < n; i++ {
		dp[i] = 1
	}

	for i := 1; i < m; i++ {
		for j := 1; j < n; j++ {
			dp[j] += dp[j-1]
		}
	}
	return dp[n-1]
}

func TestUniquePaths(t *testing.T) {
	tests := []struct {
		name     string
		m        int
		n        int
		expected int
	}{
		{
			name:     "示例1",
			m:        3,
			n:        7,
			expected: 28,
		},
		{
			name:     "示例2",
			m:        3,
			n:        2,
			expected: 3,
		},
		{
			name:     "1x1",
			m:        1,
			n:        1,
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := uniquePaths(tt.m, tt.n)
			if got != tt.expected {
				t.Errorf("uniquePaths() = %v, want %v", got, tt.expected)
			}

			got1 := uniquePaths1(tt.m, tt.n)
			if got1 != tt.expected {
				t.Errorf("uniquePaths1() = %v, want %v", got1, tt.expected)
			}
		})
	}
}

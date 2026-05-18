package top100liked

import "testing"

// 64. 最小路径和 (Minimum Path Sum)
//
// 题目描述:
// 给定一个包含非负整数的 m x n 网格 grid，请找出一条从左上角到右下角的路径，
// 使得路径上的数字总和为最小。每次只能向下或者向右移动一步。
//
// 示例 1：
// 输入：grid = [[1,3,1],[1,5,1],[4,2,1]]
// 输出：7
// 解释：路径 1→3→1→1→1 的总和最小
//
// 示例 2：
// 输入：grid = [[1,2,3],[4,5,6]]
// 输出：12
//
// 提示：dp[i][j] = min(dp[i-1][j], dp[i][j-1]) + grid[i][j]

func minPathSum(grid [][]int) int {
	m, n := len(grid), len(grid[0])
	dp := make([][]int, m)
	for i := 0; i < m; i++ {
		dp[i] = make([]int, n)
	}
	dp[0][0] = grid[0][0]
	for i := 1; i < m; i++ {
		dp[i][0] = dp[i-1][0] + grid[i][0]
	}
	for i := 1; i < n; i++ {
		dp[0][i] = dp[0][i-1] + grid[0][i]
	}
	for i := 1; i < m; i++ {
		for j := 1; j < n; j++ {
			dp[i][j] = min(dp[i-1][j], dp[i][j-1]) + grid[i][j]
		}
	}
	return dp[m-1][n-1]
}

// 空间优化
func minPathSum1(grid [][]int) int {
	m, n := len(grid), len(grid[0])
	dp := make([]int, n)

	// 初始化第一行
	dp[0] = grid[0][0]
	for i := 1; i < n; i++ {
		dp[i] = dp[i-1] + grid[0][i]
	}
	for i := 1; i < m; i++ {
		dp[0] += grid[i][0]
		for j := 1; j < n; j++ {
			dp[j] = min(dp[j-1], dp[j]) + grid[i][j]
		}
	}
	return dp[n-1]
}

func TestMinPathSum(t *testing.T) {
	tests := []struct {
		name     string
		grid     [][]int
		expected int
	}{
		{name: "示例1", grid: [][]int{{1, 3, 1}, {1, 5, 1}, {4, 2, 1}}, expected: 7},
		{name: "示例2", grid: [][]int{{1, 2, 3}, {4, 5, 6}}, expected: 12},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := minPathSum(tt.grid)
			if got != tt.expected {
				t.Errorf("minPathSum() = %v, want %v", got, tt.expected)
			}
		})

		t.Run(tt.name, func(t *testing.T) {
			got := minPathSum1(tt.grid)
			if got != tt.expected {
				t.Errorf("minPathSum1() = %v, want %v", got, tt.expected)
			}
		})
	}
}

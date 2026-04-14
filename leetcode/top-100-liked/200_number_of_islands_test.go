package top100liked

import (
	"testing"
)

// 200. 岛屿数量 (Number of Islands)
//
// 题目描述:
// 给你一个由 '1'（陆地）和 '0'（水）组成的的二维网格，请你计算网格中岛屿的数量。
// 岛屿总是被水包围，并且每座岛屿只能由水平方向和/或垂直方向上相邻的陆地连接形成。
// 此外，你可以假设该网格的四条边均被水包围。
//
// 示例 1：
// 输入：grid = [
//   ["1","1","1","1","0"],
//   ["1","1","0","1","0"],
//   ["1","1","0","0","0"],
//   ["0","0","0","0","0"]
// ]
// 输出：1
//
// 示例 2：
// 输入：grid = [
//   ["1","1","0","0","0"],
//   ["1","1","0","0","0"],
//   ["0","0","1","0","0"],
//   ["0","0","0","1","1"]
// ]
// 输出：3

func numIslands(grid [][]byte) int {
	res := 0
	for i := range grid {
		for j := range grid[i] {
			if grid[i][j] == '1' {
				res++
				dfs(grid, i, j)
			}
		}
	}
	return res
}

func dfs(grid [][]byte, i, j int) {
	// 越界或者不是陆地
	if i < 0 || i >= len(grid) || j < 0 || j >= len(grid[0]) || grid[i][j] != '1' {
		return
	}

	// 标记为已访问（沉岛），避免重复遍历
	grid[i][j] = '0'
	// 向四个方向扩散
	dfs(grid, i+1, j)
	dfs(grid, i-1, j)
	dfs(grid, i, j+1)
	dfs(grid, i, j-1)
}

func TestNumIslands(t *testing.T) {
	tests := []struct {
		name     string
		grid     [][]byte
		expected int
	}{
		{
			name: "示例1",
			grid: [][]byte{
				{'1', '1', '1', '1', '0'},
				{'1', '1', '0', '1', '0'},
				{'1', '1', '0', '0', '0'},
				{'0', '0', '0', '0', '0'},
			},
			expected: 1,
		},
		{
			name: "示例2",
			grid: [][]byte{
				{'1', '1', '0', '0', '0'},
				{'1', '1', '0', '0', '0'},
				{'0', '0', '1', '0', '0'},
				{'0', '0', '0', '1', '1'},
			},
			expected: 3,
		},
		{
			name:     "空网格",
			grid:     [][]byte{},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 深拷贝 grid，避免修改原测试数据
			grid := make([][]byte, len(tt.grid))
			for i := range tt.grid {
				grid[i] = make([]byte, len(tt.grid[i]))
				copy(grid[i], tt.grid[i])
			}
			got := numIslands(grid)
			if got != tt.expected {
				t.Errorf("numIslands() = %v, want %v", got, tt.expected)
			}
		})
	}
}

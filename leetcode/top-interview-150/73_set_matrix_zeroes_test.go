package topinterview150

import (
	"testing"
)

// 73. 矩阵置零 (Set Matrix Zeroes)
//
// 题目描述:
// 给定一个 m x n 的矩阵，如果一个元素为 0 ，则将其所在行和列的所有元素都设为 0 。请使用 原地 算法。
//
// 示例 1：
// 输入：matrix = [[1,1,1],[1,0,1],[1,1,1]]
// 输出：[[1,0,1],[0,0,0],[1,0,1]]
//
// 示例 2：
// 输入：matrix = [[0,1,2,0],[3,4,5,2],[1,3,1,5]]
// 输出：[[0,0,0,0],[0,4,5,0],[0,3,1,0]]

func setZeroes(matrix [][]int) {
	m, n := len(matrix), len(matrix[0])

	row0HasZero, col0HasZero := false, false

	// 1. 检查第一行和第一列本身是否含有 0
	for j := 0; j < n; j++ {
		if matrix[0][j] == 0 {
			row0HasZero = true
			break
		}
	}
	for i := 0; i < m; i++ {
		if matrix[i][0] == 0 {
			col0HasZero = true
			break
		}
	}

	// 2. 用第一行和第一列记录其他格子的 0 情况
	for i := 1; i < m; i++ {
		for j := 1; j < n; j++ {
			if matrix[i][j] == 0 {
				matrix[i][0] = 0
				matrix[0][j] = 0
			}
		}
	}

	// 3. 根据标记位，将非首行首列的元素置零
	for i := 1; i < m; i++ {
		for j := 1; j < n; j++ {
			if matrix[i][0] == 0 || matrix[0][j] == 0 {
				matrix[i][j] = 0
			}
		}
	}

	// 4. 最后单独处理第一行和第一列
	if row0HasZero {
		for j := 0; j < n; j++ {
			matrix[0][j] = 0
		}
	}
	if col0HasZero {
		for i := 0; i < m; i++ {
			matrix[i][0] = 0
		}
	}

}

func TestSetZeroes(t *testing.T) {
	tests := []struct {
		name     string
		matrix   [][]int
		expected [][]int
	}{
		{
			"Example 1",
			[][]int{{1, 1, 1}, {1, 0, 1}, {1, 1, 1}},
			[][]int{{1, 0, 1}, {0, 0, 0}, {1, 0, 1}},
		},
		{
			"Example 2",
			[][]int{{0, 1, 2, 0}, {3, 4, 5, 2}, {1, 3, 1, 5}},
			[][]int{{0, 0, 0, 0}, {0, 4, 5, 0}, {0, 3, 1, 0}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setZeroes(tt.matrix)
			// if !reflect.DeepEqual(tt.matrix, tt.expected) {
			// 	t.Errorf("setZeroes() = %v, want %v", tt.matrix, tt.expected)
			// }
		})
	}
}

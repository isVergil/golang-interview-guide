package top100liked

import (
	"reflect"
	"testing"
)

// 73. 矩阵置零 (Set Matrix Zeroes)
//
// 题目描述:
// 给定一个 m x n 的矩阵，如果一个元素为 0，则将其所在行和列的所有元素都设为 0。请使用原地算法。
//
// 示例 1：
// 输入：matrix = [[1,1,1],[1,0,1],[1,1,1]]
// 输出：[[1,0,1],[0,0,0],[1,0,1]]
//
// 示例 2：
// 输入：matrix = [[0,1,2,0],[3,4,5,2],[1,3,1,5]]
// 输出：[[0,0,0,0],[0,4,5,0],[0,3,1,0]]
//
// 提示：用第一行和第一列作为标记数组，O(1) 空间

func setZeroes(matrix [][]int) {
	m, n := len(matrix), len(matrix[0])

	// 第一行列本身是否有 0
	firstRowZero, firstColZero := false, false
	for i := 0; i < m; i++ {
		if matrix[i][0] == 0 {
			firstColZero = true
			break
		}
	}
	for i := 0; i < n; i++ {
		if matrix[0][i] == 0 {
			firstRowZero = true
			break
		}
	}

	// 扫描行列，用第一行第一列做标记
	for i := 1; i < m; i++ {
		for j := 1; j < n; j++ {
			if matrix[i][j] == 0 {
				matrix[i][0] = 0
				matrix[0][j] = 0
			}
		}
	}

	// 置零
	for i := 1; i < m; i++ {
		for j := 1; j < n; j++ {
			if matrix[i][0] == 0 || matrix[0][j] == 0 {
				matrix[i][j] = 0
			}
		}
	}

	// 处理第一行
	if firstRowZero {
		for i := 0; i < n; i++ {
			matrix[0][i] = 0
		}
	}

	// 处理第一列
	if firstColZero {
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
			name:     "示例1",
			matrix:   [][]int{{1, 1, 1}, {1, 0, 1}, {1, 1, 1}},
			expected: [][]int{{1, 0, 1}, {0, 0, 0}, {1, 0, 1}},
		},
		{
			name:     "示例2",
			matrix:   [][]int{{0, 1, 2, 0}, {3, 4, 5, 2}, {1, 3, 1, 5}},
			expected: [][]int{{0, 0, 0, 0}, {0, 4, 5, 0}, {0, 3, 1, 0}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setZeroes(tt.matrix)
			if !reflect.DeepEqual(tt.matrix, tt.expected) {
				t.Errorf("setZeroes() = %v, want %v", tt.matrix, tt.expected)
			}
		})
	}
}

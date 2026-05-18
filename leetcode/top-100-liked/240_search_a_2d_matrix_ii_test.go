package top100liked

import "testing"

// 240. 搜索二维矩阵 II (Search a 2D Matrix II)
//
// 题目描述:
// 编写一个高效的算法来搜索 m x n 矩阵 matrix 中的一个目标值 target。
// 该矩阵具有以下特性：
// - 每行的元素从左到右升序排列
// - 每列的元素从上到下升序排列
//
// 示例 1：
// 输入：matrix = [[1,4,7,11,15],[2,5,8,12,19],[3,6,9,16,22],[10,13,14,17,24],[18,21,23,26,30]], target = 5
// 输出：true
//
// 示例 2：
// 同上矩阵，target = 20
// 输出：false
//
// 提示：从右上角开始搜索，大了往左走，小了往下走

func searchMatrix240(matrix [][]int, target int) bool {
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		return false
	}

	// 从右上角往左下找
	m, n := len(matrix), len(matrix[0])
	row, col := 0, n-1
	for row < m && col >= 0 {
		if matrix[row][col] == target {
			return true
		} else if matrix[row][col] > target {
			col--
		} else {
			row++
		}
	}
	return false
}

func TestSearchMatrix240(t *testing.T) {
	matrix := [][]int{
		{1, 4, 7, 11, 15},
		{2, 5, 8, 12, 19},
		{3, 6, 9, 16, 22},
		{10, 13, 14, 17, 24},
		{18, 21, 23, 26, 30},
	}

	tests := []struct {
		name     string
		target   int
		expected bool
	}{
		{name: "存在-5", target: 5, expected: true},
		{name: "不存在-20", target: 20, expected: false},
		{name: "存在-边界值", target: 30, expected: true},
		{name: "存在-左上角", target: 1, expected: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := searchMatrix240(matrix, tt.target)
			if got != tt.expected {
				t.Errorf("searchMatrix240() = %v, want %v", got, tt.expected)
			}
		})
	}
}

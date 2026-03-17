package topinterview150

import (
	"testing"
)

// 74. 搜索二维矩阵 (Search a 2D Matrix)
//
// 题目描述:
// 给你一个满足下述两条属性的 m x n 整数矩阵：
// 每行中的整数从左到右按非严格递增顺序排列。
// 每行的第一个整数大于前一行的最后一个整数。
// 给你一个整数 target ，如果 target 在矩阵中，返回 true ；否则，返回 false 。
//
// 示例 1：
// 输入：matrix = [[1,3,5,7],[10,11,16,20],[23,30,34,60]], target = 3
// 输出：true
//
// 示例 2：
// 输入：matrix = [[1,3,5,7],[10,11,16,20],[23,30,34,60]], target = 13
// 输出：false

func searchMatrix(matrix [][]int, target int) bool {
	panic("not implemented")
}

func TestSearchMatrix(t *testing.T) {
	tests := []struct {
		name     string
		matrix   [][]int
		target   int
		expected bool
	}{
		{"Example 1", [][]int{{1, 3, 5, 7}, {10, 11, 16, 20}, {23, 30, 34, 60}}, 3, true},
		{"Example 2", [][]int{{1, 3, 5, 7}, {10, 11, 16, 20}, {23, 30, 34, 60}}, 13, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := searchMatrix(tt.matrix, tt.target); got != tt.expected {
			// 	t.Errorf("searchMatrix() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

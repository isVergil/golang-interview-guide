package top100liked

import "testing"

// 74. 搜索二维矩阵 (Search a 2D Matrix)
//
// 题目描述:
// 给你一个满足下述两条属性的 m x n 整数矩阵：
// - 每行中的整数从左到右按非递减顺序排列
// - 每行的第一个整数大于前一行的最后一个整数
// 给你一个整数 target，如果 target 在矩阵中，返回 true；否则返回 false。
//
// 示例 1：
// 输入：matrix = [[1,3,5,7],[10,11,16,20],[23,30,34,60]], target = 3
// 输出：true
//
// 示例 2：
// 同上矩阵，target = 13
// 输出：false
//
// 提示：把二维矩阵看成一维有序数组，一次二分搜索

func searchMatrix74(matrix [][]int, target int) bool {
	m, n := len(matrix), len(matrix[0])
	left, right := 0, m*n-1 // 当成长度为 m*n 的一维数组

	for left <= right {
		mid := left + (right-left)/2
		// 一维下标转二维坐标
		val := matrix[mid/n][mid%n]

		if val == target {
			return true
		} else if val < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	return false
}

func TestSearchMatrix74(t *testing.T) {
	matrix := [][]int{{1, 3, 5, 7}, {10, 11, 16, 20}, {23, 30, 34, 60}}

	tests := []struct {
		name     string
		target   int
		expected bool
	}{
		{name: "存在", target: 3, expected: true},
		{name: "不存在", target: 13, expected: false},
		{name: "边界-首", target: 1, expected: true},
		{name: "边界-尾", target: 60, expected: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := searchMatrix74(matrix, tt.target)
			if got != tt.expected {
				t.Errorf("searchMatrix74() = %v, want %v", got, tt.expected)
			}
		})
	}
}

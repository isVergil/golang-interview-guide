package top100liked

import (
	"reflect"
	"testing"
)

// 54. 螺旋矩阵 (Spiral Matrix)
//
// 题目描述:
// 给你一个 m 行 n 列的矩阵 matrix，请按照顺时针螺旋顺序，返回矩阵中的所有元素。
//
// 示例 1：
// 输入：matrix = [[1,2,3],[4,5,6],[7,8,9]]
// 输出：[1,2,3,6,9,8,7,4,5]
//
// 示例 2：
// 输入：matrix = [[1,2,3,4],[5,6,7,8],[9,10,11,12]]
// 输出：[1,2,3,4,8,12,11,10,9,5,6,7]
//
// 提示：维护上下左右四个边界，按圈层遍历

func spiralOrder(matrix [][]int) []int {
	m, n := len(matrix), len(matrix[0])
	res := make([]int, 0, m*n)
	top, bottom, left, right := 0, m-1, 0, n-1
	for top <= bottom && left <= right {
		// 左到右
		for i := left; i <= right; i++ {
			res = append(res, matrix[top][i])
		}
		top++

		// 上到下
		for i := top; i <= bottom; i++ {
			res = append(res, matrix[i][right])
		}
		right--

		// 右到左 防止只有一行时重复
		if top <= bottom {
			for i := right; i >= left; i-- {
				res = append(res, matrix[bottom][i])
			}
			bottom--
		}

		// 下到上 防止只有一列时重复
		if left <= right {
			for i := bottom; i >= top; i-- {
				res = append(res, matrix[i][left])
			}
			left++
		}
	}
	return res
}

func TestSpiralOrder(t *testing.T) {
	tests := []struct {
		name     string
		matrix   [][]int
		expected []int
	}{
		{
			name:     "3x3",
			matrix:   [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
			expected: []int{1, 2, 3, 6, 9, 8, 7, 4, 5},
		},
		{
			name:     "3x4",
			matrix:   [][]int{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10, 11, 12}},
			expected: []int{1, 2, 3, 4, 8, 12, 11, 10, 9, 5, 6, 7},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := spiralOrder(tt.matrix)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("spiralOrder() = %v, want %v", got, tt.expected)
			}
		})
	}
}

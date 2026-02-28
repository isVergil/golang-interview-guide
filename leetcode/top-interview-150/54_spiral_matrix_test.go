package topinterview150

import (
	"testing"
)

// 54. 螺旋矩阵 (Spiral Matrix)
//
// 题目描述:
// 给你一个 m 行 n 列的矩阵 matrix ，请按照 顺时针螺旋顺序 ，返回矩阵中的所有元素。
//
// 示例 1：
// 输入：matrix = [[1,2,3],[4,5,6],[7,8,9]]
// 输出：[1,2,3,6,9,8,7,4,5]
//
// 示例 2：
// 输入：matrix = [[1,2,3,4],[5,6,7,8],[9,10,11,12]]
// 输出：[1,2,3,4,8,12,11,10,9,5,6,7]

func spiralOrder(matrix [][]int) []int {
	if len(matrix) == 0 {
		return []int{}
	}

	// 1 初始化 4 个边界
	top, bottom := 0, len(matrix)-1
	left, right := 0, len(matrix[0])-1

	// 预分配切片空间
	res := make([]int, 0, (bottom+1)*(right+1))

	for {
		// 从左往右 在 top 行 遍历 left - right
		for i := left; i <= right; i++ {
			res = append(res, matrix[top][i])
		}
		top++
		if top > bottom {
			break
		}

		// 从上往下 在 right 列 遍历 top - bottom
		for i := top; i <= bottom; i++ {
			res = append(res, matrix[i][right])
		}
		right--
		if left > right {
			break
		}

		// 从右往左 在 bottom 行 遍历 right - left
		for i := right; i >= left; i-- {
			res = append(res, matrix[bottom][i])
		}
		bottom--
		if top > bottom {
			break
		}

		// 从下往上 在 left 列 遍历 bottom - top
		for i := bottom; i >= top; i-- {
			res = append(res, matrix[i][left])
		}
		left++
		if left > right {
			break
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
			"Example 1",
			[][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
			[]int{1, 2, 3, 6, 9, 8, 7, 4, 5},
		},
		{
			"Example 2",
			[][]int{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10, 11, 12}},
			[]int{1, 2, 3, 4, 8, 12, 11, 10, 9, 5, 6, 7},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// got := spiralOrder(tt.matrix)
			// if !reflect.DeepEqual(got, tt.expected) {
			// 	t.Errorf("spiralOrder() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

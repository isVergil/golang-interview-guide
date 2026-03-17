package topinterview150

import (
	"testing"
)

// 221. 最大正方形 (Maximal Square)
//
// 题目描述:
// 在一个由 '0' 和 '1' 组成的 m x n 二维二进制矩阵中，找出只包含 '1' 的最大正方形，并返回其面积。
//
// 示例 1：
// 输入：matrix = [["1","0","1","0","0"],["1","0","1","1","1"],["1","1","1","1","1"],["1","0","0","1","0"]]
// 输出：4
//
// 示例 2：
// 输入：matrix = [["0","1"],["1","0"]]
// 输出：1
//
// 示例 3：
// 输入：matrix = [["0"]]
// 输出：0

func maximalSquare(matrix [][]byte) int {
	panic("not implemented")
}

func TestMaximalSquare(t *testing.T) {
	tests := []struct {
		name     string
		matrix   [][]byte
		expected int
	}{
		{
			"Example 1",
			[][]byte{
				{'1', '0', '1', '0', '0'},
				{'1', '0', '1', '1', '1'},
				{'1', '1', '1', '1', '1'},
				{'1', '0', '0', '1', '0'},
			},
			4,
		},
		{
			"Example 2",
			[][]byte{{'0', '1'}, {'1', '0'}},
			1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := maximalSquare(tt.matrix); got != tt.expected {
			// 	t.Errorf("maximalSquare() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

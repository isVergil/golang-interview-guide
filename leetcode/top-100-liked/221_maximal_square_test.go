package top100liked

import (
	"reflect"
	"testing"
)

// 221. 最大正方形 (Maximal Square)
//
// 题目描述:
// 在一个由 '0' 和 '1' 组成的二维矩阵内，找到只包含 '1' 的最大正方形，并返回其面积。
//
// 示例 1：
// 输入：matrix = [["1","0","1","0","0"],["1","0","1","1","1"],["1","1","1","1","1"],["1","0","0","1","0"]]
// 输出：4
//
// 示例 2：
// 输入：matrix = [["0","1"],["1","0"]]
// 输出：1
//
// 提示：dp[i][j] = 以 (i,j) 为右下角的最大正方形边长

// 二维 dp
func maximalSquare(matrix [][]byte) int {
	m, n := len(matrix), len(matrix[0])
	dp := make([][]int, m)
	for i := 0; i < m; i++ {
		dp[i] = make([]int, n)
	}
	maxSide := 0
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if matrix[i][j] == '0' {
				continue
			}
			if i == 0 || j == 0 {
				dp[i][j] = 1
			} else {
				dp[i][j] = min(dp[i-1][j-1], min(dp[i-1][j], dp[i][j-1])) + 1
			}
			if dp[i][j] > maxSide {
				maxSide = dp[i][j]
			}
		}
	}

	return maxSide * maxSide
}

// 一维 dp
// 时间 O(mn)，空间 O(n)
func maximalSquare1(matrix [][]byte) int {
	m := len(matrix)
	if m == 0 {
		return 0
	}
	n := len(matrix[0])
	dp := make([]int, n)
	maxSide := 0
	prev := 0 // 相当于 dp[i-1][j-1]

	for i := 0; i < m; i++ {
		prev = 0 // 每行开始时，左上角是 0（相当于越界）
		for j := 0; j < n; j++ {
			temp := dp[j] // 暂存当前 dp[j]，它代表 dp[i-1][j]（即将被更新成 dp[i][j]）
			if matrix[i][j] == '1' {
				if j == 0 {
					dp[j] = 1
				} else {
					// prev 是 dp[i-1][j-1], dp[j] 是 dp[i-1][j], dp[j-1] 是 dp[i][j-1]
					dp[j] = min(prev, min(dp[j], dp[j-1])) + 1
				}
				if dp[j] > maxSide {
					maxSide = dp[j]
				}
			} else {
				dp[j] = 0
			}
			prev = temp // 更新 prev 为刚才的 dp[i-1][j]，供下一个 j 使用
		}
	}
	return maxSide * maxSide
}

func TestMaximalSquare(t *testing.T) {
	tests := []struct {
		name     string
		matrix   [][]byte
		expected int
	}{
		{
			name: "示例1",
			matrix: [][]byte{
				{'1', '0', '1', '0', '0'},
				{'1', '0', '1', '1', '1'},
				{'1', '1', '1', '1', '1'},
				{'1', '0', '0', '1', '0'},
			},
			expected: 4,
		},
		{
			name:     "示例2",
			matrix:   [][]byte{{'0', '1'}, {'1', '0'}},
			expected: 1,
		},
		{
			name:     "全0",
			matrix:   [][]byte{{'0'}},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := maximalSquare(tt.matrix)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("maximalSquare() = %v, want %v", got, tt.expected)
			}
		})

		t.Run(tt.name, func(t *testing.T) {
			got := maximalSquare1(tt.matrix)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("maximalSquare1() = %v, want %v", got, tt.expected)
			}
		})
	}
}

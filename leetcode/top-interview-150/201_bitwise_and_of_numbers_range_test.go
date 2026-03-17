package topinterview150

import (
	"testing"
)

// 201. 数字范围按位与 (Bitwise AND of Numbers Range)
//
// 题目描述:
// 给你两个整数 left 和 right ，表示区间 [left, right] ，返回此区间内所有数字 按位与 的结果（包含 left 、right 端点）。
//
// 示例 1：
// 输入：left = 5, right = 7
// 输出：4
//
// 示例 2：
// 输入：left = 0, right = 0
// 输出：0
//
// 示例 3：
// 输入：left = 1, right = 2147483647
// 输出：0

func rangeBitwiseAnd(left int, right int) int {
	panic("not implemented")
}

func TestRangeBitwiseAnd(t *testing.T) {
	tests := []struct {
		name     string
		left     int
		right    int
		expected int
	}{
		{"Example 1", 5, 7, 4},
		{"Example 2", 0, 0, 0},
		{"Example 3", 1, 2147483647, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := rangeBitwiseAnd(tt.left, tt.right); got != tt.expected {
			// 	t.Errorf("rangeBitwiseAnd() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

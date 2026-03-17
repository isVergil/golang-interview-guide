package topinterview150

import (
	"testing"
)

// 172. 阶乘后的零 (Factorial Trailing Zeroes)
//
// 题目描述:
// 给定一个整数 n ，返回 n! 结果中尾随零的数量。
// 提示 n! = n * (n - 1) * (n - 2) * ... * 3 * 2 * 1
//
// 示例 1：
// 输入：n = 3
// 输出：0
// 解释：3! = 6 ，不含尾随零。
//
// 示例 2：
// 输入：n = 5
// 输出：1
// 解释：5! = 120 ，有一个尾随零。
//
// 示例 3：
// 输入：n = 0
// 输出：0

func trailingZeroes(n int) int {
	panic("not implemented")
}

func TestTrailingZeroes(t *testing.T) {
	tests := []struct {
		name     string
		n        int
		expected int
	}{
		{"Example 1", 3, 0},
		{"Example 2", 5, 1},
		{"Example 3", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := trailingZeroes(tt.n); got != tt.expected {
			// 	t.Errorf("trailingZeroes() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

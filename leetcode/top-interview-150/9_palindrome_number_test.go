package topinterview150

import (
	"testing"
)

// 9. 回文数 (Palindrome Number)
//
// 题目描述:
// 给你一个整数 x ，如果 x 是一个回文整数，返回 true ；否则，返回 false 。
// 回文数是指正序（从左向右）和倒序（从右向左）读都是一样的整数。
// 例如，121 是回文，而 123 不是。
//
// 示例 1：
// 输入：x = 121
// 输出：true
//
// 示例 2：
// 输入：x = -121
// 输出：false
// 解释：从左向右读, 为 -121 。 从右向左读, 为 121- 。因此它不是一个回文数。
//
// 示例 3：
// 输入：x = 10
// 输出：false
// 解释：从右向左读, 为 01 。因此它不是一个回文数。

func isPalindromeNumber(x int) bool {
	panic("not implemented")
}

func TestIsPalindromeNumber(t *testing.T) {
	tests := []struct {
		name     string
		x        int
		expected bool
	}{
		{"Example 1", 121, true},
		{"Example 2", -121, false},
		{"Example 3", 10, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := isPalindromeNumber(tt.x); got != tt.expected {
			// 	t.Errorf("isPalindromeNumber() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

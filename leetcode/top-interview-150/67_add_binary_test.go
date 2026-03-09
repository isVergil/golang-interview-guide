package topinterview150

import (
	"testing"
)

// 67. 二进制求和 (Add Binary)
//
// 题目描述:
// 给你两个二进制字符串 a 和 b ，以二进制字符串的形式返回它们的和。
//
// 示例 1：
// 输入:a = "11", b = "1"
// 输出："100"
//
// 示例 2：
// 输入：a = "1010", b = "1011"
// 输出："10101"

func addBinary(a string, b string) string {
	panic("not implemented")
}

func TestAddBinary(t *testing.T) {
	tests := []struct {
		name     string
		a        string
		b        string
		expected string
	}{
		{"Example 1", "11", "1", "100"},
		{"Example 2", "1010", "1011", "10101"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := addBinary(tt.a, tt.b); got != tt.expected {
			// 	t.Errorf("addBinary() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

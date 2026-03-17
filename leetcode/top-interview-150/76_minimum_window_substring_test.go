package topinterview150

import (
	"testing"
)

// 76. 最小覆盖子串 (Minimum Window Substring)
//
// 题目描述:
// 给你一个字符串 s 、一个字符串 t 。返回 s 中包含 t 所有字符的最小子串。如果 s 中不存在符合条件的子串，则返回空字符串 "" 。
//
// 示例 1：
// 输入：s = "ADOBECODEBANC", t = "ABC"
// 输出："BANC"
//
// 示例 2：
// 输入：s = "a", t = "a"
// 输出："a"
//
// 示例 3:
// 输入: s = "a", t = "aa"
// 输出: ""

func minWindow(s string, t string) string {
	panic("not implemented")
}

func TestMinWindow(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		t        string
		expected string
	}{
		{"Example 1", "ADOBECODEBANC", "ABC", "BANC"},
		{"Example 2", "a", "a", "a"},
		{"Example 3", "a", "aa", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := minWindow(tt.s, tt.t); got != tt.expected {
			// 	t.Errorf("minWindow() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

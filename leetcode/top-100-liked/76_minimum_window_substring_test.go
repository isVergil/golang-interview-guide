package top100liked

import (
	"testing"
)

// 76. 最小覆盖子串 (Minimum Window Substring)
//
// 题目描述:
// 给你一个字符串 s 、一个字符串 t 。返回 s 中涵盖 t 所有字符的最小子串。
// 如果 s 中不存在涵盖 t 所有字符的子串，则返回空字符串 ""。
//
// 示例 1：
// 输入：s = "ADOBECODEBANC", t = "ABC"
// 输出："BANC"
//
// 示例 2：
// 输入：s = "a", t = "a"
// 输出："a"
//
// 示例 3：
// 输入：s = "a", t = "aa"
// 输出：""

func minWindow(s string, t string) string {
	need := make(map[byte]int)
	for i := 0; i < len(t); i++ {
		need[t[i]]++
	}

	remaining := len(t)
	start, minLen := 0, len(s)+1

	l := 0
	for r := 0; r < len(s); r++ {
		if need[s[r]] > 0 {
			remaining--
		}

		need[s[r]]--

		for remaining == 0 {
			if r-l+1 < minLen {
				minLen = r - l + 1
				start = l
			}
			need[s[l]]++
			if need[s[l]] > 0 {
				remaining++
			}
			l++
		}

	}

	if minLen == len(s)+1 {
		return ""
	}

	return s[start : start+minLen]
}

func TestMinWindow(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		tt       string
		expected string
	}{
		{name: "示例1", s: "ADOBECODEBANC", tt: "ABC", expected: "BANC"},
		{name: "示例2", s: "a", tt: "a", expected: "a"},
		{name: "示例3", s: "a", tt: "aa", expected: ""},
		{name: "完全匹配", s: "abc", tt: "abc", expected: "abc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := minWindow(tt.s, tt.tt)
			if got != tt.expected {
				t.Errorf("minWindow(%q, %q) = %q, want %q", tt.s, tt.tt, got, tt.expected)
			}
		})
	}
}

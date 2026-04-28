package top100liked

import (
	"testing"
)

// 5. 最长回文子串 (Longest Palindromic Substring)
//
// 题目描述:
// 给你一个字符串 s，找到 s 中最长的回文子串。
//
// 示例 1：
// 输入：s = "babad"
// 输出："bab"（"aba" 同样是符合题意的答案）
//
// 示例 2：
// 输入：s = "cbbd"
// 输出："bb"

func longestPalindrome(s string) string {
	start, maxLen := 0, 1

	expand := func(l, r int) {
		for l >= 0 && r < len(s) && s[l] == s[r] {
			l--
			r++
		}
		if r-l-1 > maxLen {
			start = l + 1
			maxLen = r - l - 1
		}
	}

	for i := 0; i < len(s); i++ {
		expand(i, i)
		expand(i, i+1)
	}

	return s[start : start+maxLen]
}

func TestLongestPalindrome(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		expected []string
	}{
		{name: "示例1", s: "babad", expected: []string{"bab", "aba"}},
		{name: "示例2", s: "cbbd", expected: []string{"bb"}},
		{name: "单字符", s: "a", expected: []string{"a"}},
		{name: "全相同", s: "aaaa", expected: []string{"aaaa"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := longestPalindrome(tt.s)
			valid := false
			for _, exp := range tt.expected {
				if got == exp {
					valid = true
					break
				}
			}
			if !valid {
				t.Errorf("longestPalindrome(%q) = %q, want one of %v", tt.s, got, tt.expected)
			}
		})
	}
}

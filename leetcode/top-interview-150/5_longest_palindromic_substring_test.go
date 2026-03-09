package topinterview150

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
// 输出："bab"
// 解释："aba" 同样是符合题意的答案。
//
// 示例 2：
// 输入：s = "cbbd"
// 输出："bb"

func longestPalindrome(s string) string {
	if len(s) < 2 {
		return s
	}

	start, end := 0, 0
	for i := 0; i < len(s); i++ {
		len1 := expandCenter(s, i, i)
		len2 := expandCenter(s, i, i+1)

		maxLen := len1
		if len2 > maxLen {
			maxLen = len2
		}

		if maxLen > end-start {
			start = i - (maxLen-1)/2
			end = i + maxLen/2
		}
	}

	return s[start : end+1]
}

func expandCenter(s string, l, r int) int {
	for l >= 0 && r < len(s) && s[l] == s[r] {
		l--
		r++
	}

	return r - l - 1
}

func TestLongestPalindrome(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		expected []string // Possible answers
	}{
		{"Example 1", "babad", []string{"bab", "aba"}},
		{"Example 2", "cbbd", []string{"bb"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// got := longestPalindrome(tt.s)
			// found := false
			// for _, exp := range tt.expected {
			// 	if got == exp {
			// 		found = true
			// 		break
			// 	}
			// }
			// if !found {
			// 	t.Errorf("longestPalindrome() = %v, want one of %v", got, tt.expected)
			// }
		})
	}
}

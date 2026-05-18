package top100liked

import (
	"reflect"
	"testing"
)

// 647. 回文子串 (Palindromic Substrings)
//
// 题目描述:
// 给你一个字符串 s，请你统计并返回这个字符串中回文子串的数目。
// 具有不同开始位置或结束位置的子串，即使是由相同的字符组成，也会被视作不同的子串。
//
// 示例 1：
// 输入：s = "abc"
// 输出：3
// 解释：三个回文子串: "a", "b", "c"
//
// 示例 2：
// 输入：s = "aaa"
// 输出：6
// 解释：6个回文子串: "a", "a", "a", "aa", "aa", "aaa"
//
// 提示：中心扩展法，每个位置分别以奇数和偶数长度向两边扩展

func countSubstrings(s string) int {
	n := len(s)
	count := 0
	for i := 0; i < n; i++ {
		count += expandCountSubstrings(s, i, i)
		count += expandCountSubstrings(s, i, i+1)
	}
	return count
}

func expandCountSubstrings(s string, l, r int) int {
	res := 0
	for l >= 0 && r < len(s) && s[l] == s[r] {
		res++
		l++
		r--
	}
	return res
}

func TestCountSubstrings(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		expected int
	}{
		{name: "无重复", s: "abc", expected: 3},
		{name: "全相同", s: "aaa", expected: 6},
		{name: "单字符", s: "a", expected: 1},
		{name: "回文串", s: "aba", expected: 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := countSubstrings(tt.s)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("countSubstrings() = %v, want %v", got, tt.expected)
			}
		})
	}
}

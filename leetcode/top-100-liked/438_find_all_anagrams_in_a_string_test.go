package top100liked

import (
	"testing"
)

// 438. 找到字符串中所有字母异位词 (Find All Anagrams in a String)
//
// 题目描述:
// 给定两个字符串 s 和 p，找到 s 中所有是 p 的 字母异位词 的子串，返回这些子串的起始索引。不考虑答案输出的顺序。
// 异位词 指由相同字母重排列形成的字符串（包括相同的字符串）。
//
// 示例 1:
// 输入: s = "cbaebabacd", p = "abc"
// 输出: [0,6]
// 解释:
// 起始索引等于 0 的子串是 "cba", 它是 "abc" 的异位词。
// 起始索引等于 6 的子串是 "bac", 它是 "abc" 的异位词。
//
// 示例 2:
// 输入: s = "abab", p = "ab"
// 输出: [0,1,2]
// 解释:
// 起始索引等于 0 的子串是 "ab", 它是 "ab" 的异位词。
// 起始索引等于 1 的子串是 "ba", 它是 "ab" 的异位词。
// 起始索引等于 2 的子串是 "ab", 它是 "ab" 的异位词。

func findAnagrams(s string, p string) []int {
	panic("not implemented")
}

func TestFindAnagrams(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		p        string
		expected []int
	}{
		{"Example 1", "cbaebabacd", "abc", []int{0, 6}},
		{"Example 2", "abab", "ab", []int{0, 1, 2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := findAnagrams(tt.s, tt.p); !reflect.DeepEqual(got, tt.expected) {
			// 	t.Errorf("findAnagrams() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

package top100liked

import (
	"reflect"
	"sort"
	"testing"
)

// 438. 找到字符串中所有字母异位词 (Find All Anagrams in a String)
//
// 题目描述:
// 给定两个字符串 s 和 p，找到 s 中所有 p 的异位词的子串，返回这些子串的起始索引。
// 不考虑答案输出的顺序。异位词指由相同字母重排列形成的字符串（包括相同的字符串）。
//
// 示例 1：
// 输入：s = "cbaebabacd", p = "abc"
// 输出：[0,6]
// 解释：起始索引等于 0 的子串是 "cba"，它是 "abc" 的异位词；
//       起始索引等于 6 的子串是 "bac"，它是 "abc" 的异位词。
//
// 示例 2：
// 输入：s = "abab", p = "ab"
// 输出：[0,1,2]
//
// 提示：滑动窗口 + 字符计数数组

func findAnagrams(s string, p string) []int {
	n, k := len(s), len(p)
	if n < k {
		return nil
	}
	var res []int
	var cnt [26]int
	for _, v := range p {
		cnt[v-'a']++
	}

	l := 0
	for i := 0; i < n; i++ {
		cnt[s[i]-'a']--

		// 窗口超过 k
		if i-l+1 > k {
			cnt[s[l]-'a']++
			l++
		}

		// 窗口长度等于 k 且 cnt 全为 0
		if i-l+1 == k && allZero(cnt) {
			res = append(res, l)
		}
	}
	return res
}

func allZero(cnt [26]int) bool {
	for _, v := range cnt {
		if v != 0 {
			return false
		}
	}
	return true
}

func TestFindAnagrams(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		p        string
		expected []int
	}{
		{name: "示例1", s: "cbaebabacd", p: "abc", expected: []int{0, 6}},
		{name: "示例2", s: "abab", p: "ab", expected: []int{0, 1, 2}},
		{name: "无匹配", s: "hello", p: "xyz", expected: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findAnagrams(tt.s, tt.p)
			sort.Ints(got)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("findAnagrams() = %v, want %v", got, tt.expected)
			}
		})
	}
}

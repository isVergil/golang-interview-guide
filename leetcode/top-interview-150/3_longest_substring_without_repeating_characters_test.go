package topinterview150

import (
	"testing"
)

// 3. 无重复字符的最长子串 (Longest Substring Without Repeating Characters)
//
// 题目描述:
// 给定一个字符串 s ，请你找出其中不含有重复字符的 最长子串 的长度。
//
// 示例 1:
// 输入: s = "abcabcbb"
// 输出: 3
// 解释: 因为无重复字符的最长子串是 "abc"，所以其长度为 3。
//
// 示例 2:
// 输入: s = "bbbbb"
// 输出: 1
// 解释: 因为无重复字符的最长子串是 "b"，所以其长度为 1。
//
// 示例 3:
// 输入: s = "pwwkew"
// 输出: 3
// 解释: 因为无重复字符的最长子串是 "wke"，所以其长度为 3。
//      请注意，你的答案必须是 子串 的长度，"pwke" 是一个子序列，不是子串。

func lengthOfLongestSubstring(s string) int {
	// 上一次出现的下标
	lastChar := [128]int{}
	res, l := 0, 0

	for idx := range lastChar {
		lastChar[idx] = -1
	}

	for r := 0; r < len(s); r++ {
		// 如果字符出现过，且在当前窗口左边界之内，规避 abba 情况
		if lastChar[s[r]] >= l {
			l = lastChar[s[r]] + 1
		}

		lastChar[s[r]] = r

		if r-l+1 > res {
			res = r - l + 1
		}
	}

	return res
}

func TestLengthOfLongestSubstring(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		expected int
	}{
		{"Example 1", "abcabcbb", 3},
		{"Example 2", "bbbbb", 1},
		{"Example 3", "pwwkew", 3},
		{"Empty String", "", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := lengthOfLongestSubstring(tt.s); got != tt.expected {
			// 	t.Errorf("lengthOfLongestSubstring() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

package top100liked

import (
	"testing"
)

// 3. 无重复字符的最长子串 (Longest Substring Without Repeating Characters)
//
// 题目描述:
// 给定一个字符串 s ，请你找出其中不含有重复字符的最长子串的长度。
//
// 示例 1：
// 输入: s = "abcabcbb"
// 输出: 3
// 解释: 因为无重复字符的最长子串是 "abc"，所以其长度为 3。
//
// 示例 2：
// 输入: s = "bbbbb"
// 输出: 1
// 解释: 因为无重复字符的最长子串是 "b"，所以其长度为 1。
//
// 示例 3：
// 输入: s = "pwwkew"
// 输出: 3
// 解释: 因为无重复字符的最长子串是 "wke"，所以其长度为 3。

func lengthOfLongestSubstring(s string) int {
	// 数组代替 map 处理，字符出现的位置
	lastOccurred := [128]int{}

	// 初始化位置为 -1，因为 0 是有效索引
	for i := range lastOccurred {
		lastOccurred[i] = -1
	}

	res, start := 0, 0

	for end, char := range s {
		// 出现的位置超过了起始 start
		if lastPos := lastOccurred[char]; lastPos >= start {
			start = lastPos + 1
		}

		// 更新当前字符的位置
		lastOccurred[char] = end

		// 更新最大长度
		if end-start+1 > res {
			res = end - start + 1
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
		{
			name:     "示例1",
			s:        "abcabcbb",
			expected: 3,
		},
		{
			name:     "示例2",
			s:        "bbbbb",
			expected: 1,
		},
		{
			name:     "示例3",
			s:        "pwwkew",
			expected: 3,
		},
		{
			name:     "空字符串",
			s:        "",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := lengthOfLongestSubstring(tt.s)
			if got != tt.expected {
				t.Errorf("lengthOfLongestSubstring() = %v, want %v", got, tt.expected)
			}
		})
	}
}

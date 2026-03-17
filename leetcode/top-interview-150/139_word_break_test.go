package topinterview150

import (
	"testing"
)

// 139. 单词拆分 (Word Break)
//
// 题目描述:
// 给你一个字符串 s 和一个字符串列表 wordDict 作为字典。请你判断是否可以利用字典中出现的单词拼接出 s 。
// 注意：不要求字典中出现的单词全部都使用，并且字典中的单词可以重复使用。
//
// 示例 1：
// 输入: s = "leetcode", wordDict = ["leet", "code"]
// 输出: true
// 解释: 返回 true 因为 "leetcode" 可以由 "leet" 和 "code" 拼接成。
//
// 示例 2：
// 输入: s = "applepenapple", wordDict = ["apple", "pen"]
// 输出: true
// 解释: 返回 true 因为 "applepenapple" 可以由 "apple" "pen" "apple" 拼接成。
//      注意，你可以重复使用字典中的单词。
//
// 示例 3：
// 输入: s = "catsandog", wordDict = ["cats", "dog", "sand", "and", "cat"]
// 输出: false

func wordBreak(s string, wordDict []string) bool {
	panic("not implemented")
}

func TestWordBreak(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		wordDict []string
		expected bool
	}{
		{"Example 1", "leetcode", []string{"leet", "code"}, true},
		{"Example 2", "applepenapple", []string{"apple", "pen"}, true},
		{"Example 3", "catsandog", []string{"cats", "dog", "sand", "and", "cat"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := wordBreak(tt.s, tt.wordDict); got != tt.expected {
			// 	t.Errorf("wordBreak() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

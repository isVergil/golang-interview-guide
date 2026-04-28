package top100liked

import (
	"testing"
)

// 139. 单词拆分 (Word Break)
//
// 题目描述:
// 给你一个字符串 s 和一个字符串列表 wordDict 作为字典。
// 如果可以利用字典中出现的一个或多个单词拼接出 s 则返回 true。
// 注意：不要求字典中出现的单词全部都使用，并且字典中的单词可以重复使用。
//
// 示例 1：
// 输入：s = "leetcode", wordDict = ["leet","code"]
// 输出：true（"leetcode" 可以由 "leet" + "code" 拼成）
//
// 示例 2：
// 输入：s = "applepenapple", wordDict = ["apple","pen"]
// 输出：true（"applepenapple" 可以由 "apple" + "pen" + "apple" 拼成，单词可复用）
//
// 示例 3：
// 输入：s = "catsandog", wordDict = ["cats","dog","sand","and","cat"]
// 输出：false

func wordBreak(s string, wordDict []string) bool {
	dict := make(map[string]bool)
	maxLen := 0
	for _, w := range wordDict {
		dict[w] = true
		if len(w) > maxLen {
			maxLen = len(w)
		}
	}

	n := len(s)

	// dp[i] 表示 s 的前 i 个字符能否被拆分
	dp := make([]bool, n+1)
	dp[0] = true

	for i := 1; i <= n; i++ {
		// 从位置 i 往前截，截 1 个字符、2 个字符...最多截 maxLen 个
		for wordLen := 1; wordLen <= maxLen && wordLen <= i; wordLen++ {
			// 截出来的词是 s[i-wordLen : i]
			if dp[i-wordLen] && dict[s[i-wordLen:i]] {
				dp[i] = true
				break
			}
		}
	}

	return dp[n]

}

func TestWordBreak(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		wordDict []string
		expected bool
	}{
		{name: "示例1", s: "leetcode", wordDict: []string{"leet", "code"}, expected: true},
		{name: "示例2", s: "applepenapple", wordDict: []string{"apple", "pen"}, expected: true},
		{name: "示例3", s: "catsandog", wordDict: []string{"cats", "dog", "sand", "and", "cat"}, expected: false},
		{name: "空串", s: "", wordDict: []string{"a"}, expected: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := wordBreak(tt.s, tt.wordDict)
			if got != tt.expected {
				t.Errorf("wordBreak(%q) = %v, want %v", tt.s, got, tt.expected)
			}
		})
	}
}

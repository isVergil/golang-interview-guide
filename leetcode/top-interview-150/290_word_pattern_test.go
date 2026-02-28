package topinterview150

import (
	"strings"
	"testing"
)

// 290. 单词规律 (Word Pattern)
//
// 题目描述:
// 给定一种规律 pattern 和一个字符串 s ，判断 s 是否遵循相同的规律。
// 这里的 遵循 指完全匹配，例如， pattern 里的每个字母和字符串 s 中的每个非空单词之间存在着双向连接的映射规律。
//
// 示例 1:
// 输入: pattern = "abba", s = "dog cat cat dog"
// 输出: true
//
// 示例 2:
// 输入: pattern = "abba", s = "dog cat cat fish"
// 输出: false
//
// 示例 3:
// 输入: pattern = "aaaa", s = "dog cat cat dog"
// 输出: false

func wordPattern(pattern string, s string) bool {
	// 拆分单词，strings.Fields 性能优于 strings.Split(s, " ")
	words := strings.Fields(s)

	if len(pattern) != len(words) {
		return false
	}

	// 建立双向映射
	p2w := [26]string{}
	w2p := make(map[string]byte)

	for i := 0; i < len(pattern); i++ {
		p := pattern[i] - 'a'
		word := words[i]

		// 已经映射
		if p2w[p] != "" {
			if p2w[p] != word {
				return false
			}
		} else {
			// 还没映射，但单词已经被别的字符占了，说明不同构
			if _, ok := w2p[word]; ok {
				return false
			}

			// 建立双向绑定
			p2w[p] = word
			w2p[word] = pattern[i]
		}
	}
	return true

}

func TestWordPattern(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		s        string
		expected bool
	}{
		{"Example 1", "abba", "dog cat cat dog", true},
		{"Example 2", "abba", "dog cat cat fish", false},
		{"Example 3", "aaaa", "dog cat cat dog", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := wordPattern(tt.pattern, tt.s); got != tt.expected {
			// 	t.Errorf("wordPattern() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

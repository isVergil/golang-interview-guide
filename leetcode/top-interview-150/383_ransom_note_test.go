package topinterview150

import (
	"testing"
)

// 383. 赎金信 (Ransom Note)
//
// 题目描述:
// 给你两个字符串：ransomNote 和 magazine ，判断 ransomNote 能不能由 magazine 里面的字符构成。
// 如果可以，返回 true ；否则返回 false 。
// magazine 中的每个字符只能在 ransomNote 中使用一次。
//
// 示例 1：
// 输入：ransomNote = "a", magazine = "b"
// 输出：false
//
// 示例 2：
// 输入：ransomNote = "aa", magazine = "ab"
// 输出：false
//
// 示例 3：
// 输入：ransomNote = "aa", magazine = "aab"
// 输出：true

func canConstruct(ransomNote string, magazine string) bool {
	panic("not implemented")
}

func TestCanConstruct(t *testing.T) {
	tests := []struct {
		name       string
		ransomNote string
		magazine   string
		expected   bool
	}{
		{"Example 1", "a", "b", false},
		{"Example 2", "aa", "ab", false},
		{"Example 3", "aa", "aab", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := canConstruct(tt.ransomNote, tt.magazine); got != tt.expected {
			// 	t.Errorf("canConstruct() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

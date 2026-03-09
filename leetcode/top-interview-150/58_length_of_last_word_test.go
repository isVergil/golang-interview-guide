package topinterview150

import (
	"testing"
)

// 58. 最后一个单词的长度 (Length of Last Word)
//
// 题目描述:
// 给你一个字符串 s，由若干单词组成，单词前后用一些空格隔开。返回字符串中 最后一个 单词的长度。
// 单词 是指仅由字母组成、不包含任何空格字符的最大子串。
//
// 示例 1：
// 输入：s = "Hello World"
// 输出：5
// 解释：最后一个单词是“World”，长度为 5。
//
// 示例 2：
// 输入：s = "   fly me   to   the moon  "
// 输出：4
// 解释：最后一个单词是“moon”，长度为 4。
//
// 示例 3：
// 输入：s = "luffy is still joyboy"
// 输出：6
// 解释：最后一个单词是“joyboy”，长度为 6。

func lengthOfLastWord(s string) int {
	panic("not implemented")
}

func TestLengthOfLastWord(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		expected int
	}{
		{"Example 1", "Hello World", 5},
		{"Example 2", "   fly me   to   the moon  ", 4},
		{"Example 3", "luffy is still joyboy", 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := lengthOfLastWord(tt.s); got != tt.expected {
			// 	t.Errorf("lengthOfLastWord() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

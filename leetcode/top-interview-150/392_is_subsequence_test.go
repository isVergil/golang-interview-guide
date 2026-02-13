package topinterview150

import "testing"

// 392. 判断子序列 (Is Subsequence)
//
// 题目描述:
// 给定字符串 s 和 t ，判断 s 是否为 t 的子序列。
// 字符串的一个子序列是原始字符串删除一些（也可以不删除）字符而不改变剩余字符相对位置形成的新字符串。（例如，"ace"是"abcde"的一个子序列，而"aec"不是）。
// 进阶：
// 如果有大量输入的 S，称作 S1, S2, ... , Sk 其中 k >= 10亿，你需要依次检查它们是否为 T 的子序列。在这种情况下，你会怎样改变代码？
//
// 示例 1：
// 输入：s = "abc", t = "ahbgdc"
// 输出：true
//
// 示例 2：
// 输入：s = "axc", t = "ahbgdc"
// 输出：false

func isSubsequence(s string, t string) bool {
	i, j := 0, 0
	for i < len(s) && j < len(t) {
		if s[i] == t[j] {
			i++
		}
		j++
	}
	return i == len(s)
}

func TestIsSubsequence(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		t        string
		expected bool
	}{
		{
			name:     "Example 1",
			s:        "abc",
			t:        "ahbgdc",
			expected: true,
		},
		{
			name:     "Example 2",
			s:        "axc",
			t:        "ahbgdc",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: Uncomment when implementation is ready
			// if got := isSubsequence(tt.s, tt.t); got != tt.expected {
			// 	t.Errorf("isSubsequence() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

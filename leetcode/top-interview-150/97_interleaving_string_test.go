package topinterview150

import (
	"testing"
)

// 97. 交错字符串 (Interleaving String)
//
// 题目描述:
// 给定三个字符串 s1, s2, s3，请你帮忙验证 s3 是否是由 s1 和 s2 交错 组成的。
// 两个字符串 s 和 t 的 交错 含义为：
// s = s1 + s2 + ... + sn
// t = t1 + t2 + ... + tm
// |n - m| <= 1
// 交错 是 s1 + t1 + s2 + t2 + s3 + t3 + ... 或者 t1 + s1 + t2 + s2 + t3 + s3 + ...
// 注意：a + b 意味着字符串 a 和 b 连接。
//
// 示例 1：
// 输入：s1 = "aabcc", s2 = "dbbca", s3 = "aadbbcbcac"
// 输出：true
//
// 示例 2：
// 输入：s1 = "aabcc", s2 = "dbbca", s3 = "aadbbbaccc"
// 输出：false
//
// 示例 3：
// 输入：s1 = "", s2 = "", s3 = ""
// 输出：true

func isInterleave(s1 string, s2 string, s3 string) bool {
	panic("not implemented")
}

func TestIsInterleave(t *testing.T) {
	tests := []struct {
		name     string
		s1       string
		s2       string
		s3       string
		expected bool
	}{
		{"Example 1", "aabcc", "dbbca", "aadbbcbcac", true},
		{"Example 2", "aabcc", "dbbca", "aadbbbaccc", false},
		{"Example 3", "", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := isInterleave(tt.s1, tt.s2, tt.s3); got != tt.expected {
			// 	t.Errorf("isInterleave() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

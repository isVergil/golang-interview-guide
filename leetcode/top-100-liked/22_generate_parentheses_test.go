package top100liked

import (
	"testing"
)

// 22. 括号生成 (Generate Parentheses)
//
// 题目描述:
// 数字 n 代表生成括号的对数，请你设计一个函数，用于能够生成所有可能的并且有效的括号组合。
//
// 示例 1：
// 输入：n = 3
// 输出：["((()))","(()())","(())()","()(())","()()()"]
//
// 示例 2：
// 输入：n = 1
// 输出：["()"]

func generateParenthesis(n int) []string {
	res := make([]string, 0)
	path := make([]byte, 0)
	var backtrack = func(l, r int) {}
	backtrack = func(l, r int) {
		if len(path) == 2*n {
			res = append(res, string(path))
			return
		}

		// 选左
		if l < n {
			path = append(path, '(')
			backtrack(l+1, r)
			path = path[:len(path)-1]
		}

		// 选右
		if l > r {
			path = append(path, ')')
			backtrack(l, r+1)
			path = path[:len(path)-1]
		}
	}

	backtrack(0, 0)

	return res
}

func TestGenerateParenthesis(t *testing.T) {
	tests := []struct {
		name        string
		n           int
		expectedLen int
	}{
		{
			name:        "示例1",
			n:           3,
			expectedLen: 5,
		},
		{
			name:        "示例2",
			n:           1,
			expectedLen: 1,
		},
		{
			name:        "n=2",
			n:           2,
			expectedLen: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generateParenthesis(tt.n)
			if len(got) != tt.expectedLen {
				t.Errorf("generateParenthesis() returned %d results, want %d", len(got), tt.expectedLen)
			}
		})
	}
}

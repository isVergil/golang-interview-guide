package top100liked

import (
	"testing"
)

// 20. 有效的括号 (Valid Parentheses)
//
// 题目描述:
// 给定一个只包括 '('，')'，'{'，'}'，'['，']' 的字符串 s ，判断字符串是否有效。
// 有效字符串需满足：
// 1. 左括号必须用相同类型的右括号闭合。
// 2. 左括号必须以正确的顺序闭合。
// 3. 每个右括号都有一个对应的相同类型的左括号。
//
// 示例 1：
// 输入：s = "()"
// 输出：true
//
// 示例 2：
// 输入：s = "()[]{}"
// 输出：true
//
// 示例 3：
// 输入：s = "(]"
// 输出：false

var pairs = map[rune]rune{
	')': '(',
	']': '[',
	'}': '{',
}

func isValid(s string) bool {
	stack := make([]rune, 0)
	for _, v := range s {
		if cur, ok := pairs[v]; ok {
			if len(stack) > 0 && stack[len(stack)-1] == cur {
				stack = stack[:len(stack)-1]
			} else {
				return false
			}
		} else {
			stack = append(stack, v)
		}
	}
	return len(stack) == 0
}

func TestIsValid(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		expected bool
	}{
		{
			name:     "示例1",
			s:        "()",
			expected: true,
		},
		{
			name:     "示例2",
			s:        "()[]{}",
			expected: true,
		},
		{
			name:     "示例3",
			s:        "(]",
			expected: false,
		},
		{
			name:     "嵌套括号",
			s:        "([{}])",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValid(tt.s)
			if got != tt.expected {
				t.Errorf("isValid() = %v, want %v", got, tt.expected)
			}
		})
	}
}

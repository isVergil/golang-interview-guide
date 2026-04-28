package top100liked

import (
	"testing"
)

// 32. 最长有效括号 (Longest Valid Parentheses)
//
// 题目描述:
// 给你一个只包含 '(' 和 ')' 的字符串，找出最长有效（格式正确且连续）括号子串的长度。
//
// 示例 1：
// 输入：s = "(()"
// 输出：2
//
// 示例 2：
// 输入：s = ")()())"
// 输出：4
//
// 示例 3：
// 输入：s = ""
// 输出：0

func longestValidParentheses(s string) int {
	res := 0
	// 栈里存下标，栈底是「最近一个没被匹配的右括号的下标」作为分隔符
	// 初始放 -1，表示起始分隔点，也就是虚拟点
	stack := []int{-1}
	for i := 0; i < len(s); i++ {
		if s[i] == '(' {
			// 左括号：下标入栈
			stack = append(stack, i)
		} else {
			// 右括号：弹出栈顶
			stack = stack[:len(stack)-1]
			if len(stack) == 0 {
				// 栈空了，当前 ')' 无法匹配，作为新的分隔符入栈
				stack = append(stack, i)
			} else {
				// 当前有效长度 = 当前下标 - 栈顶（分隔符或上一个未匹配的左括号）
				length := i - stack[len(stack)-1]
				if length > res {
					res = length
				}
			}
		}
	}

	return res
}

func TestLongestValidParentheses(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		expected int
	}{
		{name: "示例1", s: "(()", expected: 2},
		{name: "示例2", s: ")()())", expected: 4},
		{name: "示例3", s: "", expected: 0},
		{name: "全匹配", s: "()()", expected: 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := longestValidParentheses(tt.s)
			if got != tt.expected {
				t.Errorf("longestValidParentheses(%q) = %v, want %v", tt.s, got, tt.expected)
			}
		})
	}
}

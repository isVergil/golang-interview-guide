package top100liked

import (
	"strings"
	"testing"
)

// 394. 字符串解码 (Decode String)
//
// 题目描述:
// 给定一个经过编码的字符串，返回它解码后的字符串。
// 编码规则为: k[encoded_string]，表示方括号内部的字符串正好重复 k 次。
//
// 示例 1：
// 输入：s = "3[a]2[bc]"
// 输出："aaabcbc"
//
// 示例 2：
// 输入：s = "3[a2[c]]"
// 输出："accaccacc"
//
// 示例 3：
// 输入：s = "2[abc]3[cd]ef"
// 输出："abcabccdcdcdef"
//
// 提示：两个栈，一个存数字一个存字符串，遇到 [ 压栈，遇到 ] 弹栈拼接

func decodeString(s string) string {
	numStack := []int{}
	strStack := []string{}
	cur := ""
	num := 0
	for _, ch := range s {
		switch {
		case ch >= '0' && ch <= '9':
			num = num*10 + int(ch-'0')
		case ch == '[':
			numStack = append(numStack, num)
			strStack = append(strStack, cur)
			num = 0
			cur = ""
		case ch == ']':
			repeat := numStack[len(numStack)-1]
			numStack = numStack[:len(numStack)-1]
			prefix := strStack[len(strStack)-1]
			strStack = strStack[:len(strStack)-1]
			cur = prefix + strings.Repeat(cur, repeat)
		default:
			cur += string(ch)
		}
	}
	return cur
}

func TestDecodeString(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		expected string
	}{
		{name: "示例1", s: "3[a]2[bc]", expected: "aaabcbc"},
		{name: "嵌套", s: "3[a2[c]]", expected: "accaccacc"},
		{name: "多段拼接", s: "2[abc]3[cd]ef", expected: "abcabccdcdcdef"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := decodeString(tt.s)
			if got != tt.expected {
				t.Errorf("decodeString() = %v, want %v", got, tt.expected)
			}
		})
	}
}

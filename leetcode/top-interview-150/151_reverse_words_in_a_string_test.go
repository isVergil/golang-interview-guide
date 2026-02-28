package topinterview150

import (
	"strings"
	"testing"
)

// 151. 反转字符串中的单词 (Reverse Words in a String)
//
// 题目描述:
// 给你一个字符串 s ，请你反转字符串中 单词 的顺序。
// 单词 是由非空格字符组成的字符串。s 中至少存在一个单词。
// 返回 单词 顺序反转且 单词 之间用单个空格连接的结果字符串。
// 注意：输入字符串 s中可能会包含前导空格、尾随空格或者单词间的多个空格。返回的结果字符串中，单词间应当仅用单个空格分隔，且不包含任何额外的空格。
//
// 示例 1：
// 输入：s = "the sky is blue"
// 输出："blue is sky the"
//
// 示例 2：
// 输入：s = "  hello world  "
// 输出："world hello"
// 解释：反转后的字符串中不能包含前导空格和尾随空格。
//
// 示例 3：
// 输入：s = "a good   example"
// 输出："example good a"
// 解释：如果两个单词间有多余的空格，反转后的字符串需要将单词间的空格减少到仅有一个。

func reverseWords(s string) string {
	sArr := strings.Split(s, " ")
	res := make([]string, 0, len(sArr))
	for i := len(sArr) - 1; i >= 0; i-- {
		if sArr[i] != "" {
			res = append(res, sArr[i])
		}
	}
	return strings.Join(res, " ")
}

func TestReverseWords(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		expected string
	}{
		{"Example 1", "the sky is blue", "blue is sky the"},
		{"Example 2", "  hello world  ", "world hello"},
		{"Example 3", "a good   example", "example good a"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := reverseWords(tt.s); got != tt.expected {
			// 	t.Errorf("reverseWords() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

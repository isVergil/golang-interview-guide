package topinterview150

import (
	"testing"
)

// 6. Z 字形变换 (ZigZag Conversion)
//
// 题目描述:
// 将一个给定字符串 s 根据给定的行数 numRows ，以从上往下、从左到右进行 Z 字形排列。
// 比如输入字符串为 "PAYPALISHIRING" 行数为 3 时，排列如下：
// P   A   H   N
// A P L S I I G
// Y   I   R
// 之后，你的输出需要从左往右逐行读取，产生出一个新的字符串，比如："PAHNAPLSIIGYIR"。
//
// 示例 1：
// 输入：s = "PAYPALISHIRING", numRows = 3
// 输出："PAHNAPLSIIGYIR"
//
// 示例 2：
// 输入：s = "PAYPALISHIRING", numRows = 4
// 输出："PINALSIGYAHRPI"
// 解释：
// P     I    N
// A   L S  I G
// Y A   H R
// P     I
//
// 示例 3：
// 输入：s = "A", numRows = 1
// 输出："A"

func convert(s string, numRows int) string {
	panic("not implemented")
}

func TestConvert(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		numRows  int
		expected string
	}{
		{"Example 1", "PAYPALISHIRING", 3, "PAHNAPLSIIGYIR"},
		{"Example 2", "PAYPALISHIRING", 4, "PINALSIGYAHRPI"},
		{"Example 3", "A", 1, "A"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := convert(tt.s, tt.numRows); got != tt.expected {
			// 	t.Errorf("convert() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

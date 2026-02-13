package topinterview150

import (
	"testing"
)

// 125. 验证回文串 (Valid Palindrome)
//
// 题目描述:
// 如果在将所有大写字符转换为小写字符、并移除所有非字母数字字符之后，短语正着读和反着读都一样，则可以认为该短语是一个 回文串 。
// 字母和数字都属于字母数字字符。
// 给你一个字符串 s，如果它是 回文串 ，返回 true ；否则，返回 false 。
//
// 示例 1：
// 输入: s = "A man, a plan, a canal: Panama"
// 输出: true
// 解释："amanaplanacanalpanama" 是回文串。
//
// 示例 2：
// 输入: s = "race a car"
// 输出: false
// 解释："raceacar" 不是回文串。
//
// 示例 3：
// 输入: s = " "
// 输出: true
// 解释：在移除非字母数字字符之后，s 是一个空字符串 "" 。
// 由于空字符串正着反着读都一样，所以是回文串。

func isPalindrome(s string) bool {
	i, j := 0, len(s)-1
	for i < j {
		// if !unicode.IsLetter(rune(s[i])) && !unicode.IsDigit(rune(s[i])) {
		if !isLetter(s[i]) && !isDigit(s[i]) {
			i++
			// } else if !unicode.IsLetter(rune(s[j])) && !unicode.IsDigit(rune(s[j])) {
		} else if !isLetter(s[j]) && !isDigit(s[j]) {
			j--
			// } else if unicode.ToLower(rune(s[i])) == unicode.ToLower(rune(s[j])) {
		} else if toLower(s[i]) == toLower(s[j]) {
			i++
			j--
		} else {
			return false
		}
	}
	return true
}

// 判断是否为字母 (大写或小写)
func isLetter(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}

// 判断是否为数字
func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

// 'A' 是 65，'a' 是 97。
func toLower(b byte) byte {
	if b >= 'A' && b <= 'Z' {
		//return b | 32 // 或者 b | ' ' (空格的 ASCII 也是 32)
		return b + 32
	}
	return b
}

func TestIsPalindrome(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		expected bool
	}{
		{
			name:     "Example 1",
			s:        "A man, a plan, a canal: Panama",
			expected: true,
		},
		{
			name:     "Example 2",
			s:        "race a car",
			expected: false,
		},
		{
			name:     "Example 3",
			s:        " ",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: Uncomment when implementation is ready
			// if got := isPalindrome(tt.s); got != tt.expected {
			// 	t.Errorf("isPalindrome() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

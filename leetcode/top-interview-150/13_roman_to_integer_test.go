package topinterview150

import (
	"testing"
)

// 13. 罗马数字转整数 (Roman to Integer)
//
// 题目描述:
// 给定一个罗马数字，将其转换成整数。
//
// 示例 1:
// 输入: s = "III"
// 输出: 3

func romanToInt(s string) int {
	// 细节优化 1：使用 switch 替代 map，减少哈希开销
	// 细节优化 2：直接遍历字节切片，不涉及字符串转换
	symbolValue := func(b byte) int {
		switch b {
		case 'I':
			return 1
		case 'V':
			return 5
		case 'X':
			return 10
		case 'L':
			return 50
		case 'C':
			return 100
		case 'D':
			return 500
		case 'M':
			return 1000
		default:
			return 0
		}
	}

	ans := 0
	n := len(s)
	// 遍历到倒数第二个字符
	for i := 0; i < n; i++ {
		value := symbolValue(s[i])
		// 细节优化 3：提前获取下一个字符的值进行比较
		if i < n-1 && value < symbolValue(s[i+1]) {
			ans -= value
		} else {
			ans += value
		}
	}
	return ans
}

func TestRomanToInt(t *testing.T) {
	// 罗马数字转整数测试
}

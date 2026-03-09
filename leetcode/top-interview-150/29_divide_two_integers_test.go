package topinterview150

import (
	"testing"
)

// 29. 两数相除 (Divide Two Integers)
//
// 题目描述:
// 给你两个整数，被除数 dividend 和除数 divisor。将两数相除，要求 不使用 乘法、除法和取余运算。
// 整数除法应该向零截断，也就是截去其小数部分。例如，truncate(8.345) = 8 以及 truncate(-2.7335) = -2。
// 返回被除数 dividend 除以除数 divisor 得到的 商 。
// 注意：假设我们的环境只能存储 32 位 有符号整数，其数值范围是 [−2^31,  2^31 − 1]。本题中，如果商 严格大于 2^31 − 1 ，则返回 2^31 − 1 ；如果商 严格小于 -2^31 ，则返回 -2^31 。
//
// 示例 1:
// 输入: dividend = 10, divisor = 3
// 输出: 3
// 解释: 10/3 = 3.33333... ，向零截断后得到 3 。
//
// 示例 2:
// 输入: dividend = 7, divisor = -3
// 输出: -2
// 解释: 7/-3 = -2.33333... ，向零截断后得到 -2 。

func divide(dividend int, divisor int) int {
	panic("not implemented")
}

func TestDivide(t *testing.T) {
	tests := []struct {
		name     string
		dividend int
		divisor  int
		expected int
	}{
		{"Example 1", 10, 3, 3},
		{"Example 2", 7, -3, -2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := divide(tt.dividend, tt.divisor); got != tt.expected {
			// 	t.Errorf("divide() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

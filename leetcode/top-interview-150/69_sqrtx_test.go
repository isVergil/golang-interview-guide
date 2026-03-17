package topinterview150

import (
	"testing"
)

// 69. x 的平方根 (Sqrt(x))
//
// 题目描述:
// 给你一个非负整数 x ，计算并返回 x 的 算术平方根 。
// 由于返回类型是整数，结果只保留 整数部分 ，小数部分将被 舍去 。
// 注意：不允许使用任何内置指数函数和算符，例如 pow(x, 0.5) 或者 x ** 0.5 。
//
// 示例 1：
// 输入：x = 4
// 输出：2
//
// 示例 2：
// 输入：x = 8
// 输出：2
// 解释：8 的算术平方根是 2.82842..., 由于返回类型是整数，小数部分将被舍去。

func mySqrt(x int) int {
	panic("not implemented")
}

func TestMySqrt(t *testing.T) {
	tests := []struct {
		name     string
		x        int
		expected int
	}{
		{"Example 1", 4, 2},
		{"Example 2", 8, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := mySqrt(tt.x); got != tt.expected {
			// 	t.Errorf("mySqrt() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

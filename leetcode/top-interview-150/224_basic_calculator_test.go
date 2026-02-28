package topinterview150

import (
	"testing"
)

// 224. 基本计算器 (Basic Calculator)
//
// 题目描述:
// 给你一个字符串表达式 s ，请你实现一个基本计算器来计算并返回它的值。
// 注意:不允许使用任何将字符串作为数学表达式计算的内置函数，比如 eval() 。
//
// 示例 1：
// 输入：s = "1 + 1"
// 输出：2
//
// 示例 2：
// 输入：s = " 2-1 + 2 "
// 输出：3
//
// 示例 3：
// 输入：s = "(1+(4+5+2)-3)+(6+8)"
// 输出：23

func calculate(s string) int {
	// ops 栈存储当前括号级别的全局符号状态
	ops := []int{1}
	sign := 1
	res := 0
	n := len(s)

	for i := 0; i < n; i++ {
		switch s[i] {
		case ' ':
			continue
		case '+':
			// 当前符号 = 1 * 括号外的环境符号
			sign = ops[len(ops)-1]
		case '-':
			// 当前符号 = -1 * 括号外的环境符号
			sign = -ops[len(ops)-1]
		case '(':
			// 进入括号，把当前的 sign 变成环境符号压栈
			ops = append(ops, sign)
		case ')':
			// 退出括号，弹出环境符号
			ops = ops[:len(ops)-1]
		default:
			// 处理数字（考虑到数字可能是多位数）
			num := 0
			for i < n && s[i] >= '0' && s[i] <= '9' {
				num = num*10 + int(s[i]-'0')
				i++
			}
			// 累加结果
			res += sign * num
			i-- // for 循环还会 i++，所以这里要减回来
		}
	}
	return res
}

func TestCalculate(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		expected int
	}{
		{"Example 1", "1 + 1", 2},
		{"Example 2", " 2-1 + 2 ", 3},
		{"Example 3", "(1+(4+5+2)-3)+(6+8)", 23},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := calculate(tt.s); got != tt.expected {
			// 	t.Errorf("calculate() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

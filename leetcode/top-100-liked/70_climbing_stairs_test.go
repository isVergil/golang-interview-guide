package top100liked

import (
	"testing"
)

// 70. 爬楼梯 (Climbing Stairs)
//
// 题目描述:
// 假设你正在爬楼梯。需要 n 阶你才能到达楼顶。
// 每次你可以爬 1 或 2 个台阶。你有多少种不同的方法可以爬到楼顶呢？
//
// 示例 1：
// 输入：n = 2
// 输出：2
// 解释：有两种方法可以爬到楼顶。
// 1. 1 阶 + 1 阶
// 2. 2 阶
//
// 示例 2：
// 输入：n = 3
// 输出：3
// 解释：有三种方法可以爬到楼顶。
// 1. 1 阶 + 1 阶 + 1 阶
// 2. 1 阶 + 2 阶
// 3. 2 阶 + 1 阶

func climbStairs(n int) int {
	if n <= 2 {
		return n
	}
	res, step1, step2 := 0, 1, 2
	for i := 3; i <= n; i++ {
		res = step1 + step2
		step1 = step2
		step2 = res
	}
	return res
}

func TestClimbStairs(t *testing.T) {
	tests := []struct {
		name     string
		n        int
		expected int
	}{
		{
			name:     "示例1",
			n:        2,
			expected: 2,
		},
		{
			name:     "示例2",
			n:        3,
			expected: 3,
		},
		{
			name:     "n=1",
			n:        1,
			expected: 1,
		},
		{
			name:     "n=10",
			n:        10,
			expected: 89,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := climbStairs(tt.n)
			if got != tt.expected {
				t.Errorf("climbStairs() = %v, want %v", got, tt.expected)
			}
		})
	}
}

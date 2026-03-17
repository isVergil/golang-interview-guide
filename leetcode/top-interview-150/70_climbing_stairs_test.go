package topinterview150

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
	panic("not implemented")
}

func TestClimbStairs(t *testing.T) {
	tests := []struct {
		name     string
		n        int
		expected int
	}{
		{"Example 1", 2, 2},
		{"Example 2", 3, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := climbStairs(tt.n); got != tt.expected {
			// 	t.Errorf("climbStairs() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

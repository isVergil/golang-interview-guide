package topinterview150

import (
	"testing"
)

// 322. 零钱兑换 (Coin Change)
//
// 题目描述:
// 给你一个整数数组 coins ，表示不同面额的硬币；以及一个整数 amount ，表示总金额。
// 计算并返回可以凑成总金额所需的 最少的硬币个数 。如果没有任何一种硬币组合能组成总金额，返回 -1 。
// 你可以认为每种硬币的数量是无限的。
//
// 示例 1：
// 输入：coins = [1, 2, 5], amount = 11
// 输出：3
// 解释：11 = 5 + 5 + 1
//
// 示例 2：
// 输入：coins = [2], amount = 3
// 输出：-1
//
// 示例 3：
// 输入：coins = [1], amount = 0
// 输出：0

func coinChange(coins []int, amount int) int {
	panic("not implemented")
}

func TestCoinChange(t *testing.T) {
	tests := []struct {
		name     string
		coins    []int
		amount   int
		expected int
	}{
		{"Example 1", []int{1, 2, 5}, 11, 3},
		{"Example 2", []int{2}, 3, -1},
		{"Example 3", []int{1}, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := coinChange(tt.coins, tt.amount); got != tt.expected {
			// 	t.Errorf("coinChange() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

package top100liked

import (
	"testing"
)

// 322. 零钱兑换 (Coin Change)
//
// 题目描述:
// 给你一个整数数组 coins ，表示不同面额的硬币；以及一个整数 amount ，表示总金额。
// 计算并返回可以凑成总金额所需的最少的硬币个数。如果没有任何一种硬币组合能组成总金额，返回 -1。
// 你可以认为每种硬币的数量是无限的。
//
// 示例 1：
// 输入：coins = [1, 2, 5], amount = 11
// 输出：3（11 = 5 + 5 + 1）
//
// 示例 2：
// 输入：coins = [2], amount = 3
// 输出：-1
//
// 示例 3：
// 输入：coins = [1], amount = 0
// 输出：0

func coinChange(coins []int, amount int) int {
	coinMap := make(map[int]bool)
	dp := make([]int, amount+1)
	for _, coin := range coins {
		coinMap[coin] = true
	}
	for i := 0; i <= amount; i++ {
		if coinMap[i] {
			dp[i] = 1
		}
	}
	for i := 1; i <= amount; i++ {
		for j := 0; j < i; j++ {
			if coinMap[i-j] && dp[j] != 0 {
				if dp[i] == 0 {
					dp[i] = dp[j] + 1
				} else {
					dp[i] = min(dp[i], dp[j]+1)
				}
			}
		}
	}
	if dp[amount] == 0 && amount != 0 {
		return -1
	}
	return dp[amount]
}

func coinChange1(coins []int, amount int) int {
	inf := amount + 1
	dp := make([]int, amount+1)
	for i := range dp {
		dp[i] = inf
	}
	dp[0] = 0
	for i := 1; i <= amount; i++ {
		for _, coin := range coins {
			if coin <= i {
				dp[i] = min(dp[i], dp[i-coin]+1)
			}
		}
	}

	if dp[amount] == inf {
		return -1
	}

	return dp[amount]

}

func TestCoinChange(t *testing.T) {
	tests := []struct {
		name     string
		coins    []int
		amount   int
		expected int
	}{
		{name: "示例1", coins: []int{1, 2, 5}, amount: 11, expected: 3},
		{name: "示例2", coins: []int{2}, amount: 3, expected: -1},
		{name: "示例3", coins: []int{1}, amount: 0, expected: 0},
		{name: "大面额", coins: []int{1, 5, 10, 25}, amount: 30, expected: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := coinChange(tt.coins, tt.amount)
			if got != tt.expected {
				t.Errorf("coinChange() = %v, want %v", got, tt.expected)
			}
		})

		t.Run(tt.name, func(t *testing.T) {
			got := coinChange1(tt.coins, tt.amount)
			if got != tt.expected {
				t.Errorf("coinChange1() = %v, want %v", got, tt.expected)
			}
		})
	}
}

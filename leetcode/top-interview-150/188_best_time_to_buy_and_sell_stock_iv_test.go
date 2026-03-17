package topinterview150

import (
	"testing"
)

// 188. 买卖股票的最佳时机 IV (Best Time to Buy and Sell Stock IV)
//
// 题目描述:
// 给你一个整数数组 prices 和一个整数 k ，其中 prices[i] 是某支股票第 i 天的价格。
// 设计一个算法来计算你所能获取的最大利润。你最多可以完成 k 笔交易。也就是说，你最多可以买入 k 次，卖出 k 次。
// 注意：你不能同时参与多笔交易（你必须在再次购买前出售掉之前的股票）。
//
// 示例 1：
// 输入：k = 2, prices = [2,4,1]
// 输出：2
// 解释：在第 1 天 (价格 = 2) 买入，在第 2 天 (价格 = 4) 卖出，这笔交易所能获得利润 = 4-2 = 2 。
//
// 示例 2：
// 输入：k = 2, prices = [3,2,6,5,0,3]
// 输出：7
// 解释：在第 2 天 (价格 = 2) 买入，在第 3 天 (价格 = 6) 卖出, 这笔交易所能获得利润 = 6-2 = 4 。
//      随后，在第 5 天 (价格 = 0) 买入，在第 6 天 (价格 = 3) 卖出, 这笔交易所能获得利润 = 3-0 = 3 。

func maxProfitIV(k int, prices []int) int {
	panic("not implemented")
}

func TestMaxProfitIV(t *testing.T) {
	tests := []struct {
		name     string
		k        int
		prices   []int
		expected int
	}{
		{"Example 1", 2, []int{2, 4, 1}, 2},
		{"Example 2", 2, []int{3, 2, 6, 5, 0, 3}, 7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := maxProfitIV(tt.k, tt.prices); got != tt.expected {
			// 	t.Errorf("maxProfitIV() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

package topinterview150

import (
	"testing"
)

// 123. 买卖股票的最佳时机 III (Best Time to Buy and Sell Stock III)
//
// 题目描述:
// 给定一个数组 prices ，它的第 i 个元素 prices[i] 是一支给定的股票在第 i 天的价格。
// 设计一个算法来计算你所能获取的最大利润。你最多可以完成 两笔 交易。
// 注意：你不能同时参与多笔交易（你必须在再次购买前出售掉之前的股票）。
//
// 示例 1:
// 输入：prices = [3,3,5,0,0,3,1,4]
// 输出：6
// 解释：在第 4 天 (价格 = 0) 买入，在第 6 天 (价格 = 3) 卖出，这笔交易所能获得利润 = 3-0 = 3 。
//      随后，在第 7 天 (价格 = 1) 买入，在第 8 天 (价格 = 4) 卖出，这笔交易所能获得利润 = 4-1 = 3 。
//
// 示例 2：
// 输入：prices = [1,2,3,4,5]
// 输出：4
// 解释：在第 1 天 (价格 = 1) 买入，在第 5 天 (价格 = 5) 卖出, 这笔交易所能获得利润 = 5-1 = 4 。   
//      注意你不能在第 1 天和第 2 天接连购买股票，之后再将它们卖出。因为这样属于同时参与了多笔交易，你必须在再次购买前出售掉之前的股票。
//
// 示例 3：
// 输入：prices = [7,6,4,3,1]
// 输出：0
// 解释：在这个情况下, 没有交易完成, 所以最大利润为 0。

func maxProfitIII(prices []int) int {
	panic("not implemented")
}

func TestMaxProfitIII(t *testing.T) {
	tests := []struct {
		name     string
		prices   []int
		expected int
	}{
		{"Example 1", []int{3, 3, 5, 0, 0, 3, 1, 4}, 6},
		{"Example 2", []int{1, 2, 3, 4, 5}, 4},
		{"Example 3", []int{7, 6, 4, 3, 1}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := maxProfitIII(tt.prices); got != tt.expected {
			// 	t.Errorf("maxProfitIII() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

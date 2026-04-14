package top100liked

import (
	"testing"
)

// 121. 买卖股票的最佳时机 (Best Time to Buy and Sell Stock)
//
// 题目描述:
// 给定一个数组 prices ，它的第 i 个元素 prices[i] 表示一支给定股票第 i 天的价格。
// 你只能选择某一天买入这只股票，并选择在未来的某一个不同的日子卖出该股票。
// 设计一个算法来计算你所能获取的最大利润。返回你可以从这笔交易中获取的最大利润。
// 如果你不能获取任何利润，返回 0。
//
// 示例 1：
// 输入：[7,1,5,3,6,4]
// 输出：5
// 解释：在第 2 天（股票价格 = 1）的时候买入，在第 5 天（股票价格 = 6）的时候卖出，最大利润 = 6-1 = 5。
//
// 示例 2：
// 输入：prices = [7,6,4,3,1]
// 输出：0
// 解释：在这种情况下, 没有交易完成, 所以最大利润为 0。

func maxProfit(prices []int) int {
	minPrice := prices[0] // 记录遍历过程中的历史最低价
	maxProfit := 0        // 记录最大利润
	for _, p := range prices[1:] {
		// 以当前价卖出，看能否刷新最大利润
		if p-minPrice > maxProfit {
			maxProfit = p - minPrice
		}
		// 更新历史最低价
		if p < minPrice {
			minPrice = p
		}
	}
	return maxProfit
}

func TestMaxProfit(t *testing.T) {
	tests := []struct {
		name     string
		prices   []int
		expected int
	}{
		{
			name:     "示例1",
			prices:   []int{7, 1, 5, 3, 6, 4},
			expected: 5,
		},
		{
			name:     "示例2",
			prices:   []int{7, 6, 4, 3, 1},
			expected: 0,
		},
		{
			name:     "单个元素",
			prices:   []int{5},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := maxProfit(tt.prices)
			if got != tt.expected {
				t.Errorf("maxProfit() = %v, want %v", got, tt.expected)
			}
		})
	}
}

package topinterview150

import "testing"

// 121. 买卖股票的最佳时机 (Best Time to Buy and Sell Stock)
//
// 题目描述:
// 给定一个数组 prices ，它的第 i 个元素 prices[i] 表示一支给定股票第 i 天的价格。
// 你只能选择 "某一天" 买入这只股票，并选择在 "未来的某一个不同的日子" 卖出该股票。
// 设计一个算法来计算你所能获取的最大利润。
// 返回你可以从这笔交易中获取的最大利润。如果你不能获取任何利润，返回 0 。
//
// 示例 1：
// 输入：[7,1,5,3,6,4]
// 输出：5
//
// 示例 2：
// 输入：prices = [7,6,4,3,1]
// 输出：0

func maxProfit(prices []int) int {
	min, res := prices[0], 0
	for _, p := range prices {
		if min > p {
			min = p
		} else if res < p-min {
			res = p - min

		}
	}
	return res
}

func TestMaxProfit(t *testing.T) {
	tests := []struct {
		name     string
		prices   []int
		expected int
	}{
		{
			name:     "Example 1",
			prices:   []int{7, 1, 5, 3, 6, 4},
			expected: 5,
		},
		{
			name:     "Example 2",
			prices:   []int{7, 6, 4, 3, 1},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: Uncomment when implementation is ready
			if got := maxProfit(tt.prices); got != tt.expected {
				t.Errorf("maxProfit() = %v, want %v", got, tt.expected)
			}
		})
	}
}

package topinterview150

import (
	"testing"
)

// 135. 分发糖果 (Candy)
//
// 题目描述:
// n 个孩子站成一排。给你一个整数数组 ratings 表示每个孩子的评分。
// 你需要按照以下要求，给这些孩子分发糖果：
// 每个孩子至少分配到 1 个糖果。
// 相邻两个孩子评分更高的孩子会获得更多的糖果。
// 请弄清楚需要准备的 最少糖果数目 。
//
// 示例 1:
// 输入：ratings = [1,0,2]
// 输出：5
// 解释：你可以分别给第一个、第二个、第三个孩子分发 2、1、2 颗糖果。
//
// 示例 2:
// 输入：ratings = [1,2,2]
// 输出：4
// 解释：你可以分别给第一个、第二个、第三个孩子分发 1、2通1 颗糖果。
//      第三个孩子只得到 1 颗糖果，这满足题面定义的两个条件。

func candy(ratings []int) int {
	n := len(ratings)
	if n <= 1 {
		return n
	}

	// 初始化每个孩子一颗糖
	candies := make([]int, n)
	for i := range candies {
		candies[i] = 1
	}

	// 从左往右遍历
	// 满足左规则：如果右边评分更高，就在左边基础上 +1
	for i := 1; i < n; i++ {
		if ratings[i] > ratings[i-1] {
			candies[i] = candies[i-1] + 1
		}
	}

	// 从右往左遍历
	// 满足右规则：如果左边评分更高，且左边目前的糖果不比右边多
	for i := n - 2; i >= 0; i-- {
		if ratings[i] > ratings[i+1] && candies[i] <= candies[i+1] {
			candies[i] = candies[i+1] + 1
		}
	}

	res := 0
	for _, v := range candies {
		res += v
	}

	return res
}

func TestCandy(t *testing.T) {
	tests := []struct {
		name     string
		ratings  []int
		expected int
	}{
		{"Example 1", []int{1, 0, 2}, 5},
		{"Example 2", []int{1, 2, 2}, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := candy(tt.ratings); got != tt.expected {
			// 	t.Errorf("candy() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

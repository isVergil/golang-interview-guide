package topinterview150

import (
	"testing"
)

// 198. 打家劫舍 (House Robber)
//
// 题目描述:
// 你是一个专业的小偷，计划偷窃沿街的房屋。每间房内都藏有一定的现金，影响你偷窃的唯一制约因素就是相邻的房屋装有相互连通的防盗系统，如果两间相邻的房屋在同一晚上被小偷闯入，系统会自动报警。
// 给定一个代表每个房屋存放金额的非负整数数组，计算你 不触动警报装置的情况下 ，一夜之内能够偷窃到的最高金额。
//
// 示例 1：
// 输入：[1,2,3,1]
// 输出：4
// 解释：偷窃 1 号房屋 (金额 = 1) ，然后偷窃 3 号房屋 (金额 = 3)。
//      偷窃到的最高金额 = 1 + 3 = 4 。
//
// 示例 2：
// 输入：[2,7,9,3,1]
// 输出：12
// 解释：偷窃 1 号房屋 (金额 = 2), 偷窃 3 号房屋 (金额 = 9)，接着偷窃 5 号房屋 (金额 = 1)。
//      偷窃到的最高金额 = 2 + 9 + 1 = 12 。

func rob(nums []int) int {
	panic("not implemented")
}

func TestRob(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected int
	}{
		{"Example 1", []int{1, 2, 3, 1}, 4},
		{"Example 2", []int{2, 7, 9, 3, 1}, 12},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := rob(tt.nums); got != tt.expected {
			// 	t.Errorf("rob() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

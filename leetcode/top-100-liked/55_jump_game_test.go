package top100liked

import (
	"testing"
)

// 55. 跳跃游戏 (Jump Game)
//
// 题目描述:
// 给你一个非负整数数组 nums ，你最初位于数组的第一个下标。数组中的每个元素代表你在该位置可以跳跃的最大长度。
// 判断你是否能够到达最后一个下标。
//
// 示例 1：
// 输入：nums = [2,3,1,1,4]
// 输出：true
// 解释：可以先跳 1 步，从下标 0 到达下标 1, 然后再从下标 1 跳 3 步到达最后一个下标。
//
// 示例 2：
// 输入：nums = [3,2,1,0,4]
// 输出：false
// 解释：无论怎样，总会到达下标为 3 的位置。但该下标的最大跳跃长度是 0，所以永远不可能到达最后一个下标。

func canJump(nums []int) bool {
	step := 0

	// 不得不跳
	for i := 0; i <= step; i++ {
		if nums[i]+i >= step {
			step = nums[i] + i
		}
		if step >= len(nums)-1 {
			return true
		}
	}

	return false
}

func TestCanJump(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected bool
	}{
		{
			name:     "示例1",
			nums:     []int{2, 3, 1, 1, 4},
			expected: true,
		},
		{
			name:     "示例2",
			nums:     []int{3, 2, 1, 0, 4},
			expected: false,
		},
		{
			name:     "单元素",
			nums:     []int{0},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := canJump(tt.nums)
			if got != tt.expected {
				t.Errorf("canJump() = %v, want %v", got, tt.expected)
			}
		})
	}
}

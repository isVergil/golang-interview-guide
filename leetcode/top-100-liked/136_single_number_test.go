package top100liked

import (
	"testing"
)

// 136. 只出现一次的数字 (Single Number)
//
// 题目描述:
// 给你一个非空整数数组 nums ，除了某个元素只出现一次以外，其余每个元素均出现两次。
// 找出那个只出现了一次的元素。你必须设计并实现线性时间复杂度的算法，且该算法只使用常量额外空间。
//
// 示例 1：
// 输入：nums = [2,2,1]
// 输出：1
//
// 示例 2：
// 输入：nums = [4,1,2,1,2]
// 输出：4
//
// 示例 3：
// 输入：nums = [1]
// 输出：1
// 异或 相同为 0 不同为 1 
// a ^ a = 0
// a ^ 0 = a
func singleNumber(nums []int) int {
	res := nums[0]
	for i := 1; i < len(nums); i++ {
		res ^= nums[i]
	}
	return res
}

func TestSingleNumber(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected int
	}{
		{
			name:     "示例1",
			nums:     []int{2, 2, 1},
			expected: 1,
		},
		{
			name:     "示例2",
			nums:     []int{4, 1, 2, 1, 2},
			expected: 4,
		},
		{
			name:     "示例3",
			nums:     []int{1},
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := singleNumber(tt.nums)
			if got != tt.expected {
				t.Errorf("singleNumber() = %v, want %v", got, tt.expected)
			}
		})
	}
}

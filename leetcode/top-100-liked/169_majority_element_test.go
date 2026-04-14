package top100liked

import (
	"testing"
)

// 169. 多数元素 (Majority Element)
//
// 题目描述:
// 给定一个大小为 n 的数组 nums ，返回其中的多数元素。多数元素是指在数组中出现次数大于 ⌊n/2⌋ 的元素。
// 你可以假设数组是非空的，并且给定的数组总是存在多数元素。
//
// 示例 1：
// 输入：nums = [3,2,3]
// 输出：3
//
// 示例 2：
// 输入：nums = [2,2,1,1,1,2,2]
// 输出：2

func majorityElement(nums []int) int {
	// 摩尔投票法：候选人和票数
	res, voted := nums[0], 1
	for i := 1; i < len(nums); i++ {
		if voted == 0 {
			// 票数归零，换新候选人，自带一票
			res = nums[i]
			voted = 1
		} else if res == nums[i] {
			// 遇到相同元素，票数+1
			voted++
		} else {
			// 遇到不同元素，票数-1（互相抵消）
			voted--
		}
	}
	return res
}

func TestMajorityElement(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected int
	}{
		{
			name:     "示例1",
			nums:     []int{3, 2, 3},
			expected: 3,
		},
		{
			name:     "示例2",
			nums:     []int{2, 2, 1, 1, 1, 2, 2},
			expected: 2,
		},
		{
			name:     "单个元素",
			nums:     []int{1},
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := majorityElement(tt.nums)
			if got != tt.expected {
				t.Errorf("majorityElement() = %v, want %v", got, tt.expected)
			}
		})
	}
}

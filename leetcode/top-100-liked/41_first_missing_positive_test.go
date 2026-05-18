package top100liked

import "testing"

// 41. 缺失的第一个正数 (First Missing Positive)
//
// 题目描述:
// 给你一个未排序的整数数组 nums，请你找出其中没有出现的最小的正整数。
// 要求时间复杂度 O(n)，空间复杂度 O(1)。
//
// 示例 1：
// 输入：nums = [1,2,0]
// 输出：3
//
// 示例 2：
// 输入：nums = [3,4,-1,1]
// 输出：2
//
// 示例 3：
// 输入：nums = [7,8,9,11,12]
// 输出：1
//
// 提示：原地哈希，把数字 x 放到下标 x-1 的位置上

func firstMissingPositive(nums []int) int {
	n := len(nums)
	for i := 0; i < n; i++ {
		for nums[i] > 0 && nums[i] <= n && nums[i] != nums[nums[i]-1] {
			nums[i], nums[nums[i]-1] = nums[nums[i]-1], nums[i]
		}
	}
	for i := 0; i < n; i++ {
		if nums[i] != i+1 {
			return i + 1
		}
	}
	return n + 1
}

func TestFirstMissingPositive(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected int
	}{
		{name: "示例1", nums: []int{1, 2, 0}, expected: 3},
		{name: "示例2", nums: []int{3, 4, -1, 1}, expected: 2},
		{name: "示例3", nums: []int{7, 8, 9, 11, 12}, expected: 1},
		{name: "连续", nums: []int{1, 2, 3, 4, 5}, expected: 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := firstMissingPositive(tt.nums)
			if got != tt.expected {
				t.Errorf("firstMissingPositive() = %v, want %v", got, tt.expected)
			}
		})
	}
}

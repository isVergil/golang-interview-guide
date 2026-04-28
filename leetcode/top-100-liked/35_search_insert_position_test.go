package top100liked

import (
	"testing"
)

// 35. 搜索插入位置 (Search Insert Position)
//
// 题目描述:
// 给定一个排序数组和一个目标值，在数组中找到目标值，并返回其索引。
// 如果目标值不存在于数组中，返回它将会被按顺序插入的位置。
// 请必须使用时间复杂度为 O(log n) 的算法。
//
// 示例 1：
// 输入：nums = [1,3,5,6], target = 5
// 输出：2
//
// 示例 2：
// 输入：nums = [1,3,5,6], target = 2
// 输出：1
//
// 示例 3：
// 输入：nums = [1,3,5,6], target = 7
// 输出：4

func searchInsert(nums []int, target int) int {
	l, r := 0, len(nums)-1
	for l <= r {

		mid := l + (r-l)/2
		if nums[mid] > target {
			r = mid - 1
		} else if nums[mid] < target {
			l = mid + 1
		} else {
			return mid
		}
	}

	// l 的含义：始终指向"第一个 >= target 的位置"
	return l
}

func TestSearchInsert(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		target   int
		expected int
	}{
		{name: "示例1", nums: []int{1, 3, 5, 6}, target: 5, expected: 2},
		{name: "示例2", nums: []int{1, 3, 5, 6}, target: 2, expected: 1},
		{name: "示例3", nums: []int{1, 3, 5, 6}, target: 7, expected: 4},
		{name: "插入头部", nums: []int{1, 3, 5, 6}, target: 0, expected: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := searchInsert(tt.nums, tt.target)
			if got != tt.expected {
				t.Errorf("searchInsert() = %v, want %v", got, tt.expected)
			}
		})
	}
}

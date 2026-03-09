package topinterview150

import (
	"testing"
)

// 35. 搜索插入位置 (Search Insert Position)
//
// 题目描述:
// 给定一个排序数组和一个目标值，在数组中找到目标值，并返回其索引。如果目标值不存在于数组中，返回它将会被按顺序插入的位置。
// 请必须使用时间复杂度为 O(log n) 的算法。
//
// 示例 1:
// 输入: nums = [1,3,5,6], target = 5
// 输出: 2
//
// 示例 2:
// 输入: nums = [1,3,5,6], target = 2
// 输出: 1
//
// 示例 3:
// 输入: nums = [1,3,5,6], target = 7
// 输出: 4

func searchInsert(nums []int, target int) int {
	panic("not implemented")
}

func TestSearchInsert(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		target   int
		expected int
	}{
		{"Example 1", []int{1, 3, 5, 6}, 5, 2},
		{"Example 2", []int{1, 3, 5, 6}, 2, 1},
		{"Example 3", []int{1, 3, 5, 6}, 7, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := searchInsert(tt.nums, tt.target); got != tt.expected {
			// 	t.Errorf("searchInsert() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

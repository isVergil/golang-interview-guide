package top100liked

import (
	"reflect"
	"testing"
)

// 34. 在排序数组中查找元素的第一个和最后一个位置 (Find First and Last Position of Element in Sorted Array)
//
// 题目描述:
// 给你一个按照非递减顺序排列的整数数组 nums，和一个目标值 target。
// 请你找出给定目标值在数组中的开始位置和结束位置。
// 如果数组中不存在目标值 target，返回 [-1, -1]。你必须设计并实现时间复杂度为 O(log n) 的算法。
//
// 示例 1：
// 输入：nums = [5,7,7,8,8,10], target = 8
// 输出：[3,4]
//
// 示例 2：
// 输入：nums = [5,7,7,8,8,10], target = 6
// 输出：[-1,-1]
//
// 示例 3：
// 输入：nums = [], target = 0
// 输出：[-1,-1]

func searchRange(nums []int, target int) []int {
	return []int{findFirst(nums, target), findLast(nums, target)}
}

// findFirst 找第一个等于 target 的位置
func findFirst(nums []int, target int) int {
	l, r := 0, len(nums)-1
	res := -1
	for l <= r {
		mid := l + (r-l)/2
		if nums[mid] < target {
			l = mid + 1
		} else if nums[mid] > target {
			r = mid - 1
		} else {
			// 找到了，但不急着返回，记录答案，继续往左找
			res = mid
			r = mid - 1
		}
	}
	return res
}

// findLast 找最后一个等于 target 的位置
func findLast(nums []int, target int) int {
	l, r := 0, len(nums)-1
	res := -1
	for l <= r {
		mid := l + (r-l)/2
		if nums[mid] < target {
			l = mid + 1
		} else if nums[mid] > target {
			r = mid - 1
		} else {
			// 找到了，记录答案，继续往右找
			res = mid
			l = mid + 1
		}
	}
	return res
}

func TestSearchRange(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		target   int
		expected []int
	}{
		{name: "示例1", nums: []int{5, 7, 7, 8, 8, 10}, target: 8, expected: []int{3, 4}},
		{name: "示例2", nums: []int{5, 7, 7, 8, 8, 10}, target: 6, expected: []int{-1, -1}},
		{name: "示例3", nums: []int{}, target: 0, expected: []int{-1, -1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := searchRange(tt.nums, tt.target)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("searchRange() = %v, want %v", got, tt.expected)
			}
		})
	}
}

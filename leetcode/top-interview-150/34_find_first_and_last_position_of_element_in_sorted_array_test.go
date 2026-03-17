package topinterview150

import (
	"testing"
)

// 34. 在排序数组中查找元素的第一个和最后一个位置 (Find First and Last Position of Element in Sorted Array)
//
// 题目描述:
// 给你一个按照非递减顺序排列的整数数组 nums，和一个目标值 target。请你找出给定目标值在数组中的开始位置和结束位置。
// 如果数组中不存在目标值 target，返回 [-1, -1]。
// 你必须设计并实现时间复杂度为 O(log n) 的算法解决此问题。
//
// 示例 1：
// 输入：nums = [5,7,7,8,8,10], target = 8
// 输出：[3,4]

func searchRange(nums []int, target int) []int {
	// 初始结果
	res := []int{-1, -1}

	// 找第一个位置
	first := binarySearch(nums, target, true)
	// 如果连第一个都没找到，直接返回 [-1, -1]
	if first == -1 {
		return res
	}

	// 找最后一个位置
	last := binarySearch(nums, target, false)

	return []int{first, last}
}

// leftBound 为 true 找第一个，为 false 找最后一个
func binarySearch(nums []int, target int, leftBound bool) int {
	left, right := 0, len(nums)-1
	res := -1

	for left <= right {
		mid := left + (right-left)/2
		if nums[mid] == target {
			res = mid // 记录当前位置
			if leftBound {
				right = mid - 1 // 尝试往左找更早的
			} else {
				left = mid + 1 // 尝试往右找更晚的
			}
		} else if nums[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return res
}

func TestSearchRange(t *testing.T) {
	// 查找元素位置测试
}

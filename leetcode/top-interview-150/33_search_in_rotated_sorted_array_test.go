package topinterview150

import (
	"testing"
)

// 33. 搜索旋转排序数组 (Search in Rotated Sorted Array)
//
// 题目描述:
// 整数数组 nums 按升序排列，数组中的值 互不相同 。
// 在传递给函数之前，nums 在预先未知的某个下标 k（0 <= k < nums.length）上进行了 旋转。
// 给你 旋转后 的数组 nums 和一个整数 target ，如果 nums 中存在这个目标值 target ，则返回它的下标，否则返回 -1 。
// 你必须设计一个时间复杂度为 O(log n) 的算法解决此问题。
//
// 示例 1：
// 输入：nums = [4,5,6,7,0,1,2], target = 0
// 输出：4

func search(nums []int, target int) int {
	left, right := 0, len(nums)-1

	for left <= right {
		mid := left + (right-left)/2 // 细节：防止溢出
		if nums[mid] == target {
			return mid
		}

		// 1. 判断左半部分是否有序
		if nums[left] <= nums[mid] {
			// 如果 target 在左侧有序区间内
			if nums[left] <= target && target < nums[mid] {
				right = mid - 1
			} else {
				left = mid + 1
			}
		} else {
			// 2. 否则，右半部分必然是有序的
			// 如果 target 在右侧有序区间内
			if nums[mid] < target && target <= nums[right] {
				left = mid + 1
			} else {
				right = mid - 1
			}
		}
	}

	return -1
}

func TestSearch(t *testing.T) {
	// 搜索旋转排序数组测试
}

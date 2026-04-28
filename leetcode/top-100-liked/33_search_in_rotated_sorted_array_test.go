package top100liked

import (
	"testing"
)

// 33. 搜索旋转排序数组 (Search in Rotated Sorted Array)
//
// 题目描述:
// 整数数组 nums 按升序排列，数组中的值互不相同。在传递给函数之前，nums 在预先未知的某个下标 k 上进行了旋转。
// 给你旋转后的数组 nums 和一个整数 target ，如果 nums 中存在这个目标值，则返回它的下标，否则返回 -1 。
// 你必须设计一个时间复杂度为 O(log n) 的算法解决此问题。
//
// 示例 1：
// 输入：nums = [4,5,6,7,0,1,2], target = 0
// 输出：4
//
// 示例 2：
// 输入：nums = [4,5,6,7,0,1,2], target = 3
// 输出：-1
//
// 示例 3：
// 输入：nums = [1], target = 0
// 输出：-1
// 二分查找：每次判断哪半边有序，缩小范围
// 每次二分，找到有序的那半边 → 能精确判断 target 在不在 → 在就缩到那边，不在就去另一边。O(log n) 次就找到了。
func search(nums []int, target int) int {
	left, right := 0, len(nums)-1
	for left <= right {
		mid := left + (right-left)/2
		if target == nums[mid] {
			return mid
		}
		// 左半边有序，单个也是有序的 所以用=号
		if nums[mid] >= nums[left] {
			// target 在左半边范围内：target 可能就在 left 位置，要包含，mid已经判断过了不用=
			if target >= nums[left] && target < nums[mid] {
				right = mid - 1
			} else {
				left = mid + 1
			}
		} else {
			// 右半边有序，跟上面类似
			if target > nums[mid] && target <= nums[right] {
				left = mid + 1
			} else {
				right = mid - 1
			}
		}
	}
	return -1
}

func TestSearch(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		target   int
		expected int
	}{
		{
			name:     "示例1",
			nums:     []int{4, 5, 6, 7, 0, 1, 2},
			target:   0,
			expected: 4,
		},
		{
			name:     "示例2",
			nums:     []int{4, 5, 6, 7, 0, 1, 2},
			target:   3,
			expected: -1,
		},
		{
			name:     "示例3",
			nums:     []int{1},
			target:   0,
			expected: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := search(tt.nums, tt.target)
			if got != tt.expected {
				t.Errorf("search() = %v, want %v", got, tt.expected)
			}
		})
	}
}

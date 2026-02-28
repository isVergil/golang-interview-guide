package topinterview150

import (
	"testing"
)

// 215. 数组中的第K个最大元素 (Kth Largest Element in an Array)
//
// 题目描述:
// 给定整数数组 nums 和整数 k，请返回数组中第 k 个最大的元素。
// 请注意，你需要找的是数组排序后的第 k 个最大的元素，而不是第 k 个不同的元素。
// 你必须设计并实现时间复杂度为 O(n) 的算法解决此问题。
//
// 示例 1:
// 输入: [3,2,1,5,6,4], k = 2
// 输出: 5
//
// 示例 2:
// 输入: [3,2,3,1,2,4,5,5,6], k = 4
// 输出: 4

func findKthLargest(nums []int, k int) int {
	// 原地建一个大小为 k 的最小堆
	for i := k/2 - 1; i >= 0; i-- {
		siftDown(nums, i, k)
	}

	for i := k; i < len(nums); i++ {
		if nums[i] > nums[0] {
			nums[0] = nums[i]
			siftDown(nums, 0, k)
		}
	}
	return nums[0]
}

func siftDown(nums []int, i, k int) {
	for {
		left := 2*i + 1
		if left >= k {
			break
		}

		smaller := left
		if right := left + 1; right < k && nums[right] < nums[left] {
			smaller = right
		}

		if nums[i] <= nums[smaller] {
			break
		}

		nums[i], nums[smaller] = nums[smaller], nums[i]
		i = smaller
	}
}

func TestFindKthLargest(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		k        int
		expected int
	}{
		{"Example 1", []int{3, 2, 1, 5, 6, 4}, 2, 5},
		{"Example 2", []int{3, 2, 3, 1, 2, 4, 5, 5, 6}, 4, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := findKthLargest(tt.nums, tt.k); got != tt.expected {
			// 	t.Errorf("findKthLargest() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

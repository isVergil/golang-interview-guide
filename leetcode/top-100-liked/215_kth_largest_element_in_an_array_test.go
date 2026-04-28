package top100liked

import (
	"math/rand"
	"testing"
)

// 215. 数组中的第K个最大元素 (Kth Largest Element in an Array)
//
// 题目描述:
// 给定整数数组 nums 和整数 k，请返回数组中第 k 个最大的元素。
// 请注意，你需要找的是数组排序后的第 k 个最大的元素，而不是第 k 个不同的元素。
// 你必须设计并实现时间复杂度为 O(n) 的算法解决此问题。
//
// 示例 1：
// 输入：nums = [3,2,1,5,6,4], k = 2
// 输出：5
//
// 示例 2：
// 输入：nums = [3,2,3,1,2,4,5,5,6], k = 4
// 输出：4

func findKthLargest(nums []int, k int) int {
	idx := len(nums) - k
	l, r := 0, len(nums)-1
	for l < r {
		pivot := partition(nums, l, r)
		if pivot < idx {
			l = pivot + 1
		} else if pivot > idx {
			r = pivot - 1
		} else {
			return nums[pivot]
		}
	}
	return nums[l]
}

func partition(nums []int, left, right int) int {
	randIdx := left + rand.Intn(right-left+1)
	nums[randIdx], nums[right] = nums[right], nums[randIdx]
	pivot := nums[right]
	i := left
	for j := left; j < right; j++ {
		if nums[j] <= pivot {
			nums[i], nums[j] = nums[j], nums[i]
			i++
		}
	}
	nums[i], nums[right] = nums[right], nums[i]
	return i

}

func TestFindKthLargest(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		k        int
		expected int
	}{
		{name: "示例1", nums: []int{3, 2, 1, 5, 6, 4}, k: 2, expected: 5},
		{name: "示例2", nums: []int{3, 2, 3, 1, 2, 4, 5, 5, 6}, k: 4, expected: 4},
		{name: "单元素", nums: []int{1}, k: 1, expected: 1},
		{name: "全相同", nums: []int{2, 2, 2}, k: 2, expected: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findKthLargest(tt.nums, tt.k)
			if got != tt.expected {
				t.Errorf("findKthLargest() = %v, want %v", got, tt.expected)
			}
		})
	}
}

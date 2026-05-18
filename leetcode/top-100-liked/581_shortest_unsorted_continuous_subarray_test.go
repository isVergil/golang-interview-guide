package top100liked

import (
	"math"
	"testing"
)

// 581. 最短无序连续子数组 (Shortest Unsorted Continuous Subarray)
//
// 题目描述:
// 给你一个整数数组 nums，你需要找出一个连续子数组，
// 如果对这个子数组进行升序排序，那么整个数组都会变为升序排序。
// 请你找出符合题意的最短子数组，并输出它的长度。
//
// 示例 1：
// 输入：nums = [2,6,4,8,10,9,15]
// 输出：5
// 解释：对 [6, 4, 8, 10, 9] 排序即可让整个数组有序
//
// 示例 2：
// 输入：nums = [1,2,3,4]
// 输出：0
//
// 提示：从左找右边界（最后一个比左侧最大值小的），从右找左边界（最后一个比右侧最小值大的）

func findUnsortedSubarray(nums []int) int {
	n := len(nums)
	l, r := -1, -2
	maxVal, minVal := math.MinInt, math.MaxInt
	for i := 0; i < n; i++ {
		if nums[i] < maxVal {
			r = i
		} else {
			maxVal = nums[i]
		}

		j := n - 1 - i
		if nums[j] > minVal {
			l = j
		} else {
			minVal = nums[j]
		}
	}
	return r - l + 1
}

func TestFindUnsortedSubarray(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected int
	}{
		{name: "示例1", nums: []int{2, 6, 4, 8, 10, 9, 15}, expected: 5},
		{name: "有序", nums: []int{1, 2, 3, 4}, expected: 0},
		{name: "逆序", nums: []int{3, 2, 1}, expected: 3},
		{name: "单元素", nums: []int{1}, expected: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findUnsortedSubarray(tt.nums)
			if got != tt.expected {
				t.Errorf("findUnsortedSubarray() = %v, want %v", got, tt.expected)
			}
		})
	}
}

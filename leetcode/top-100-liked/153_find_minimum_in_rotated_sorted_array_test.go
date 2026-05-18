package top100liked

import "testing"

// 153. 寻找旋转排序数组中的最小值 (Find Minimum in Rotated Sorted Array)
//
// 题目描述:
// 已知一个长度为 n 的数组，预先按照升序排列，经由 1 到 n 次旋转后得到输入数组。
// 给你一个元素值互不相同的数组 nums，请你找出并返回数组中的最小元素。
// 要求时间复杂度 O(log n)。
//
// 示例 1：
// 输入：nums = [3,4,5,1,2]
// 输出：1
//
// 示例 2：
// 输入：nums = [4,5,6,7,0,1,2]
// 输出：0
//
// 示例 3：
// 输入：nums = [11,13,15,17]
// 输出：11
//
// 提示：二分搜索，比较 mid 和右端点判断最小值在哪半边

func findMin(nums []int) int {
	l, r := 0, len(nums)-1
	for l < r {
		mid := l + (r-l)/2
		if nums[mid] > nums[r] {
			// mid 在左段，最小值一定在 mid 右边
			l = mid + 1
		} else {
			// mid 在右段，mid 可能就是最小值
			r = mid
		}
	}
	return nums[l]
}

func TestFindMin(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected int
	}{
		{name: "示例1", nums: []int{3, 4, 5, 1, 2}, expected: 1},
		{name: "示例2", nums: []int{4, 5, 6, 7, 0, 1, 2}, expected: 0},
		{name: "未旋转", nums: []int{11, 13, 15, 17}, expected: 11},
		{name: "两个元素", nums: []int{2, 1}, expected: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findMin(tt.nums)
			if got != tt.expected {
				t.Errorf("findMin() = %v, want %v", got, tt.expected)
			}
		})
	}
}

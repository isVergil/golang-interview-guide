package topinterview150

import (
	"testing"
)

// 153. 寻找旋转排序数组中的最小值 (Find Minimum in Rotated Sorted Array)
//
// 题目描述:
// 已知一个长度为 n 的数组，预先按照升序排列，经由 1 到 n 次 旋转 后，得到输入数组。
// 例如，原数组 nums = [0,1,2,4,5,6,7] 在变化后可能得到：
// 若旋转 4 次，则可以得到 [4,5,6,7,0,1,2]
// 若旋转 7 次，则可以得到 [0,1,2,4,5,6,7]
// 给你一个元素值 互不相同 的数组 nums ，它原来是一个升序排列的数组，并按上述情形进行了多次旋转。请你找出并返回数组中的 最小元素 。
// 你必须设计一个时间复杂度为 O(log n) 的算法解决此问题。
//
// 示例 1：
// 输入：nums = [3,4,5,1,2]
// 输出：1
// 解释：原数组为 [1,2,3,4,5] ，旋转 3 次得到输入数组。
//
// 示例 2：
// 输入：nums = [4,5,6,7,0,1,2]
// 输出：0
// 解释：原数组为 [0,1,2,4,5,6,7] ，旋转 4 次得到输入数组。
//
// 示例 3：
// 输入：nums = [11,13,15,17]
// 输出：11
// 解释：原数组为 [11,13,15,17] ，旋转 4 次得到输入数组。

func findMin(nums []int) int {
	panic("not implemented")
}

func TestFindMin(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected int
	}{
		{"Example 1", []int{3, 4, 5, 1, 2}, 1},
		{"Example 2", []int{4, 5, 6, 7, 0, 1, 2}, 0},
		{"Example 3", []int{11, 13, 15, 17}, 11},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := findMin(tt.nums); got != tt.expected {
			// 	t.Errorf("findMin() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

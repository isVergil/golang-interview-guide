package top100liked

import (
	"reflect"
	"testing"
)

// 31. 下一个排列 (Next Permutation)
//
// 题目描述:
// 整数数组的一个排列就是将其所有成员以序列或线性顺序排列。
// 实现获取下一个排列的函数，将数字重新排列成字典序中下一个更大的排列。
// 如果不存在下一个更大的排列，则将数字重新排列成最小的排列（即升序排列）。
// 必须原地修改，只允许使用额外常数空间。
//
// 示例 1：
// 输入：nums = [1,2,3]
// 输出：[1,3,2]
//
// 示例 2：
// 输入：nums = [3,2,1]
// 输出：[1,2,3]
//
// 示例 3：
// 输入：nums = [1,1,5]
// 输出：[1,5,1]

func nextPermutation(nums []int) {
	n := len(nums)

	// 从右往左 找第一个下降的位置
	i := n - 2
	for i >= 0 && nums[i] >= nums[i+1] {
		i--
	}

	// 找到了 从右往左 找第一个大于 nums[i] 的元素
	if i >= 0 {
		j := n - 1
		for nums[j] <= nums[i] {
			j--
		}

		// 交换
		nums[i], nums[j] = nums[j], nums[i]
	}

	l, r := i+1, n-1
	for l < r {
		nums[l], nums[r] = nums[r], nums[l]
		l++
		r--
	}
}

func TestNextPermutation(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected []int
	}{
		{name: "示例1", nums: []int{1, 2, 3}, expected: []int{1, 3, 2}},
		{name: "示例2", nums: []int{3, 2, 1}, expected: []int{1, 2, 3}},
		{name: "示例3", nums: []int{1, 1, 5}, expected: []int{1, 5, 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextPermutation(tt.nums)
			if !reflect.DeepEqual(tt.nums, tt.expected) {
				t.Errorf("nextPermutation() = %v, want %v", tt.nums, tt.expected)
			}
		})
	}
}

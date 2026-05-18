package top100liked

import (
	"reflect"
	"testing"
)

// 189. 轮转数组 (Rotate Array)
//
// 题目描述:
// 给定一个整数数组 nums，将数组中的元素向右轮转 k 个位置。
//
// 示例 1：
// 输入：nums = [1,2,3,4,5,6,7], k = 3
// 输出：[5,6,7,1,2,3,4]
//
// 示例 2：
// 输入：nums = [-1,-100,3,99], k = 2
// 输出：[3,99,-1,-100]
//
// 提示：三次反转法 - 全部反转，反转前k个，反转后n-k个

func rotateArray(nums []int, k int) {
	n := len(nums)
	k %= n
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		nums[i], nums[j] = nums[j], nums[i]
	}
	for i, j := 0, k-1; i < j; i, j = i+1, j-1 {
		nums[i], nums[j] = nums[j], nums[i]
	}
	for i, j := k, n-1; i < j; i, j = i+1, j-1 {
		nums[i], nums[j] = nums[j], nums[i]
	}
}

func TestRotateArray(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		k        int
		expected []int
	}{
		{name: "示例1", nums: []int{1, 2, 3, 4, 5, 6, 7}, k: 3, expected: []int{5, 6, 7, 1, 2, 3, 4}},
		{name: "示例2", nums: []int{-1, -100, 3, 99}, k: 2, expected: []int{3, 99, -1, -100}},
		{name: "k大于n", nums: []int{1, 2, 3}, k: 4, expected: []int{3, 1, 2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rotateArray(tt.nums, tt.k)
			if !reflect.DeepEqual(tt.nums, tt.expected) {
				t.Errorf("rotate() = %v, want %v", tt.nums, tt.expected)
			}
		})
	}
}

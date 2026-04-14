package top100liked

import (
	"reflect"
	"testing"
)

// 238. 除自身以外数组的乘积 (Product of Array Except Self)
//
// 题目描述:
// 给你一个整数数组 nums，返回数组 answer ，其中 answer[i] 等于 nums 中除 nums[i] 之外其余各元素的乘积。
// 题目数据保证数组 nums 之中任意元素的全部前缀元素和后缀的乘积都在 32 位整数范围内。
// 请不要使用除法，且在 O(n) 时间复杂度内完成此题。
//
// 示例 1：
// 输入: nums = [1,2,3,4]
// 输出: [24,12,8,6]
//
// 示例 2：
// 输入: nums = [-1,1,0,-3,3]
// 输出: [0,0,9,0,0]

func productExceptSelf(nums []int) []int {
	n := len(nums)
	res := make([]int, n)

	// 正向遍历，res[i] 存左侧所有元素的乘积
	res[0] = 1
	for i := 1; i < n; i++ {
		res[i] = res[i-1] * nums[i-1]
	}

	// 反向遍历，用一个变量累积右侧乘积，直接乘进 res
	right := 1
	for i := n - 2; i >= 0; i-- {
		right *= nums[i+1]
		res[i] *= right
	}

	return res
}

func TestProductExceptSelf(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected []int
	}{
		{
			name:     "示例1",
			nums:     []int{1, 2, 3, 4},
			expected: []int{24, 12, 8, 6},
		},
		{
			name:     "示例2",
			nums:     []int{-1, 1, 0, -3, 3},
			expected: []int{0, 0, 9, 0, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := productExceptSelf(tt.nums)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("productExceptSelf() = %v, want %v", got, tt.expected)
			}
		})
	}
}

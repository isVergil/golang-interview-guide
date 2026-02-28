package topinterview150

import (
	"testing"
)

// 238. 除自身以外数组的乘积 (Product of Array Except Self)
//
// 题目描述:
// 给你一个整数数组 nums，返回 数组 answer ，其中 answer[i] 等于 nums 中除 nums[i] 之外其余各元素的乘积 。
// 题目数据 保证 数组 answer之中任意元素的全部前缀乘积和后缀乘积都在  32 位 整数范围内。
// 请 不要使用除法，且在 O(n) 时间复杂度内完成此题。
//
// 示例 1:
// 输入: nums = [1,2,3,4]
// 输出: [24,12,8,6]
//
// 示例 2:
// 输入: nums = [-1,1,0,-3,3]
// 输出: [0,0,9,0,0]

func productExceptSelf(nums []int) []int {
	n := len(nums)
	res := make([]int, n)

	// 1 从左往右遍历，计算前缀积
	res[0] = 1
	for i := 1; i < n; i++ {
		res[i] = res[i-1] * nums[i-1]
	}

	// 2 从右往左遍历：计算后缀积并直接乘以之前的结果
	right := 1
	for i := n - 1; i >= 0; i-- {
		res[i] *= right
		right *= nums[i]
	}

	return res
}

func TestProductExceptSelf(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected []int
	}{
		{"Example 1", []int{1, 2, 3, 4}, []int{24, 12, 8, 6}},
		{"Example 2", []int{-1, 1, 0, -3, 3}, []int{0, 0, 9, 0, 0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// got := productExceptSelf(tt.nums)
			// if !reflect.DeepEqual(got, tt.expected) {
			// 	t.Errorf("productExceptSelf() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

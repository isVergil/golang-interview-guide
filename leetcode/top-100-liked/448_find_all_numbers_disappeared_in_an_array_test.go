package top100liked

import (
	"testing"
)

// 448. 找到所有数组中消失的数字 (Find All Numbers Disappeared in an Array)
//
// 题目描述:
// 给你一个含 n 个整数的数组 nums ，其中 nums[i] 在区间 [1, n] 内。请你找出所有在 [1, n] 范围内但没有出现在 nums 中的数字，并以数组的形式返回结果。
//
// 示例 1：
// 输入：nums = [4,3,2,7,8,2,3,1]
// 输出：[5,6]
//
// 示例 2：
// 输入：nums = [1,1]
// 输出：[2]

func findDisappearedNumbers(nums []int) []int {
	panic("not implemented")
}

func TestFindDisappearedNumbers(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected []int
	}{
		{"Example 1", []int{4, 3, 2, 7, 8, 2, 3, 1}, []int{5, 6}},
		{"Example 2", []int{1, 1}, []int{2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := findDisappearedNumbers(tt.nums); !reflect.DeepEqual(got, tt.expected) {
			// 	t.Errorf("findDisappearedNumbers() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

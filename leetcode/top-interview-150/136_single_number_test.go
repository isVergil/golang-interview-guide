package topinterview150

import (
	"testing"
)

// 136. 只出现一次的数字 (Single Number)
//
// 题目描述:
// 给你一个 非空 整数数组 nums ，除了某个元素只出现一次以外，其余每个元素均出现两次。找出那个只出现了一次的元素。
// 你必须设计并实现线性时间复杂度的算法来解决此问题，且该算法只使用常量额外空间。
//
// 示例 1 ：
// 输入：nums = [2,2,1]
// 输出：1
//
// 示例 2 ：
// 输入：nums = [4,1,2,1,2]
// 输出：4
//
// 示例 3 ：
// 输入：nums = [1]
// 输出：1

func singleNumber(nums []int) int {
	panic("not implemented")
}

func TestSingleNumber(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected int
	}{
		{"Example 1", []int{2, 2, 1}, 1},
		{"Example 2", []int{4, 1, 2, 1, 2}, 4},
		{"Example 3", []int{1}, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := singleNumber(tt.nums); got != tt.expected {
			// 	t.Errorf("singleNumber() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

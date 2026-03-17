package topinterview150

import (
	"testing"
)

// 137. 只出现一次的数字 II (Single Number II)
//
// 题目描述:
// 给你一个整数数组 nums ，除某个元素仅出现 一次 外，其余每个元素都恰出现 三次 。请你找出并返回那个只出现了一次的元素。
// 你必须设计并实现线性时间复杂度的算法来解决此问题，且该算法只使用常量额外空间。
//
// 示例 1：
// 输入：nums = [2,2,3,2]
// 输出：3
//
// 示例 2：
// 输入：nums = [0,1,0,1,0,1,99]
// 输出：99

func singleNumberII(nums []int) int {
	panic("not implemented")
}

func TestSingleNumberII(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected int
	}{
		{"Example 1", []int{2, 2, 3, 2}, 3},
		{"Example 2", []int{0, 1, 0, 1, 0, 1, 99}, 99},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := singleNumberII(tt.nums); got != tt.expected {
			// 	t.Errorf("singleNumberII() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

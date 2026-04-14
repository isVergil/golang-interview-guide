package top100liked

import (
	"reflect"
	"testing"
)

// 283. 移动零 (Move Zeroes)
//
// 题目描述:
// 给定一个数组 nums，编写一个函数将所有 0 移动到数组的末尾，同时保持非零元素的相对顺序。
// 请注意，必须在不复制数组的情况下原地对数组进行操作。
//
// 示例 1：
// 输入: nums = [0,1,0,3,12]
// 输出: [1,3,12,0,0]
//
// 示例 2：
// 输入: nums = [0]
// 输出: [0]

func moveZeroes(nums []int) {
	// slow 指向下一个非零元素应该放的位置
	slow := 0
	for fast := 0; fast < len(nums); fast++ {
		if nums[fast] != 0 {
			// 交换，将非零元素换到前面
			nums[slow], nums[fast] = nums[fast], nums[slow]
			slow++
		}
	}
}

func TestMoveZeroes(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected []int
	}{
		{
			name:     "示例1",
			nums:     []int{0, 1, 0, 3, 12},
			expected: []int{1, 3, 12, 0, 0},
		},
		{
			name:     "示例2",
			nums:     []int{0},
			expected: []int{0},
		},
		{
			name:     "无零",
			nums:     []int{1, 2, 3},
			expected: []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			moveZeroes(tt.nums)
			if !reflect.DeepEqual(tt.nums, tt.expected) {
				t.Errorf("moveZeroes() = %v, want %v", tt.nums, tt.expected)
			}
		})
	}
}

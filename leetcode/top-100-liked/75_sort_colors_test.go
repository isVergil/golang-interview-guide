package top100liked

import (
	"reflect"
	"testing"
)

// 75. 颜色分类 (Sort Colors)
//
// 题目描述:
// 给定一个包含红色、白色和蓝色、共 n 个元素的数组 nums ，原地对它们进行排序，
// 使得相同颜色的元素相邻，并按照红色、白色、蓝色顺序排列。
// 我们使用整数 0、1 和 2 分别表示红色、白色和蓝色。必须在不使用库内置的 sort 函数的情况下解决这个问题。
//
// 示例 1：
// 输入：nums = [2,0,2,1,1,0]
// 输出：[0,0,1,1,2,2]
//
// 示例 2：
// 输入：nums = [2,0,1]
// 输出：[0,1,2]

func sortColors(nums []int) {
	i, l, r := 0, 0, len(nums)-1
	for i <= r {
		switch nums[i] {
		case 0:
			nums[i], nums[l] = nums[l], nums[i]
			l++
			i++
		case 2:
			nums[i], nums[r] = nums[r], nums[i]
			r--
		default:
			i++
		}
	}
}

func TestSortColors(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected []int
	}{
		{
			name:     "示例1",
			nums:     []int{2, 0, 2, 1, 1, 0},
			expected: []int{0, 0, 1, 1, 2, 2},
		},
		{
			name:     "示例2",
			nums:     []int{2, 0, 1},
			expected: []int{0, 1, 2},
		},
		{
			name:     "全相同",
			nums:     []int{1, 1, 1},
			expected: []int{1, 1, 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sortColors(tt.nums)
			if !reflect.DeepEqual(tt.nums, tt.expected) {
				t.Errorf("sortColors() = %v, want %v", tt.nums, tt.expected)
			}
		})
	}
}

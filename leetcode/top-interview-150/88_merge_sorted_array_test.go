package topinterview150

import (
	"reflect"
	"testing"
)

// 88. 合并两个有序数组 (Merge Sorted Array)
//
// 题目描述:
// 给你两个按 "非递减顺序" 排列的整数数组 nums1 和 nums2，另有两个整数 m 和 n ，分别表示 nums1 和 nums2 中的元素数目。
// 请你 合并 nums2 到 nums1 中，使合并后的数组同样按 "非递减顺序" 排列。
//
// 注意：
// 最终，合并后数组不应由函数返回，而是存储在数组 nums1 中。
// 为了应对这种情况，nums1 的初始长度为 m + n，其中前 m 个元素表示应合并的元素，后 n 个元素为 0 ，应忽略。
// nums2 的长度为 n 。
//
// 示例 1：
// 输入：nums1 = [1,2,3,0,0,0], m = 3, nums2 = [2,5,6], n = 3
// 输出：[1,2,2,3,5,6]
//
// 示例 2：
// 输入：nums1 = [1], m = 1, nums2 = [], n = 0
// 输出：[1]

func merge(nums1 []int, m int, nums2 []int, n int) {
	idx := m + n - 1
	for n > 0 {
		if m > 0 && nums1[m-1] > nums2[n-1] {
			nums1[idx] = nums1[m-1]
			m--
		} else {
			nums1[idx] = nums2[n-1]
			n--
		}
		idx--
	}
}

func TestMerge(t *testing.T) {
	tests := []struct {
		name     string
		nums1    []int
		m        int
		nums2    []int
		n        int
		expected []int
	}{
		{
			name:     "Example 1",
			nums1:    []int{1, 2, 3, 0, 0, 0},
			m:        3,
			nums2:    []int{2, 5, 6},
			n:        3,
			expected: []int{1, 2, 2, 3, 5, 6},
		},
		{
			name:     "Example 2",
			nums1:    []int{1},
			m:        1,
			nums2:    []int{},
			n:        0,
			expected: []int{1},
		},
		{
			name:     "Example 3",
			nums1:    []int{0},
			m:        0,
			nums2:    []int{1},
			n:        1,
			expected: []int{1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: Uncomment when implementation is ready
			nums1 := make([]int, len(tt.nums1))
			copy(nums1, tt.nums1)

			merge(nums1, tt.m, tt.nums2, tt.n)

			if !reflect.DeepEqual(nums1, tt.expected) {
				t.Errorf("merge() resulted in %v, want %v", nums1, tt.expected)
			}
		})
	}
}

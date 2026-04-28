package top100liked

import (
	"math"
	"testing"
)

// 4. 寻找两个正序数组的中位数 (Median of Two Sorted Arrays)
//
// 题目描述:
// 给定两个大小分别为 m 和 n 的正序（从小到大）数组 nums1 和 nums2。
// 请你找出并返回这两个正序数组的中位数。算法的时间复杂度应该为 O(log(m+n))。
//
// 示例 1：
// 输入：nums1 = [1,3], nums2 = [2]
// 输出：2.00000
//
// 示例 2：
// 输入：nums1 = [1,2], nums2 = [3,4]
// 输出：2.50000

// findMedianSortedArrays 在较短数组上二分，O(log(min(m,n))) 时间，O(1) 空间
// 核心思路：在两个数组上各切一刀，使左半部分元素数 = 右半部分，且左边最大值 <= 右边最小值
func findMedianSortedArrays(nums1 []int, nums2 []int) float64 {
	// 保证 nums1 是较短的数组，缩小二分范围
	if len(nums1) > len(nums2) {
		nums1, nums2 = nums2, nums1
	}
	m, n := len(nums1), len(nums2)
	half := (m + n + 1) / 2 // 左半部分应有的元素总数

	left, right := 0, m
	for left <= right {
		cut1 := (left + right) / 2 // nums1 的切割位置
		cut2 := half - cut1        // nums2 的切割位置

		// 切割线左右的四个值，越界时用 ±∞
		l1, r1 := math.MinInt64, math.MaxInt64
		if cut1 > 0 {
			l1 = nums1[cut1-1]
		}
		if cut1 < m {
			r1 = nums1[cut1]
		}
		l2, r2 := math.MinInt64, math.MaxInt64
		if cut2 > 0 {
			l2 = nums2[cut2-1]
		}
		if cut2 < n {
			r2 = nums2[cut2]
		}

		if l1 <= r2 && l2 <= r1 {
			// 找到正确切割位置
			if (m+n)%2 == 1 {
				return float64(max(l1, l2))
			}
			return float64(max(l1, l2)+min(r1, r2)) / 2.0
		} else if l1 > r2 {
			right = cut1 - 1 // nums1 左边太大，左移
		} else {
			left = cut1 + 1 // nums1 左边太小，右移
		}
	}
	return 0
}

func TestFindMedianSortedArrays(t *testing.T) {
	tests := []struct {
		name     string
		nums1    []int
		nums2    []int
		expected float64
	}{
		{name: "示例1", nums1: []int{1, 3}, nums2: []int{2}, expected: 2.0},
		{name: "示例2", nums1: []int{1, 2}, nums2: []int{3, 4}, expected: 2.5},
		{name: "单元素", nums1: []int{}, nums2: []int{1}, expected: 1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findMedianSortedArrays(tt.nums1, tt.nums2)
			if got != tt.expected {
				t.Errorf("findMedianSortedArrays() = %v, want %v", got, tt.expected)
			}
		})
	}
}

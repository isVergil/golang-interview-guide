package topinterview150

import (
	"testing"
)

// 4. 寻找两个正序数组的中位数 (Median of Two Sorted Arrays)
//
// 题目描述:
// 给定两个大小分别为 m 和 n 的正序（从小到大）数组 nums1 和 nums2。请你找出并返回这两个正序数组的 中位数 。
// 算法的时间复杂度应该为 O(log (m+n)) 。
//
// 示例 1：
// 输入：nums1 = [1,3], nums2 = [2]
// 输出：2.00000
// 解释：合并数组 = [1,2,3] ，中位数 2
//
// 示例 2：
// 输入：nums1 = [1,2], nums2 = [3,4]
// 输出：2.50000
// 解释：合并数组 = [1,2,3,4] ，中位数 (2 + 3) / 2 = 2.5

func findMedianSortedArrays(nums1 []int, nums2 []int) float64 {
	m, n := len(nums1), len(nums2)
	totalLen := m + n

	// 前一个和当前的数
	prev, cur := 0, 0
	p1, p2 := 0, 0

	for i := 0; i <= totalLen/2; i++ {
		prev = cur
		if p1 < m && (p2 >= n || nums1[p1] < nums2[p2]) {
			cur = nums1[p1]
			p1++
		} else {
			cur = nums2[p2]
			p2++
		}
	}

	if totalLen%2 == 0 {
		return float64(cur+prev) / 2.0
	}

	return float64(cur)

}

func TestFindMedianSortedArrays(t *testing.T) {
	tests := []struct {
		name     string
		nums1    []int
		nums2    []int
		expected float64
	}{
		{"Example 1", []int{1, 3}, []int{2}, 2.0},
		{"Example 2", []int{1, 2}, []int{3, 4}, 2.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := findMedianSortedArrays(tt.nums1, tt.nums2); got != tt.expected {
			// 	t.Errorf("findMedianSortedArrays() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

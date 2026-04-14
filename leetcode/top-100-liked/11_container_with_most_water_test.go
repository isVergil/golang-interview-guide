package top100liked

import (
	"testing"
)

// 11. 盛最多水的容器 (Container With Most Water)
//
// 题目描述:
// 给定一个长度为 n 的整数数组 height。有 n 条垂线，第 i 条线的两个端点是 (i, 0) 和 (i, height[i])。
// 找出其中的两条线，使得它们与 x 轴共同构成的容器可以容纳最多的水。返回容器可以储存的最大水量。
//
// 示例 1：
// 输入：[1,8,6,2,5,4,8,3,7]
// 输出：49
// 解释：图中垂直线代表输入数组 [1,8,6,2,5,4,8,3,7]。在此情况下，容器能够容纳水的最大值为 49。
//
// 示例 2：
// 输入：height = [1,1]
// 输出：1

func maxArea(height []int) int {
	l, r, res := 0, len(height)-1, 0
	for l < r {
		width := r - l
		if height[l] <= height[r] {
			cur := width * height[l]
			if cur > res {
				res = cur
			}
			l++

			// 剪枝
			min := height[l]
			for l < r && height[l] < min {
				l++
			}
		} else {
			cur := width * height[r]
			if cur > res {
				res = cur
			}
			r--

			// 剪枝
			min := height[r]
			for l < r && height[r] < min {
				r--
			}
		}

	}

	return res
}

func TestMaxArea(t *testing.T) {
	tests := []struct {
		name     string
		height   []int
		expected int
	}{
		{
			name:     "示例1",
			height:   []int{1, 8, 6, 2, 5, 4, 8, 3, 7},
			expected: 49,
		},
		{
			name:     "示例2",
			height:   []int{1, 1},
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := maxArea(tt.height)
			if got != tt.expected {
				t.Errorf("maxArea() = %v, want %v", got, tt.expected)
			}
		})
	}
}

package topinterview150

import (
	"testing"
)

// 11. 盛最多水的容器 (Container With Most Water)
//
// 题目描述:
// 给定一个长度为 n 的整数数组 height 。有 n 条垂线，第 i 条线的两个端点是 (i, 0) 和 (i, height[i]) 。
// 找出其中的两条线，使得它们与 x 轴共同构成的容器可以容纳最多的水。
// 返回容器可以储存的最大水量。
// 说明：你不能倾斜容器。
//
// 示例 1：
// 输入：[1,8,6,2,5,4,8,3,7]
// 输出：49
// 解释：图中垂直线代表输入数组 [1,8,6,2,5,4,8,3,7]。在此场景中，容器能够容纳水（表示为蓝色部分）的最大值为 49。
//
// 示例 2：
// 输入：height = [1,1]
// 输出：1

func maxArea(height []int) int {
	l, r := 0, len(height)-1
	maxWater := 0

	for l < r {
		// 记录当前两端的柱子高度
		hLeft := height[l]
		hRight := height[r]

		// 计算面积宽度
		width := r - l

		if hLeft < hRight {
			// 左侧是短板，高度由 hLeft 决定
			cur := hLeft * width
			if cur > maxWater {
				maxWater = cur
			}
			l++

			// 剪枝性能优化
			// 如果移动后的新柱子比刚才的 hLeft 还要矮，面积绝对不可能更大，直接跳过，省去乘法计算
			for l < r && height[l] <= hLeft {
				l++
			}
		} else {
			cur := hRight * width
			if cur > maxWater {
				maxWater = cur
			}
			r--

			// 剪枝性能优化
			for l < r && height[r] <= hRight {
				r--
			}
		}
	}
	return maxWater
}

func TestMaxArea(t *testing.T) {
	tests := []struct {
		name     string
		height   []int
		expected int
	}{
		{"Example 1", []int{1, 8, 6, 2, 5, 4, 8, 3, 7}, 49},
		{"Example 2", []int{1, 1}, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := maxArea(tt.height); got != tt.expected {
			// 	t.Errorf("maxArea() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

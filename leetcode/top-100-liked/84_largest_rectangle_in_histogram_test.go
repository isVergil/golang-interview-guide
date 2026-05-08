package top100liked

import (
	"testing"
)

// 84. 柱状图中最大的矩形 (Largest Rectangle in Histogram)
//
// 题目描述:
// 给定 n 个非负整数，用来表示柱状图中各个柱子的高度。每个柱子彼此相邻，且宽度为 1。
// 求在该柱状图中，能够勾勒出来的矩形的最大面积。
//
// 示例 1：
// 输入：heights = [2,1,5,6,2,3]
// 输出：10（5和6两根柱子组成的矩形，宽度2，高度5）
//
// 示例 2：
// 输入：heights = [2,4]
// 输出：4

func largestRectangleArea(heights []int) int {
	n := len(heights)
	stack := make([]int, 0, n)
	maxArea := 0
	for i := 0; i <= n; i++ {
		curHeight := 0
		if i < n {
			curHeight = heights[i]
		}

		// 当前柱子比栈顶矮
		for len(stack) > 0 && curHeight < heights[stack[len(stack)-1]] {
			// 弹出栈顶 计算面积
			topIdx := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			h := heights[topIdx]

			// 计算宽度
			width := 0
			if len(stack) == 0 {
				width = i
			} else {
				width = i - stack[len(stack)-1] - 1
			}
			if h*width > maxArea {
				maxArea = h * width
			}
		}
		stack = append(stack, i)
	}
	return maxArea
}

func TestLargestRectangleArea(t *testing.T) {
	tests := []struct {
		name     string
		heights  []int
		expected int
	}{
		{name: "示例1", heights: []int{2, 1, 5, 6, 2, 3}, expected: 10},
		{name: "示例2", heights: []int{2, 4}, expected: 4},
		{name: "单柱", heights: []int{1}, expected: 1},
		{name: "递增", heights: []int{1, 2, 3, 4, 5}, expected: 9},
		{name: "等高", heights: []int{3, 3, 3}, expected: 9},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := largestRectangleArea(tt.heights)
			if got != tt.expected {
				t.Errorf("largestRectangleArea() = %v, want %v", got, tt.expected)
			}
		})
	}
}

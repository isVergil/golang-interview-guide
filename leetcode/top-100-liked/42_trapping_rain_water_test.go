package top100liked

import (
	"testing"
)

// 42. 接雨水 (Trapping Rain Water)
//
// 题目描述:
// 给定 n 个非负整数表示每个宽度为 1 的柱子的高度图，计算按此排列的柱子，下雨之后能接多少雨水。
//
// 示例 1：
// 输入：height = [0,1,0,2,1,0,1,3,2,1,2,1]
// 输出：6
//
// 示例 2：
// 输入：height = [4,2,0,3,2,5]
// 输出：9

// 哪边矮，哪边就是瓶颈，处理哪边。
func trap(height []int) int {
	l, r := 0, len(height)-1
	lMax, rMax := 0, 0
	res := 0
	for l < r {
		if height[l] < height[r] {
			// 右边有比当前更高的柱子兜底，瓶颈在左边
			if height[l] > lMax {
				lMax = height[l] // 更新左边最高
			} else {
				res += lMax - height[l]
			}
			l++
		} else {
			// 左边有比当前更高的柱子兜底，瓶颈在右边
			if height[r] > rMax {
				rMax = height[r]
			} else {
				res += rMax - height[r]
			}
			r--
		}
	}
	return res
}

func TestTrap(t *testing.T) {
	tests := []struct {
		name     string
		height   []int
		expected int
	}{
		{name: "示例1", height: []int{0, 1, 0, 2, 1, 0, 1, 3, 2, 1, 2, 1}, expected: 6},
		{name: "示例2", height: []int{4, 2, 0, 3, 2, 5}, expected: 9},
		{name: "空数组", height: []int{}, expected: 0},
		{name: "无法接水", height: []int{1, 2, 3}, expected: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := trap(tt.height)
			if got != tt.expected {
				t.Errorf("trap() = %v, want %v", got, tt.expected)
			}
		})
	}
}

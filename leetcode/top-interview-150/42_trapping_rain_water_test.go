package topinterview150

import (
	"testing"
)

// 42. 接雨水 (Trapping Rain Water)
//
// 题目描述:
// 给定 n 个非负整数表示每个宽度为 1 的柱子的高度图，计算按此排列的柱子，下雨之后能接多少雨水。
//
// 示例 1:
// 输入：height = [0,1,0,2,1,0,1,3,2,1,2,1]
// 输出：6
// 解释：上面是由数组 [0,1,0,2,1,0,1,3,2,1,2,1] 表示的高度图，在这种情况下，可以接 6 个单位的雨水（蓝色部分表示雨水）。
//
// 示例 2:
// 输入：height = [4,2,0,3,2,5]
// 输出：9

func trap(height []int) int {
	if len(height) < 3 {
		return 0
	}

	l, r := 0, len(height)-1
	lMax, rMax := 0, 0
	res := 0
	for l < r {
		if height[l] > lMax {
			lMax = height[l]
		}
		if height[r] > rMax {
			rMax = height[r]
		}

		// 哪边短就结算哪边
		if lMax < rMax {
			res += lMax - height[l]
			l++
		} else {
			res += rMax - height[r]
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
		{"Example 1", []int{0, 1, 0, 2, 1, 0, 1, 3, 2, 1, 2, 1}, 6},
		{"Example 2", []int{4, 2, 0, 3, 2, 5}, 9},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := trap(tt.height); got != tt.expected {
			// 	t.Errorf("trap() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

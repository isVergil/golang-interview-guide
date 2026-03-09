package topinterview150

import (
	"testing"
)

// 149. 直线上最多的点数 (Max Points on a Line)
//
// 题目描述:
// 给你一个数组 points ，其中 points[i] = [xi, yi] 表示 X-Y 平面上的一个点。求最多有多少个点在同一条直线上。
//
// 示例 1：
// 输入：points = [[1,1],[2,2],[3,3]]
// 输出：3
//
// 示例 2：
// 输入：points = [[1,1],[3,2],[5,3],[4,1],[2,3],[1,4]]
// 输出：4

func maxPoints(points [][]int) int {
	panic("not implemented")
}

func TestMaxPoints(t *testing.T) {
	tests := []struct {
		name     string
		points   [][]int
		expected int
	}{
		{"Example 1", [][]int{{1, 1}, {2, 2}, {3, 3}}, 3},
		{"Example 2", [][]int{{1, 1}, {3, 2}, {5, 3}, {4, 1}, {2, 3}, {1, 4}}, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := maxPoints(tt.points); got != tt.expected {
			// 	t.Errorf("maxPoints() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

package topinterview150

import (
	"sort"
	"testing"
)

// 452. 用最少数量的箭引爆气球 (Minimum Number of Arrows to Burst Balloons)
//
// 题目描述:
// 有一些球形气球贴在一堵用 XY 平面表示的墙面上。墙面上有许多球形气球。对于每个气球，提供的输入是水平方向上，气球直径的开始和结束坐标。由于是球形，所以y坐标并不重要。
// 气球的宽度可以覆盖 xstart 到 xend 的范围。
// 弓箭可以沿着 x 轴从不同点完全垂直地射出。在坐标 x 处射出一支箭，若有一个气球的直径的开始和结束坐标为 xstart，xend， 且满足  xstart ≤ x ≤ xend，则该气球会被引爆。
// 可以射出的弓箭的数量没有限制。 弓箭一旦被射出之后，可以无限地前进。
// 给你一个数组 points ，其中 points [i] = [xstart, xend] ，返回引爆所有气球所必须射出的最小弓箭数。
//
// 示例 1：
// 输入：points = [[10,16],[2,8],[1,6],[7,12]]
// 输出：2
// 解释：气球可以用2支箭来爆破:
// -在x = 6处射出箭，击破气球[2,8]和[1,6]。
// -在x = 11处射出箭，击破气球[10,16]和[7,12]。
//
// 示例 2：
// 输入：points = [[1,2],[3,4],[5,6],[7,8]]
// 输出：4
//
// 示例 3：
// 输入：points = [[1,2],[2,3],[3,4],[4,5]]
// 输出：2

func findMinArrowShots(points [][]int) int {
	sort.Slice(points, func(i, j int) bool {
		return points[i][1] < points[j][1]
	})

	last, cnt := points[0][1], 1
	for i := 1; i < len(points); i++ {
		if last < points[i][0] {
			last = points[i][1]
			cnt++
		}
	}
	return cnt
}

func TestFindMinArrowShots(t *testing.T) {
	tests := []struct {
		name     string
		points   [][]int
		expected int
	}{
		{
			name:     "Example 1",
			points:   [][]int{{10, 16}, {2, 8}, {1, 6}, {7, 12}},
			expected: 2,
		},
		{
			name:     "Example 2",
			points:   [][]int{{1, 2}, {3, 4}, {5, 6}, {7, 8}},
			expected: 4,
		},
		{
			name:     "Example 3",
			points:   [][]int{{1, 2}, {2, 3}, {3, 4}, {4, 5}},
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: Uncomment when implementation is ready
			// if got := findMinArrowShots(tt.points); got != tt.expected {
			// 	t.Errorf("findMinArrowShots() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

package topinterview150

import (
	"sort"
	"testing"
)

// 56. 合并区间 (Merge Intervals)
//
// 题目描述:
// 以数组 intervals 表示若干个区间的集合，其中单个区间为 intervals[i] = [starti, endi] 。
// 请你合并所有重叠的区间，并返回 一个不重叠的区间数组，该数组需恰好覆盖输入中的所有区间 。
//
// 示例 1：
// 输入：intervals = [[1,3],[2,6],[8,10],[15,18]]
// 输出：[[1,6],[8,10],[15,18]]
// 解释：区间 [1,3] 和 [2,6] 重叠, 将它们合并为 [1,6].
//
// 示例 2：
// 输入：intervals = [[1,4],[4,5]]
// 输出：[[1,5]]
// 解释：区间 [1,4] 和 [4,5] 可被视为重叠区间。

func mergeIntervals(intervals [][]int) [][]int {
	// < 2无需排序
	if len(intervals) < 2 {
		return intervals
	}

	// 按左端点排序
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})

	// 初始化结果集，先放入第一个区间
	merged := [][]int{}
	merged = append(merged, intervals[0])

	// 遍历
	for i := 1; i < len(intervals); i++ {
		last := merged[len(merged)-1]
		// 有重叠就合并
		if last[1] >= intervals[i][0] {
			if last[1] < intervals[i][1] {
				last[1] = intervals[i][1]
			}
		} else {
			merged = append(merged, intervals[i])
		}
	}
	return merged
}

func TestMergeIntervals(t *testing.T) {
	tests := []struct {
		name      string
		intervals [][]int
		expected  [][]int
	}{
		{
			name:      "Example 1",
			intervals: [][]int{{1, 3}, {2, 6}, {8, 10}, {15, 18}},
			expected:  [][]int{{1, 6}, {8, 10}, {15, 18}},
		},
		{
			name:      "Example 2",
			intervals: [][]int{{1, 4}, {4, 5}},
			expected:  [][]int{{1, 5}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: Uncomment when implementation is ready
			// if got := merge(tt.intervals); !reflect.DeepEqual(got, tt.expected) {
			// 	t.Errorf("merge() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

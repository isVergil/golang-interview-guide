package topinterview150

import (
	"sort"
	"testing"
)

// 57. 插入区间 (Insert Interval)
//
// 题目描述:
// 给你一个 无重叠的 ，按照区间起始端点排序的区间列表。
// 在列表中插入一个新的区间，你需要确保列表中的区间仍然有序且不重叠（如果有必要的话，可以合并区间）。
//
// 示例 1：
// 输入：intervals = [[1,3],[6,9]], newInterval = [2,5]
// 输出：[[1,5],[6,9]]
//
// 示例 2：
// 输入：intervals = [[1,2],[3,5],[6,7],[8,10],[12,16]], newInterval = [4,8]
// 输出：[[1,2],[3,10],[12,16]]
// 解释：这是因为新的区间 [4,8] 与 [3,5],[6,7],[8,10] 重叠。

// 1 newInterval 加入 intervals 形成新数组，再按合并排序做
func insert(intervals [][]int, newInterval []int) [][]int {
	// 先加入 newInterval 到intervals 问题就变成了合并区间了
	intervals = append(intervals, newInterval)

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

// 2 本身就是无重叠的 ，按照区间起始端点排序的区间列表，遍历取得 left right
func insert1(intervals [][]int, newInterval []int) [][]int {
	merged := [][]int{}
	i := 0
	n := len(intervals)

	//  1 处理左侧无重叠部分
	// 只要当前区间的右端点 < 新区间的左端点，就直接加入
	for i < n && intervals[i][1] < newInterval[0] {
		merged = append(merged, intervals[i])
		i++
	}

	// 2 处理重叠
	// 重叠条件：当前区间的左端点 <= 新区间的右端点
	for i < n && intervals[i][0] <= newInterval[1] {
		// 更新新区间的左边界为两者中的最小值
		newInterval[0] = min(newInterval[0], intervals[i][0])
		// 更新新区间的右边界为两者中的最大值
		newInterval[1] = max(newInterval[1], intervals[i][1])
		i++
	}
	// 将合并后的“巨型新区间”放入结果
	merged = append(merged, newInterval)

	// 3 处理右侧无重叠部分
	for i < n {
		merged = append(merged, intervals[i])
		i++
	}

	return merged
}

func TestInsertInterval(t *testing.T) {
	tests := []struct {
		name        string
		intervals   [][]int
		newInterval []int
		expected    [][]int
	}{
		{
			name:        "Example 1",
			intervals:   [][]int{{1, 3}, {6, 9}},
			newInterval: []int{2, 5},
			expected:    [][]int{{1, 5}, {6, 9}},
		},
		{
			name:        "Example 2",
			intervals:   [][]int{{1, 2}, {3, 5}, {6, 7}, {8, 10}, {12, 16}},
			newInterval: []int{4, 8},
			expected:    [][]int{{1, 2}, {3, 10}, {12, 16}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: Uncomment when implementation is ready
			// if got := insert(tt.intervals, tt.newInterval); !reflect.DeepEqual(got, tt.expected) {
			// 	t.Errorf("insert() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

package topinterview150

import (
	"testing"
)

// 207. 课程表 (Course Schedule)
//
// 题目描述:
// 你这个学期必须选修 numCourses 门课程，记为 0 到 numCourses - 1 。
// 在选修某些课程之前需要一些先修课程。 先修课程按数组 prerequisites 给出，其中 prerequisites[i] = [ai, bi] ，表示如果要学习课程 ai 则 必须 先学习课程  bi 。
// 例如，先修课程对 [0, 1] 表示：想要学习课程 0 ，你需要先完成课程 1 。
// 请你判断是否可能完成所有课程的学习？如果可以，返回 true ；否则，返回 false 。
//
// 示例 1：
// 输入：numCourses = 2, prerequisites = [[1,0]]
// 输出：true
// 解释：总共有 2 门课程。学习课程 1 之前，你需要完成课程 0 。这是可能的。
//
// 示例 2：
// 输入：numCourses = 2, prerequisites = [[1,0],[0,1]]
// 输出：false
// 解释：总共有 2 门课程。学习课程 1 之前，你需要先完成​课程 0 ；​并且学习课程 0 之前，你还应先完成课程 1 。这是不可能的。

func canFinish(numCourses int, prerequisites [][]int) bool {
	panic("not implemented")
}

func TestCanFinish(t *testing.T) {
	tests := []struct {
		name          string
		numCourses    int
		prerequisites [][]int
		expected      bool
	}{
		{"Example 1", 2, [][]int{{1, 0}}, true},
		{"Example 2", 2, [][]int{{1, 0}, {0, 1}}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := canFinish(tt.numCourses, tt.prerequisites); got != tt.expected {
			// 	t.Errorf("canFinish() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

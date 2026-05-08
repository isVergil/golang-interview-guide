package top100liked

import (
	"testing"
)

// 207. 课程表 (Course Schedule)
//
// 题目描述:
// 你这个学期必须选修 numCourses 门课程，记为 0 到 numCourses-1。
// 在选修某些课程之前需要一些先修课程，prerequisites[i] = [ai, bi] 表示如果要学习课程 ai
// 则必须先学习课程 bi。请你判断是否可能完成所有课程的学习？
//
// 示例 1：
// 输入：numCourses = 2, prerequisites = [[1,0]]
// 输出：true（先修课程0，再修课程1，可以完成）
//
// 示例 2：
// 输入：numCourses = 2, prerequisites = [[1,0],[0,1]]
// 输出：false（课程0和1互相依赖，无法完成）

func canFinish(numCourses int, prerequisites [][]int) bool {
	// 建图
	graph := make([][]int, numCourses)  // graph[i] 后续的课程列表
	inDegree := make([]int, numCourses) // inDegree[i] 的先修课程

	for _, course := range prerequisites {
		cur, pre := course[0], course[1]
		graph[pre] = append(graph[pre], cur)
		inDegree[cur]++
	}

	// 入度为 0 （先修的）的入队列
	queue := make([]int, 0)
	for i, v := range inDegree {
		if v == 0 {
			queue = append(queue, i)
		}
	}

	// BFS
	count := 0
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		count++
		for _, next := range graph[cur] {
			inDegree[next]--
			if inDegree[next] == 0 {
				queue = append(queue, next)
			}
		}
	}

	return count == numCourses

}

func TestCanFinish(t *testing.T) {
	tests := []struct {
		name          string
		numCourses    int
		prerequisites [][]int
		expected      bool
	}{
		{name: "示例1", numCourses: 2, prerequisites: [][]int{{1, 0}}, expected: true},
		{name: "示例2-循环依赖", numCourses: 2, prerequisites: [][]int{{1, 0}, {0, 1}}, expected: false},
		{name: "无先修", numCourses: 1, prerequisites: [][]int{}, expected: true},
		{name: "线性依赖", numCourses: 3, prerequisites: [][]int{{1, 0}, {2, 1}}, expected: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := canFinish(tt.numCourses, tt.prerequisites)
			if got != tt.expected {
				t.Errorf("canFinish() = %v, want %v", got, tt.expected)
			}
		})
	}
}

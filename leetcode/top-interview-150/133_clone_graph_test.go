package topinterview150

import (
	"testing"
)

// 133. 克隆图 (Clone Graph)
//
// 题目描述:
// 给你无向 连通 图中一个节点的引用，请你返回该图的 深拷贝（克隆）。
// 图中的每个节点都包含它的值 val（int） 和其邻居的列表（list[Node]）。
//
// 示例 1：
// 输入：adjList = [[2,4],[1,3],[2,4],[1,3]]
// 输出：[[2,4],[1,3],[2,4],[1,3]]

// Node133 defines a graph node.
type Node133 struct {
	Val       int
	Neighbors []*Node133
}

func cloneGraph(node *Node133) *Node133 {
	panic("not implemented")
}

func TestCloneGraph(t *testing.T) {
	// 克隆图测试
}

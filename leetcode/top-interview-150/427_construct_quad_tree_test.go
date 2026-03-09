package topinterview150

import (
	"testing"
)

// 427. 建立四叉树 (Construct Quad Tree)
//
// 题目描述:
// 给你一个 n * n 矩阵 grid ，矩阵中每个单元格的值不是 0 就是 1 。
// 请你返回一个表示矩阵的 四叉树 的根结点。
//
// Node427 defines a Quad Tree node.
type Node427 struct {
	Val         bool
	IsLeaf      bool
	TopLeft     *Node427
	TopRight    *Node427
	BottomLeft  *Node427
	BottomRight *Node427
}

func constructQuadTree(grid [][]int) *Node427 {
	panic("not implemented")
}

func TestConstructQuadTree(t *testing.T) {
	// 四叉树构建测试
}

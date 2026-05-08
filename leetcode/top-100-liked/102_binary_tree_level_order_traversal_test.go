package top100liked

import (
	"reflect"
	"testing"
)

// 102. 二叉树的层序遍历 (Binary Tree Level Order Traversal)
//
// 题目描述:
// 给你二叉树的根节点 root，返回其节点值的层序遍历（即逐层地，从左到右访问所有节点）。
//
// 示例 1：
// 输入：root = [3,9,20,null,null,15,7]
// 输出：[[3],[9,20],[15,7]]
//
// 示例 2：
// 输入：root = [1]
// 输出：[[1]]

func levelOrder(root *TreeNode) [][]int {
	queue := make([]*TreeNode, 0)
	if root != nil {
		queue = append(queue, root)
	}
	res := make([][]int, 0)
	for len(queue) > 0 {
		size := len(queue)
		curLevelVal := make([]int, 0, size)
		for i := 0; i < size; i++ {
			curNode := queue[i]
			curLevelVal = append(curLevelVal, curNode.Val)
			if curNode.Left != nil {
				queue = append(queue, curNode.Left)
			}
			if curNode.Right != nil {
				queue = append(queue, curNode.Right)
			}
		}
		res = append(res, curLevelVal)
		queue = queue[size:]
	}
	return res
}

func TestLevelOrder(t *testing.T) {
	tests := []struct {
		name     string
		root     *TreeNode
		expected [][]int
	}{
		{name: "示例1", root: &TreeNode{Val: 3, Left: &TreeNode{Val: 9}, Right: &TreeNode{Val: 20, Left: &TreeNode{Val: 15}, Right: &TreeNode{Val: 7}}}, expected: [][]int{{3}, {9, 20}, {15, 7}}},
		{name: "单节点", root: &TreeNode{Val: 1}, expected: [][]int{{1}}},
		{name: "空树", root: nil, expected: [][]int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := levelOrder(tt.root)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("levelOrder() = %v, want %v", got, tt.expected)
			}
		})
	}
}

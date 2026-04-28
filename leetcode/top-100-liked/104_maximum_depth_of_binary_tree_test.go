package top100liked

import (
	"testing"
)

// 104. 二叉树的最大深度 (Maximum Depth of Binary Tree)
//
// 题目描述:
// 给定一个二叉树 root ，返回其最大深度。
// 二叉树的最大深度是指从根节点到最远叶子节点的最长路径上的节点数。
//
// 示例 1：
// 输入：root = [3,9,20,null,null,15,7]
// 输出：3
//
// 示例 2：
// 输入：root = [1,null,2]
// 输出：2
// 递归
func maxDepth(root *TreeNode) int {
	if root == nil {
		return 0
	}
	left := maxDepth(root.Left)
	right := maxDepth(root.Right)
	if left > right {
		return left + 1
	}
	return right + 1
}

// 非递归
func maxDepth1(root *TreeNode) int {
	if root == nil {
		return 0
	}
	depth := 0
	queue := []*TreeNode{root}
	for len(queue) > 0 {
		depth++
		size := len(queue)
		for i := 0; i < size; i++ {
			node := queue[0]
			queue = queue[1:]
			if node.Left != nil {
				queue = append(queue, node.Left)
			}
			if node.Right != nil {
				queue = append(queue, node.Right)
			}
		}
	}
	return depth
}

func TestMaxDepth(t *testing.T) {
	tests := []struct {
		name     string
		root     *TreeNode
		expected int
	}{
		{
			name: "示例1",
			root: &TreeNode{
				Val:  3,
				Left: &TreeNode{Val: 9},
				Right: &TreeNode{
					Val:   20,
					Left:  &TreeNode{Val: 15},
					Right: &TreeNode{Val: 7},
				},
			},
			expected: 3,
		},
		{
			name: "示例2",
			root: &TreeNode{
				Val:   1,
				Right: &TreeNode{Val: 2},
			},
			expected: 2,
		},
		{
			name:     "空树",
			root:     nil,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := maxDepth(tt.root)
			if got != tt.expected {
				t.Errorf("maxDepth() = %v, want %v", got, tt.expected)
			}
		})

		t.Run(tt.name, func(t *testing.T) {
			got := maxDepth1(tt.root)
			if got != tt.expected {
				t.Errorf("maxDepth1() = %v, want %v", got, tt.expected)
			}
		})
	}
}

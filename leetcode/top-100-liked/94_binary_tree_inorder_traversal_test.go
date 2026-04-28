package top100liked

import (
	"reflect"
	"testing"
)

// 94. 二叉树的中序遍历 (Binary Tree Inorder Traversal)
//
// 题目描述:
// 给定一个二叉树的根节点 root ，返回它的中序遍历。
//
// 示例 1：
// 输入：root = [1,null,2,3]
// 输出：[1,3,2]
//
// 示例 2：
// 输入：root = []
// 输出：[]
//
// 示例 3：
// 输入：root = [1]
// 输出：[1]

func inorderTraversal(root *TreeNode) []int {
	var res []int
	var dfs func(*TreeNode)
	dfs = func(node *TreeNode) {
		if node == nil {
			return
		}
		dfs(node.Left)
		res = append(res, node.Val)
		dfs(node.Right)
	}
	dfs(root)
	return res
}

func inorderTraversal1(root *TreeNode) []int {
	var res []int
	var stack []*TreeNode
	cur := root
	for cur != nil || len(stack) > 0 {
		for cur != nil {
			stack = append(stack, cur.Left)
			cur = cur.Left
		}
		cur = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		res = append(res, cur.Val)
		cur = cur.Right
	}
	return res
}

func TestInorderTraversal(t *testing.T) {
	tests := []struct {
		name     string
		root     *TreeNode
		expected []int
	}{
		{
			name:     "示例1",
			root:     &TreeNode{Val: 1, Right: &TreeNode{Val: 2, Left: &TreeNode{Val: 3}}},
			expected: []int{1, 3, 2},
		},
		{
			name:     "空树",
			root:     nil,
			expected: nil,
		},
		{
			name:     "单节点",
			root:     &TreeNode{Val: 1},
			expected: []int{1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := inorderTraversal(tt.root)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("inorderTraversal() = %v, want %v", got, tt.expected)
			}
		})

		t.Run(tt.name, func(t *testing.T) {
			got := inorderTraversal1(tt.root)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("inorderTraversal1() = %v, want %v", got, tt.expected)
			}
		})
	}
}

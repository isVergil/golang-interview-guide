package top100liked

import (
	"reflect"
	"testing"
)

// 105. 从前序与中序遍历序列构造二叉树 (Construct Binary Tree from Preorder and Inorder Traversal)
//
// 题目描述:
// 给定两个整数数组 preorder 和 inorder，其中 preorder 是二叉树的先序遍历，
// inorder 是同一棵树的中序遍历，请构造二叉树并返回其根节点。
//
// 示例 1：
// 输入：preorder = [3,9,20,15,7], inorder = [9,3,15,20,7]
// 输出：[3,9,20,null,null,15,7]
//
// 示例 2：
// 输入：preorder = [-1], inorder = [-1]
// 输出：[-1]
//
// 提示：前序第一个是根，在中序中找到根的位置，左边是左子树，右边是右子树，递归构建

func buildTree(preorder []int, inorder []int) *TreeNode {
	inorderIdx := make(map[int]int, len(inorder))
	for i, v := range inorder {
		inorderIdx[v] = i
	}

	var build func(preLeft, preRight, inLeft, inRight int) *TreeNode
	build = func(preLeft, preRight, inLeft, inRight int) *TreeNode {
		if preLeft > preRight {
			return nil
		}

		// 前序第一个就是根
		rootVal := preorder[preLeft]
		root := &TreeNode{Val: rootVal}

		// 在中序中找到根的位置
		rootIdx := inorderIdx[rootVal]

		// 左子树的节点个数
		leftSize := rootIdx - inLeft

		// 拆分左右子树
		root.Left = build(preLeft+1, preLeft+leftSize, inLeft, rootIdx-1)
		root.Right = build(preLeft+leftSize+1, preRight, rootIdx+1, inRight)
		return root
	}

	return build(0, len(preorder)-1, 0, len(inorder)-1)
}

func TestBuildTree(t *testing.T) {
	// 辅助：层序遍历用于验证
	levelOrder := func(root *TreeNode) []int {
		if root == nil {
			return nil
		}
		res := []int{}
		queue := []*TreeNode{root}
		for len(queue) > 0 {
			node := queue[0]
			queue = queue[1:]
			res = append(res, node.Val)
			if node.Left != nil {
				queue = append(queue, node.Left)
			}
			if node.Right != nil {
				queue = append(queue, node.Right)
			}
		}
		return res
	}

	tests := []struct {
		name     string
		preorder []int
		inorder  []int
		expected []int
	}{
		{
			name:     "示例1",
			preorder: []int{3, 9, 20, 15, 7},
			inorder:  []int{9, 3, 15, 20, 7},
			expected: []int{3, 9, 20, 15, 7},
		},
		{
			name:     "单节点",
			preorder: []int{-1},
			inorder:  []int{-1},
			expected: []int{-1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := buildTree(tt.preorder, tt.inorder)
			got := levelOrder(root)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("buildTree() level order = %v, want %v", got, tt.expected)
			}
		})
	}
}

package top100liked

import (
	"testing"
)

// 236. 二叉树的最近公共祖先 (Lowest Common Ancestor of a Binary Tree)
//
// 题目描述:
// 给定一个二叉树, 找到该树中两个指定节点的最近公共祖先。
// 最近公共祖先的定义为：对于有根树 T 的两个节点 p、q，最近公共祖先表示为一个节点 x，
// 满足 x 是 p、q 的祖先且 x 的深度尽可能大（一个节点也可以是它自己的祖先）。
//
// 示例 1：
// 输入：root = [3,5,1,6,2,0,8,null,null,7,4], p = 5, q = 1
// 输出：3
//
// 示例 2：
// 输入：root = [3,5,1,6,2,0,8,null,null,7,4], p = 5, q = 4
// 输出：5
//
// 示例 3：
// 输入：root = [1,2], p = 1, q = 2
// 输出：1

func lowestCommonAncestor(root, p, q *TreeNode) *TreeNode {
	if root == nil || root == p || root == q {
		return root
	}

	left := lowestCommonAncestor(root.Left, p, q)
	right := lowestCommonAncestor(root.Right, p, q)
	if left != nil && right != nil {
		return root
	}
	if left != nil {
		return left
	}
	return right
}

func TestLowestCommonAncestor(t *testing.T) {
	// 构建测试树:
	//       3
	//      / \
	//     5   1
	//    / \ / \
	//   6  2 0  8
	//     / \
	//    7   4
	node4 := &TreeNode{Val: 4}
	node7 := &TreeNode{Val: 7}
	node2 := &TreeNode{Val: 2, Left: node7, Right: node4}
	node6 := &TreeNode{Val: 6}
	node5 := &TreeNode{Val: 5, Left: node6, Right: node2}
	node0 := &TreeNode{Val: 0}
	node8 := &TreeNode{Val: 8}
	node1 := &TreeNode{Val: 1, Left: node0, Right: node8}
	root := &TreeNode{Val: 3, Left: node5, Right: node1}

	tests := []struct {
		name     string
		p        *TreeNode
		q        *TreeNode
		expected *TreeNode
	}{
		{name: "示例1: p=5,q=1", p: node5, q: node1, expected: root},
		{name: "示例2: p=5,q=4", p: node5, q: node4, expected: node5},
		{name: "示例3: p=6,q=4", p: node6, q: node4, expected: node5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := lowestCommonAncestor(root, tt.p, tt.q)
			if got != tt.expected {
				t.Errorf("lowestCommonAncestor() = %v, want %v", got.Val, tt.expected.Val)
			}
		})
	}
}

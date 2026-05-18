package top100liked

import "testing"

// 543. 二叉树的直径 (Diameter of Binary Tree)
//
// 题目描述:
// 给你一棵二叉树的根节点，返回该树的直径。
// 二叉树的直径是指树中任意两个节点之间最长路径的长度（边数）。
// 这条路径可能经过也可能不经过根节点。
//
// 示例 1：
// 输入：root = [1,2,3,4,5]
// 输出：3
// 解释：取路径 [4,2,1,3] 或 [5,2,1,3] 的长度
//
// 示例 2：
// 输入：root = [1,2]
// 输出：1
//
// 提示：对每个节点，经过它的最长路径 = 左子树深度 + 右子树深度，递归更新全局最大值

func diameterOfBinaryTree(root *TreeNode) int {
	diameter := 0
	var depth func(node *TreeNode) int
	depth = func(node *TreeNode) int {
		if node == nil {
			return 0
		}
		l := depth(node.Left)  // 左子树深度
		r := depth(node.Right) // 右子树深度
		// 以当前节点为"拐点"的路径长度 = 左深 + 右深（边数）
		if l+r > diameter {
			diameter = l + r
		}
		// 返回给父节点的是"以当前节点为顶"的最大深度
		return max(l, r) + 1
	}
	depth(root)
	return diameter
}

func TestDiameterOfBinaryTree(t *testing.T) {
	tests := []struct {
		name     string
		root     *TreeNode
		expected int
	}{
		{
			name: "示例1",
			root: &TreeNode{Val: 1,
				Left:  &TreeNode{Val: 2, Left: &TreeNode{Val: 4}, Right: &TreeNode{Val: 5}},
				Right: &TreeNode{Val: 3},
			},
			expected: 3,
		},
		{
			name:     "两个节点",
			root:     &TreeNode{Val: 1, Left: &TreeNode{Val: 2}},
			expected: 1,
		},
		{
			name:     "单节点",
			root:     &TreeNode{Val: 1},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := diameterOfBinaryTree(tt.root)
			if got != tt.expected {
				t.Errorf("diameterOfBinaryTree() = %v, want %v", got, tt.expected)
			}
		})
	}
}

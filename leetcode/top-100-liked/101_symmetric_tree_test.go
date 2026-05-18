package top100liked

import (
	"testing"
)

// 101. 对称二叉树 (Symmetric Tree)
//
// 题目描述:
// 给你一个二叉树的根节点 root，检查它是否轴对称。
//
// 示例 1：
// 输入：root = [1,2,2,3,4,4,3]
// 输出：true
//
// 示例 2：
// 输入：root = [1,2,2,null,3,null,3]
// 输出：false
//
// 提示：递归比较左右子树是否镜像（左的左 vs 右的右，左的右 vs 右的左）

func isSymmetric(root *TreeNode) bool {
	var isMirror func(left, right *TreeNode) bool
	isMirror = func(left, right *TreeNode) bool {
		if left == nil && right == nil {
			return true
		}
		if left == nil || right == nil {
			return false
		}
		if left.Val != right.Val {
			return false
		}
		return isMirror(left.Left, right.Right) && isMirror(left.Right, right.Left)
	}
	if root == nil {
		return true
	}
	return isMirror(root.Left, root.Right)
}

func TestIsSymmetric(t *testing.T) {
	tests := []struct {
		name     string
		root     *TreeNode
		expected bool
	}{
		{
			name: "对称树",
			root: &TreeNode{Val: 1,
				Left:  &TreeNode{Val: 2, Left: &TreeNode{Val: 3}, Right: &TreeNode{Val: 4}},
				Right: &TreeNode{Val: 2, Left: &TreeNode{Val: 4}, Right: &TreeNode{Val: 3}},
			},
			expected: true,
		},
		{
			name: "非对称树",
			root: &TreeNode{Val: 1,
				Left:  &TreeNode{Val: 2, Right: &TreeNode{Val: 3}},
				Right: &TreeNode{Val: 2, Right: &TreeNode{Val: 3}},
			},
			expected: false,
		},
		{
			name:     "空树",
			root:     nil,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isSymmetric(tt.root)
			if got != tt.expected {
				t.Errorf("isSymmetric() = %v, want %v", got, tt.expected)
			}
		})
	}
}

package top100liked

import (
	"math"
	"testing"
)

// 98. 验证二叉搜索树 (Validate Binary Search Tree)
//
// 题目描述:
// 给你一个二叉树的根节点 root，判断其是否是一个有效的二叉搜索树。
// - 节点的左子树只包含小于当前节点的数
// - 节点的右子树只包含大于当前节点的数
// - 所有左子树和右子树自身必须也是二叉搜索树
//
// 示例 1：
// 输入：root = [2,1,3]
// 输出：true
//
// 示例 2：
// 输入：root = [5,1,4,null,null,3,6]
// 输出：false
//
// 提示：递归传递上下界 (lower, upper)

func isValidBST(root *TreeNode) bool {
	var validate func(*TreeNode, int, int) bool
	validate = func(node *TreeNode, lower, upper int) bool {
		if node == nil {
			return true
		}
		if node.Val <= lower || node.Val >= upper {
			return false
		}
		return validate(node.Left, lower, node.Val) && validate(node.Right, node.Val, upper)
	}
	return validate(root, math.MinInt64, math.MaxInt64)
}

func TestIsValidBST(t *testing.T) {
	tests := []struct {
		name     string
		root     *TreeNode
		expected bool
	}{
		{
			name:     "有效BST",
			root:     &TreeNode{Val: 2, Left: &TreeNode{Val: 1}, Right: &TreeNode{Val: 3}},
			expected: true,
		},
		{
			name: "无效BST-右子树有小值",
			root: &TreeNode{Val: 5,
				Left: &TreeNode{Val: 1},
				Right: &TreeNode{Val: 4,
					Left:  &TreeNode{Val: 3},
					Right: &TreeNode{Val: 6},
				},
			},
			expected: false,
		},
		{
			name:     "单节点",
			root:     &TreeNode{Val: 1},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidBST(tt.root)
			if got != tt.expected {
				t.Errorf("isValidBST() = %v, want %v", got, tt.expected)
			}
		})
	}
}

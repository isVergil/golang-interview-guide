package top100liked

import (
	"math"
	"testing"
)

// 124. 二叉树中的最大路径和 (Binary Tree Maximum Path Sum)
//
// 题目描述:
// 二叉树中的路径被定义为一条节点序列，序列中每对相邻节点之间都存在一条边。
// 同一个节点在一条路径序列中至多出现一次。该路径至少包含一个节点，且不一定经过根节点。
// 路径和是路径中各节点值的总和。给你一个二叉树的根节点 root，返回其最大路径和。
//
// 示例 1：
// 输入：root = [1,2,3]
// 输出：6
// 解释：最优路径是 2 -> 1 -> 3，路径和为 2 + 1 + 3 = 6
//
// 示例 2：
// 输入：root = [-10,9,20,null,null,15,7]
// 输出：42
// 解释：最优路径是 15 -> 20 -> 7，路径和为 15 + 20 + 7 = 42
//
// 提示：递归返回"经过当前节点向下的最大贡献值"，同时更新全局最大路径和

func maxPathSum(root *TreeNode) int {
	maxSum := math.MinInt64

	var dfs func(node *TreeNode) int
	dfs = func(node *TreeNode) int {
		if node == nil {
			return 0
		}
		left := max(0, dfs(node.Left))
		right := max(0, dfs(node.Right))

		maxSum = max(maxSum, left+node.Val+right)

		return node.Val + max(left, right)
	}
	dfs(root)

	return maxSum
}

func TestMaxPathSum(t *testing.T) {
	tests := []struct {
		name     string
		root     *TreeNode
		expected int
	}{
		{
			name:     "示例1",
			root:     &TreeNode{Val: 1, Left: &TreeNode{Val: 2}, Right: &TreeNode{Val: 3}},
			expected: 6,
		},
		{
			name: "示例2",
			root: &TreeNode{Val: -10,
				Left: &TreeNode{Val: 9},
				Right: &TreeNode{Val: 20,
					Left:  &TreeNode{Val: 15},
					Right: &TreeNode{Val: 7},
				},
			},
			expected: 42,
		},
		{
			name:     "单负数",
			root:     &TreeNode{Val: -3},
			expected: -3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := maxPathSum(tt.root)
			if got != tt.expected {
				t.Errorf("maxPathSum() = %v, want %v", got, tt.expected)
			}
		})
	}
}

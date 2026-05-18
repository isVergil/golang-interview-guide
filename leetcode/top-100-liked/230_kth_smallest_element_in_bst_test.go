package top100liked

import "testing"

// 230. 二叉搜索树中第 K 小的元素 (Kth Smallest Element in a BST)
//
// 题目描述:
// 给定一个二叉搜索树的根节点 root，和一个整数 k，
// 请你设计一个算法查找其中第 k 个最小元素（从 1 开始计数）。
//
// 示例 1：
// 输入：root = [3,1,4,null,2], k = 1
// 输出：1
//
// 示例 2：
// 输入：root = [5,3,6,2,4,null,null,1], k = 3
// 输出：3
//
// 提示：BST 中序遍历是有序的，遍历到第 k 个就是答案
// 递归
func kthSmallest(root *TreeNode, k int) int {
	cnt, res := 0, 0
	var inorder func(*TreeNode) bool
	inorder = func(node *TreeNode) bool {
		if node == nil {
			return false
		}
		if inorder(node.Left) {
			return true
		}
		cnt++
		if cnt == k {
			res = node.Val
			return true
		}
		return inorder(node.Right)
	}
	inorder(root)
	return res
}

// 迭代
func kthSmallest1(root *TreeNode, k int) int {
	stack := make([]*TreeNode, 0)
	node := root
	for node != nil || len(stack) > 0 {
		for node != nil {
			stack = append(stack, node)
			node = node.Left
		}
		node = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		k--
		if k == 0 {
			return node.Val
		}
		node = node.Right
	}
	return -1
}

func TestKthSmallest(t *testing.T) {
	tests := []struct {
		name     string
		root     *TreeNode
		k        int
		expected int
	}{
		{
			name:     "示例1",
			root:     &TreeNode{Val: 3, Left: &TreeNode{Val: 1, Right: &TreeNode{Val: 2}}, Right: &TreeNode{Val: 4}},
			k:        1,
			expected: 1,
		},
		{
			name: "示例2",
			root: &TreeNode{Val: 5,
				Left: &TreeNode{Val: 3,
					Left:  &TreeNode{Val: 2, Left: &TreeNode{Val: 1}},
					Right: &TreeNode{Val: 4},
				},
				Right: &TreeNode{Val: 6},
			},
			k:        3,
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := kthSmallest(tt.root, tt.k)
			if got != tt.expected {
				t.Errorf("kthSmallest() = %v, want %v", got, tt.expected)
			}
		})

		t.Run(tt.name, func(t *testing.T) {
			got := kthSmallest1(tt.root, tt.k)
			if got != tt.expected {
				t.Errorf("kthSmallest1() = %v, want %v", got, tt.expected)
			}
		})
	}
}

package top100liked

import (
	"testing"
)

// 226. 翻转二叉树 (Invert Binary Tree)
//
// 题目描述:
// 给你一棵二叉树的根节点 root ，翻转这棵二叉树，并返回其根节点。
//
// 示例 1：
// 输入：root = [4,2,7,1,3,6,9]
// 输出：[4,7,2,9,6,3,1]
//
// 示例 2：
// 输入：root = [2,1,3]
// 输出：[2,3,1]
//
// 示例 3：
// 输入：root = []
// 输出：[]

func invertTree(root *TreeNode) *TreeNode {
	if root == nil {
		return nil
	}

	// 先交换还是后交换都可以，因为只是交换指针，顺序不影响结果
	// 交换左右子树
	root.Left, root.Right = root.Right, root.Left

	// 递归翻转左右子树
	invertTree(root.Left)
	invertTree(root.Right)

	return root
}

func invertTree1(root *TreeNode) *TreeNode {
	if root == nil {
		return nil
	}
	// 用栈模拟递归
	stack := []*TreeNode{root}
	for len(stack) > 0 {
		// 弹出栈顶节点
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		// 交换左右子树
		node.Left, node.Right = node.Right, node.Left
		// 子节点入栈，继续处理
		if node.Left != nil {
			stack = append(stack, node.Left)
		}
		if node.Right != nil {
			stack = append(stack, node.Right)
		}
	}
	return root
}

// 辅助函数：层序遍历二叉树，用于验证结果
func levelOrder(root *TreeNode) []int {
	if root == nil {
		return nil
	}
	var result []int
	queue := []*TreeNode{root}
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		result = append(result, node.Val)
		if node.Left != nil {
			queue = append(queue, node.Left)
		}
		if node.Right != nil {
			queue = append(queue, node.Right)
		}
	}
	return result
}

func TestInvertTree(t *testing.T) {
	tests := []struct {
		name     string
		root     *TreeNode
		expected []int
	}{
		{
			name: "示例1",
			root: &TreeNode{
				Val: 4,
				Left: &TreeNode{
					Val:   2,
					Left:  &TreeNode{Val: 1},
					Right: &TreeNode{Val: 3},
				},
				Right: &TreeNode{
					Val:   7,
					Left:  &TreeNode{Val: 6},
					Right: &TreeNode{Val: 9},
				},
			},
			expected: []int{4, 7, 2, 9, 6, 3, 1},
		},
		{
			name: "示例2",
			root: &TreeNode{
				Val:   2,
				Left:  &TreeNode{Val: 1},
				Right: &TreeNode{Val: 3},
			},
			expected: []int{2, 3, 1},
		},
		{
			name:     "空树",
			root:     nil,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := levelOrder(invertTree(tt.root))
			if len(got) != len(tt.expected) {
				t.Errorf("invertTree() level order = %v, want %v", got, tt.expected)
				return
			}
			for i := range got {
				if got[i] != tt.expected[i] {
					t.Errorf("invertTree() level order = %v, want %v", got, tt.expected)
					return
				}
			}
		})
	}
}

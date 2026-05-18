package top100liked

import (
	"reflect"
	"testing"
)

// 114. 二叉树展开为链表 (Flatten Binary Tree to Linked List)
//
// 题目描述:
// 给你二叉树的根节点 root，请你将它展开为一个单链表：
// - 展开后的单链表应该同样使用 TreeNode，其中 right 子指针指向链表中下一个结点，而左子指针始终为 nil
// - 展开后的单链表应该与二叉树先序遍历顺序相同
//
// 示例 1：
// 输入：root = [1,2,5,3,4,null,6]
// 输出：[1,null,2,null,3,null,4,null,5,null,6]
//
// 示例 2：
// 输入：root = []
// 输出：[]
//
// 提示：原地展开，空间复杂度 O(1)

func flatten(root *TreeNode) {
	cur := root
	for cur != nil {
		if cur.Left != nil {
			prev := cur.Left
			for prev.Right != nil {
				prev = prev.Right
			}
			//把右子树挂到左子树的最右节点
			prev.Right = cur.Right
			cur.Right = cur.Left
			cur.Left = nil
		}
		cur = cur.Right
	}
}

func TestFlatten(t *testing.T) {
	tests := []struct {
		name     string
		root     *TreeNode
		expected []int
	}{
		{
			name: "示例1",
			root: &TreeNode{Val: 1,
				Left:  &TreeNode{Val: 2, Left: &TreeNode{Val: 3}, Right: &TreeNode{Val: 4}},
				Right: &TreeNode{Val: 5, Right: &TreeNode{Val: 6}},
			},
			expected: []int{1, 2, 3, 4, 5, 6},
		},
		{
			name:     "空树",
			root:     nil,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flatten(tt.root)
			got := []int{}
			for node := tt.root; node != nil; node = node.Right {
				got = append(got, node.Val)
			}
			if tt.root == nil {
				got = nil
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("flatten() = %v, want %v", got, tt.expected)
			}
		})
	}
}

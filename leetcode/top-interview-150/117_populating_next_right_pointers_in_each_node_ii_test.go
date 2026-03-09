package topinterview150

import (
	"testing"
)

// 117. 填充每个节点的下一个右侧节点指针 II (Populating Next Right Pointers in Each Node II)
//
// 题目描述:
// 给定一个二叉树：
// struct Node {
//   int val;
//   Node *left;
//   Node *right;
//   Node *next;
// }
// 填充它的每个 next 指针，让这个指针指向其下一个右侧节点。如果找不到下一个右侧节点，则将 next 指针设置为 NULL 。
// 初始状态下，所有 next 指针都被设置为 NULL 。
//
// 示例 1：
// 输入：root = [1,2,3,4,5,null,7]
// 输出：[1,#,2,3,#,4,5,7,#]
//
// 示例 2：
// 输入：root = []
// 输出：[]

// Node117 defines a binary tree node with a next pointer.
type Node117 struct {
	Val   int
	Left  *Node117
	Right *Node117
	Next  *Node117
}

func connect(root *Node117) *Node117 {
	panic("not implemented")
}

func TestConnect(t *testing.T) {
	// 二叉树测试
}

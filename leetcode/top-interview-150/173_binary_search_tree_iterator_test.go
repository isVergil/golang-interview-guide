package topinterview150

import (
	"testing"
)

// 173. 二叉搜索树迭代器 (Binary Search Tree Iterator)
//
// 题目描述:
// 实现一个二叉搜索树迭代器类 BSTIterator ，表示一个按中序遍历二叉搜索树（BST）的迭代器：
// BSTIterator(TreeNode root) 初始化 BSTIterator 类的一个对象。BST 的根节点 root 会作为构造函数的一部分给出。
// 指针应初始化为一个不存在于 BST 中的最小元素（即，小于树中任何元素的数字）。
// boolean hasNext() 如果向指针右侧遍历，树中依然存在节点，则返回 true ；否则返回 false 。
// int next() 将指针向右移动，然后返回指针指向处的节点的值。
//
// 示例：
// 输入
// ["BSTIterator", "next", "next", "hasNext", "next", "hasNext", "next", "hasNext"]
// [[[7, 3, 15, null, null, 9, 20]], [], [], [], [], [], [], []]
// 输出
// [null, 3, 7, true, 9, true, 15, false]

type BSTIterator struct {
}

func ConstructorBSTIterator(root *TreeNode) BSTIterator {
	panic("not implemented")
}

func (this *BSTIterator) Next() int {
	panic("not implemented")
}

func (this *BSTIterator) HasNext() bool {
	panic("not implemented")
}

func TestBSTIterator(t *testing.T) {
	// 二叉树迭代器测试
}

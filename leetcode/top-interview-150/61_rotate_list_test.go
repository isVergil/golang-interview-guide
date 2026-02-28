package topinterview150

import (
	"testing"
)

// 61. 旋转链表 (Rotate List)
//
// 题目描述:
// 给你一个链表的头节点 head ，旋转链表，将链表每个节点向右移动 k 个位置。
//
// 示例 1：
// 输入：head = [1,2,3,4,5], k = 2
// 输出：[4,5,1,2,3]
//
// 示例 2：
// 输入：head = [0,1,2], k = 4
// 输出：[2,0,1]

// 可以理解成 把 n - k%n 后面的链表移到最前面
func rotateRight(head *ListNode, k int) *ListNode {
	if head == nil || head.Next == nil || k == 0 {
		return head
	}

	// 计算链表长度
	n := 1
	cur := head
	for cur.Next != nil {
		cur = cur.Next
		n++
	}

	// 链表成环
	cur.Next = head

	// 计算需要移动的次数
	k = k % n
	steps := n - k
	for i := 0; i < steps; i++ {
		head = head.Next
	}

	newHead := head.Next
	head.Next = nil
	return newHead

}

func TestRotateRight(t *testing.T) {
	// 链表测试通常需要辅助函数
}

package top100liked

import (
	"testing"
)

// 142. 环形链表 II (Linked List Cycle II)
//
// 题目描述:
// 给定一个链表的头节点 head ，返回链表开始入环的第一个节点。如果链表无环，则返回 null。
//
// 示例 1：
// 输入：head = [3,2,0,-4], pos = 1
// 输出：返回索引为 1 的链表节点
// 解释：链表中有一个环，其尾部连接到第二个节点。
//
// 示例 2：
// 输入：head = [1,2], pos = 0
// 输出：返回索引为 0 的链表节点
// 解释：链表中有一个环，其尾部连接到第一个节点。
//
// 示例 3：
// 输入：head = [1], pos = -1
// 输出：返回 null
// 解释：链表中没有环。

func detectCycle(head *ListNode) *ListNode {
	slow, fast := head, head
	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
		if slow == fast {
			res := head
			for res != slow {
				res = res.Next
				slow = slow.Next
			}
			return res
		}
	}
	return nil
}

func TestDetectCycle(t *testing.T) {
	// 测试1：有环，入环点在索引1
	node1 := &ListNode{Val: 3}
	node2 := &ListNode{Val: 2}
	node3 := &ListNode{Val: 0}
	node4 := &ListNode{Val: -4}
	node1.Next = node2
	node2.Next = node3
	node3.Next = node4
	node4.Next = node2 // 环入口是 node2

	got := detectCycle(node1)
	if got != node2 {
		t.Errorf("示例1: expected node with val 2, got %v", got)
	}

	// 测试2：无环
	head := &ListNode{Val: 1}
	got = detectCycle(head)
	if got != nil {
		t.Errorf("示例3: expected nil, got %v", got)
	}
}

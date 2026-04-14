package top100liked

import (
	"testing"
)

// 141. 环形链表 (Linked List Cycle)
//
// 题目描述:
// 给你一个链表的头节点 head ，判断链表中是否有环。
// 如果链表中存在环 ，则返回 true 。 否则，返回 false 。
//
// 示例 1：
// 输入：head = [3,2,0,-4], pos = 1
// 输出：true
// 解释：链表中有一个环，其尾部连接到第二个节点。
//
// 示例 2：
// 输入：head = [1,2], pos = 0
// 输出：true
//
// 示例 3：
// 输入：head = [1], pos = -1
// 输出：false

func hasCycle(head *ListNode) bool {
	// 快慢指针都从头出发
	slow, fast := head, head
	for fast != nil && fast.Next != nil {
		slow = slow.Next      // 慢指针走一步
		fast = fast.Next.Next // 快指针走两步
		// 如果有环，快指针终会追上慢指针
		if slow == fast {
			return true
		}
	}
	// 快指针走到末尾，说明无环
	return false
}

func TestHasCycle(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
	}{
		{
			name:     "有环",
			expected: true,
		},
		{
			name:     "无环",
			expected: false,
		},
		{
			name:     "空链表",
			expected: false,
		},
	}

	// 构建有环链表: 3 -> 2 -> 0 -> -4 -> 2(回到第二个节点)
	node4 := &ListNode{Val: -4}
	node3 := &ListNode{Val: 0, Next: node4}
	node2 := &ListNode{Val: 2, Next: node3}
	node1 := &ListNode{Val: 3, Next: node2}
	node4.Next = node2 // 形成环

	// 构建无环链表: 1 -> 2
	noLoop := &ListNode{Val: 1, Next: &ListNode{Val: 2}}

	inputs := []*ListNode{node1, noLoop, nil}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasCycle(inputs[i])
			if got != tt.expected {
				t.Errorf("hasCycle() = %v, want %v", got, tt.expected)
			}
		})
	}
}

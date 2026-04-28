package top100liked

import (
	"testing"
)

// 25. K 个一组翻转链表 (Reverse Nodes in k-Group)
//
// 题目描述:
// 给你链表的头节点 head ，每 k 个节点一组进行翻转，请你返回修改后的链表。
// k 是一个正整数，它的值小于或等于链表的长度。如果节点总数不是 k 的整数倍，
// 那么请将最后剩余的节点保持原有顺序。不能只是单纯改变节点内部的值，而是需要实际进行节点交换。
//
// 示例 1：
// 输入：head = [1,2,3,4,5], k = 2
// 输出：[2,1,4,3,5]
//
// 示例 2：
// 输入：head = [1,2,3,4,5], k = 3
// 输出：[3,2,1,4,5]

func reverseKGroup(head *ListNode, k int) *ListNode {
	dummy := &ListNode{Next: head}
	prev := dummy
	for {
		// 先探底
		end := prev
		for i := 0; i < k; i++ {
			end = end.Next
			if end == nil {
				return dummy.Next
			}
		}

		// 当前组的头尾节点
		start := prev.Next
		nxt := end.Next

		end.Next = nil
		reverseListNode(start)

		prev.Next = end
		start.Next = nxt
		prev = start
	}
}

func reverseListNode(head *ListNode) {
	var prev *ListNode
	cur := head
	for cur != nil {
		next := cur.Next
		cur.Next = prev
		prev = cur
		cur = next
	}
}

func TestReverseKGroup(t *testing.T) {
	toList := func(nums []int) *ListNode {
		dummy := &ListNode{}
		cur := dummy
		for _, n := range nums {
			cur.Next = &ListNode{Val: n}
			cur = cur.Next
		}
		return dummy.Next
	}
	toSlice := func(head *ListNode) []int {
		var res []int
		for head != nil {
			res = append(res, head.Val)
			head = head.Next
		}
		return res
	}

	tests := []struct {
		name     string
		head     []int
		k        int
		expected []int
	}{
		{name: "示例1", head: []int{1, 2, 3, 4, 5}, k: 2, expected: []int{2, 1, 4, 3, 5}},
		{name: "示例2", head: []int{1, 2, 3, 4, 5}, k: 3, expected: []int{3, 2, 1, 4, 5}},
		{name: "k=1", head: []int{1, 2, 3}, k: 1, expected: []int{1, 2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toSlice(reverseKGroup(toList(tt.head), tt.k))
			if len(got) != len(tt.expected) {
				t.Errorf("reverseKGroup() = %v, want %v", got, tt.expected)
				return
			}
			for i := range got {
				if got[i] != tt.expected[i] {
					t.Errorf("reverseKGroup() = %v, want %v", got, tt.expected)
					return
				}
			}
		})
	}
}

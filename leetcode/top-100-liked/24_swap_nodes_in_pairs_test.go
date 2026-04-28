package top100liked

import (
	"testing"
)

// 24. 两两交换链表中的节点 (Swap Nodes in Pairs)
//
// 题目描述:
// 给你一个链表，两两交换其中相邻的节点，并返回交换后链表的头节点。
// 你必须在不修改节点内部的值的情况下完成本题（即只能进行节点交换）。
//
// 示例 1：
// 输入：head = [1,2,3,4]
// 输出：[2,1,4,3]
//
// 示例 2：
// 输入：head = []
// 输出：[]
//
// 示例 3：
// 输入：head = [1]
// 输出：[1]

func swapPairs(head *ListNode) *ListNode {
	dummyNode := &ListNode{Next: head}
	prev := dummyNode
	for prev.Next != nil && prev.Next.Next != nil {
		node1 := prev.Next
		node2 := prev.Next.Next
		prev.Next = node2
		node1.Next = node2.Next
		node2.Next = node1
		prev = node1
	}
	return dummyNode.Next
}

func TestSwapPairs(t *testing.T) {
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
		expected []int
	}{
		{name: "示例1", head: []int{1, 2, 3, 4}, expected: []int{2, 1, 4, 3}},
		{name: "空链表", head: []int{}, expected: nil},
		{name: "单节点", head: []int{1}, expected: []int{1}},
		{name: "奇数个", head: []int{1, 2, 3}, expected: []int{2, 1, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toSlice(swapPairs(toList(tt.head)))
			if len(got) != len(tt.expected) {
				t.Errorf("swapPairs() = %v, want %v", got, tt.expected)
				return
			}
			for i := range got {
				if got[i] != tt.expected[i] {
					t.Errorf("swapPairs() = %v, want %v", got, tt.expected)
					return
				}
			}
		})
	}
}

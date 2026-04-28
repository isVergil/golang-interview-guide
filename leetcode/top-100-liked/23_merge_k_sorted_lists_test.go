package top100liked

import (
	"testing"
)

// 23. 合并 K 个升序链表 (Merge k Sorted Lists)
//
// 题目描述:
// 给你一个链表数组，每个链表都已经按升序排列。请你将所有链表合并到一个升序链表中，返回合并后的链表。
//
// 示例 1：
// 输入：lists = [[1,4,5],[1,3,4],[2,6]]
// 输出：[1,1,2,3,4,4,5,6]
//
// 示例 2：
// 输入：lists = []
// 输出：[]

func mergeKLists(lists []*ListNode) *ListNode {
	n := len(lists)
	if n == 0 {
		return nil
	}
	if n == 1 {
		return lists[0]
	}
	mid := n / 2
	left := mergeKLists(lists[:mid])
	right := mergeKLists(lists[mid:])
	return mergeList(left, right)
}

// merge 合并两个有序链表（就是第 21 题）
func mergeList(l1, l2 *ListNode) *ListNode {
	dummy := &ListNode{}
	cur := dummy
	for l1 != nil && l2 != nil {
		if l1.Val <= l2.Val {
			cur.Next = l1
			l1 = l1.Next
		} else {
			cur.Next = l2
			l2 = l2.Next
		}
		cur = cur.Next
	}
	if l1 != nil {
		cur.Next = l1
	} else {
		cur.Next = l2
	}
	return dummy.Next
}

func TestMergeKLists(t *testing.T) {
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
		lists    [][]int
		expected []int
	}{
		{name: "示例1", lists: [][]int{{1, 4, 5}, {1, 3, 4}, {2, 6}}, expected: []int{1, 1, 2, 3, 4, 4, 5, 6}},
		{name: "空列表", lists: [][]int{}, expected: nil},
		{name: "单个链表", lists: [][]int{{1, 2, 3}}, expected: []int{1, 2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var lists []*ListNode
			for _, l := range tt.lists {
				lists = append(lists, toList(l))
			}
			got := toSlice(mergeKLists(lists))
			if len(got) != len(tt.expected) {
				t.Errorf("mergeKLists() = %v, want %v", got, tt.expected)
				return
			}
			for i := range got {
				if got[i] != tt.expected[i] {
					t.Errorf("mergeKLists() = %v, want %v", got, tt.expected)
					return
				}
			}
		})
	}
}

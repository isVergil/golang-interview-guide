package top100liked

import (
	"reflect"
	"testing"
)

// 148. 排序链表 (Sort List)
//
// 题目描述:
// 给你链表的头结点 head，请将其按升序排列并返回排序后的链表。
// 要求时间复杂度 O(n log n)，空间复杂度 O(1)。
//
// 示例 1：
// 输入：head = [4,2,1,3]
// 输出：[1,2,3,4]
//
// 示例 2：
// 输入：head = [-1,5,3,4,0]
// 输出：[-1,0,3,4,5]
//
// 提示：归并排序（自底向上可做到 O(1) 空间）

func sortList(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}

	// 快慢指针找中点
	slow, fast := head, head.Next
	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
	}

	// 从中点断开
	mid := slow.Next
	slow.Next = nil

	// 递归排序两半
	left := sortList(head)
	right := sortList(mid)

	// 合并两个有序链表
	return mergeTwoLists148(left, right)
}

// 合并两个有序链表
func mergeTwoLists148(l1, l2 *ListNode) *ListNode {
	dummyNode := &ListNode{}
	cur := dummyNode
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
	}
	if l2 != nil {
		cur.Next = l2
	}
	return dummyNode.Next
}

func TestSortList(t *testing.T) {
	// 辅助函数：数组转链表
	toList := func(nums []int) *ListNode {
		dummy := &ListNode{}
		cur := dummy
		for _, v := range nums {
			cur.Next = &ListNode{Val: v}
			cur = cur.Next
		}
		return dummy.Next
	}
	// 辅助函数：链表转数组
	toSlice := func(head *ListNode) []int {
		res := []int{}
		for head != nil {
			res = append(res, head.Val)
			head = head.Next
		}
		return res
	}

	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{name: "示例1", input: []int{4, 2, 1, 3}, expected: []int{1, 2, 3, 4}},
		{name: "示例2", input: []int{-1, 5, 3, 4, 0}, expected: []int{-1, 0, 3, 4, 5}},
		{name: "空链表", input: []int{}, expected: []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			head := toList(tt.input)
			got := toSlice(sortList(head))
			if len(got) == 0 {
				got = []int{}
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("sortList() = %v, want %v", got, tt.expected)
			}
		})
	}
}

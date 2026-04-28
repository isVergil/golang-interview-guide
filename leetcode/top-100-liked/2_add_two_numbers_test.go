package top100liked

import (
	"testing"
)

// 2. 两数相加 (Add Two Numbers)
//
// 题目描述:
// 给你两个非空的链表，表示两个非负的整数。它们每位数字都是按照逆序的方式存储的，
// 并且每个节点只能存储一位数字。请你将两个数相加，并以相同形式返回一个表示和的链表。
//
// 示例 1：
// 输入：l1 = [2,4,3], l2 = [5,6,4]
// 输出：[7,0,8]（342 + 465 = 807）
//
// 示例 2：
// 输入：l1 = [9,9,9,9,9,9,9], l2 = [9,9,9,9]
// 输出：[8,9,9,9,0,0,0,1]

func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
	dummyNode := &ListNode{}
	cur := dummyNode
	carry := 0
	for l1 != nil || l2 != nil || carry > 0 {
		val := carry
		if l1 != nil {
			val += l1.Val
			l1 = l1.Next
		}
		if l2 != nil {
			val += l2.Val
			l2 = l2.Next
		}
		cur.Next = &ListNode{Val: val % 10}
		cur = cur.Next
		carry = val / 10
	}
	return dummyNode.Next
}

func TestAddTwoNumbers(t *testing.T) {
	// 辅助函数：切片转链表
	toList := func(nums []int) *ListNode {
		dummy := &ListNode{}
		cur := dummy
		for _, n := range nums {
			cur.Next = &ListNode{Val: n}
			cur = cur.Next
		}
		return dummy.Next
	}
	// 辅助函数：链表转切片
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
		l1, l2   []int
		expected []int
	}{
		{name: "示例1", l1: []int{2, 4, 3}, l2: []int{5, 6, 4}, expected: []int{7, 0, 8}},
		{name: "示例2", l1: []int{9, 9, 9, 9, 9, 9, 9}, l2: []int{9, 9, 9, 9}, expected: []int{8, 9, 9, 9, 0, 0, 0, 1}},
		{name: "零加零", l1: []int{0}, l2: []int{0}, expected: []int{0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toSlice(addTwoNumbers(toList(tt.l1), toList(tt.l2)))
			if len(got) != len(tt.expected) {
				t.Errorf("addTwoNumbers() = %v, want %v", got, tt.expected)
				return
			}
			for i := range got {
				if got[i] != tt.expected[i] {
					t.Errorf("addTwoNumbers() = %v, want %v", got, tt.expected)
					return
				}
			}
		})
	}
}

package top100liked

import "testing"

// 234. 回文链表 (Palindrome Linked List)
//
// 题目描述:
// 给你一个单链表的头节点 head，请你判断该链表是否为回文链表。
// 如果是，返回 true；否则，返回 false。
// 要求时间复杂度 O(n)，空间复杂度 O(1)。
//
// 示例 1：
// 输入：head = [1,2,2,1]
// 输出：true
//
// 示例 2：
// 输入：head = [1,2]
// 输出：false
//
// 提示：快慢指针找中点 + 反转后半段 + 逐一比较

func isPalindromeList(head *ListNode) bool {
	if head == nil || head.Next == nil {
		return true
	}
	slow, fast := head, head
	for fast.Next != nil && fast.Next.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
	}
	second := reverseListForIsPalindromeList(slow.Next)
	p1, p2 := head, second
	for p2 != nil {
		if p1.Val != p2.Val {
			return false
		}
		p1 = p1.Next
		p2 = p2.Next
	}
	return true
}

// 反转链表
func reverseListForIsPalindromeList(head *ListNode) *ListNode {
	var prev *ListNode
	cur := head
	for cur != nil {
		next := cur.Next
		cur.Next = prev
		prev = cur
		cur = next
	}
	return prev
}

func TestIsPalindromeList(t *testing.T) {
	toList := func(nums []int) *ListNode {
		dummy := &ListNode{}
		cur := dummy
		for _, v := range nums {
			cur.Next = &ListNode{Val: v}
			cur = cur.Next
		}
		return dummy.Next
	}

	tests := []struct {
		name     string
		input    []int
		expected bool
	}{
		{name: "回文-偶数", input: []int{1, 2, 2, 1}, expected: true},
		{name: "非回文", input: []int{1, 2}, expected: false},
		{name: "回文-奇数", input: []int{1, 2, 1}, expected: true},
		{name: "单节点", input: []int{1}, expected: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isPalindromeList(toList(tt.input))
			if got != tt.expected {
				t.Errorf("isPalindromeList() = %v, want %v", got, tt.expected)
			}
		})
	}
}

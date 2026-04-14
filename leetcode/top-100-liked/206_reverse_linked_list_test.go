package top100liked

import (
	"reflect"
	"testing"
)

// 206. 反转链表 (Reverse Linked List)
//
// 题目描述:
// 给你单链表的头节点 head ，请你反转链表，并返回反转后的链表。
//
// 示例 1：
// 输入：head = [1,2,3,4,5]
// 输出：[5,4,3,2,1]
//
// 示例 2：
// 输入：head = [1,2]
// 输出：[2,1]
//
// 示例 3：
// 输入：head = []
// 输出：[]

func reverseList(head *ListNode) *ListNode {
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

func reverseList1(head *ListNode) *ListNode {
	// 终止条件：空链表或只有一个节点
	if head == nil || head.Next == nil {
		return head
	}

	// 递归反转后面的部分，newHead 始终是原链表的最后一个节点
	newHead := reverseList(head.Next)

	// 核心操作：让“下一个节点”指向“自己”
	head.Next.Next = head
	// 斩断原本的指向，防止环路
	head.Next = nil

	return newHead
}

func TestReverseList(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{
			name:     "示例1",
			input:    []int{1, 2, 3, 4, 5},
			expected: []int{5, 4, 3, 2, 1},
		},
		{
			name:     "示例2",
			input:    []int{1, 2},
			expected: []int{2, 1},
		},
		{
			name:     "空链表",
			input:    nil,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			head := sliceToList(tt.input)
			got := listToSlice(reverseList(head))
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("reverseList() = %v, want %v", got, tt.expected)
			}
		})
	}
}

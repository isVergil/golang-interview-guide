package topinterview150

import (
	"testing"
)

// 21. 合并两个有序链表 (Merge Two Sorted Lists)
//
// 题目描述:
// 将两个升序链表合并为一个新的 升序 链表并返回。新链表是通过拼接给定的两个链表的所有节点组成的。
//
// 示例 1：
// 输入：l1 = [1,2,4], l2 = [1,3,4]
// 输出：[1,1,2,3,4,4]
//
// 示例 2：
// 输入：l1 = [], l2 = []
// 输出：[]
//
// 示例 3：
// 输入：l1 = [], l2 = [0]
// 输出：[0]

/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */

func mergeTwoLists(list1 *ListNode, list2 *ListNode) *ListNode {
	head := &ListNode{}
	cur := head
	for list1 != nil || list2 != nil {
		if list1 == nil {
			cur.Next = list2
			break
		}
		if list2 == nil {
			cur.Next = list1
			break
		}
		if list1.Val > list2.Val {
			cur.Next = &ListNode{list2.Val, nil}
			list2 = list2.Next
		} else {
			cur.Next = &ListNode{list1.Val, nil}
			list1 = list1.Next
		}
		cur = cur.Next
	}
	return head.Next
}

func TestMergeTwoLists(t *testing.T) {
	// 辅助函数和测试逻辑与 AddTwoNumbers 类似
	tests := []struct {
		name     string
		l1       []int
		l2       []int
		expected []int
	}{
		{
			name:     "Example 1",
			l1:       []int{1, 2, 4},
			l2:       []int{1, 3, 4},
			expected: []int{1, 1, 2, 3, 4, 4},
		},
		{
			name:     "Example 2",
			l1:       []int{},
			l2:       []int{},
			expected: []int{},
		},
		{
			name:     "Example 3",
			l1:       []int{},
			l2:       []int{0},
			expected: []int{0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: Uncomment when implementation is ready
			// l1 := sliceToList(tt.l1)
			// l2 := sliceToList(tt.l2)
			// got := mergeTwoLists(l1, l2)
			// gotSlice := listToSlice(got)
			// if !reflect.DeepEqual(gotSlice, tt.expected) {
			// 	t.Errorf("mergeTwoLists() = %v, want %v", gotSlice, tt.expected)
			// }
		})
	}
}

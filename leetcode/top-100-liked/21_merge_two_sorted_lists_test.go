package top100liked

import (
	"reflect"
	"testing"
)

// 21. 合并两个有序链表 (Merge Two Sorted Lists)
//
// 题目描述:
// 将两个升序链表合并为一个新的升序链表并返回。新链表是通过拼接给定的两个链表的所有节点组成的。
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

func mergeTwoLists(list1 *ListNode, list2 *ListNode) *ListNode {
	dummy := &ListNode{}
	cur := dummy
	for list1 != nil && list2 != nil {
		if list1.Val <= list2.Val {
			cur.Next = &ListNode{Val: list1.Val}
			list1 = list1.Next
		} else {
			cur.Next = &ListNode{Val: list2.Val}
			list2 = list2.Next
		}
		cur = cur.Next
	}
	if list1 == nil {
		cur.Next = list2
	}

	if list2 == nil {
		cur.Next = list1
	}
	return dummy.Next
}

// 辅助函数：将切片转为链表
func sliceToList(nums []int) *ListNode {
	dummy := &ListNode{}
	curr := dummy
	for _, v := range nums {
		curr.Next = &ListNode{Val: v}
		curr = curr.Next
	}
	return dummy.Next
}

// 辅助函数：将链表转为切片
func listToSlice(head *ListNode) []int {
	var result []int
	for head != nil {
		result = append(result, head.Val)
		head = head.Next
	}
	return result
}

func TestMergeTwoLists(t *testing.T) {
	tests := []struct {
		name     string
		l1       []int
		l2       []int
		expected []int
	}{
		{
			name:     "示例1",
			l1:       []int{1, 2, 4},
			l2:       []int{1, 3, 4},
			expected: []int{1, 1, 2, 3, 4, 4},
		},
		{
			name:     "示例2-两个空链表",
			l1:       nil,
			l2:       nil,
			expected: nil,
		},
		{
			name:     "示例3-一个空链表",
			l1:       nil,
			l2:       []int{0},
			expected: []int{0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l1 := sliceToList(tt.l1)
			l2 := sliceToList(tt.l2)
			got := listToSlice(mergeTwoLists(l1, l2))
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("mergeTwoLists() = %v, want %v", got, tt.expected)
			}
		})
	}
}

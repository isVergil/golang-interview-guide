package top100liked

import (
	"reflect"
	"testing"
)

// 19. 删除链表的倒数第 N 个结点 (Remove Nth Node From End of List)
//
// 题目描述:
// 给你一个链表，删除链表的倒数第 n 个结点，并且返回链表的头结点。
//
// 示例 1：
// 输入：head = [1,2,3,4,5], n = 2
// 输出：[1,2,3,5]
//
// 示例 2：
// 输入：head = [1], n = 1
// 输出：[]
//
// 示例 3：
// 输入：head = [1,2], n = 1
// 输出：[1]

func removeNthFromEnd(head *ListNode, n int) *ListNode {
	// 哨兵节点，处理删除头节点的边界情况
	dummy := &ListNode{Next: head}
	slow, fast := dummy, dummy

	// 快指针先走 n+1 步，拉开间距
	for i := 0; i <= n; i++ {
		fast = fast.Next
	}

	// 快慢同步走，快指针到末尾时，慢指针刚好在目标前一个
	for fast != nil {
		slow = slow.Next
		fast = fast.Next
	}

	// 跳过目标节点
	slow.Next = slow.Next.Next
	return dummy.Next
}

func TestRemoveNthFromEnd(t *testing.T) {
	tests := []struct {
		name     string
		head     []int
		n        int
		expected []int
	}{
		{
			name:     "示例1",
			head:     []int{1, 2, 3, 4, 5},
			n:        2,
			expected: []int{1, 2, 3, 5},
		},
		{
			name:     "示例2",
			head:     []int{1},
			n:        1,
			expected: nil,
		},
		{
			name:     "示例3",
			head:     []int{1, 2},
			n:        1,
			expected: []int{1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			head := sliceToList(tt.head)
			got := listToSlice(removeNthFromEnd(head, tt.n))
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("removeNthFromEnd() = %v, want %v", got, tt.expected)
			}
		})
	}
}

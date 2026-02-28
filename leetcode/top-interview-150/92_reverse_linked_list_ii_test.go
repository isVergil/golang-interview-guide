package topinterview150

import (
	"testing"
)

// 92. 反转链表 II (Reverse Linked List II)
//
// 题目描述:
// 给你单链表的头指针 head 和两个整数 left 和 right ，其中 left <= right 。请你反转从位置 left 到位置 right 的链表节点，返回 反转后的链表 。
//
// 示例 1：
// 输入：head = [1,2,3,4,5], left = 2, right = 4
// 输出：[1,4,3,2,5]
//
// 示例 2：
// 输入：head = [5], left = 1, right = 1
// 输出：[5]

func reverseBetween(head *ListNode, left int, right int) *ListNode {
	dummyNode := &ListNode{Next: head}
	prev := dummyNode

	// 将 pre 移动到 left 的前一个位置
	for i := 0; i < left-1; i++ {
		prev = prev.Next
	}

	cur := prev.Next
	for i := 0; i < right-left; i++ {
		next := cur.Next
		cur.Next = next.Next
		next.Next = prev.Next
		prev.Next = next
	}

	return dummyNode.Next

}

func TestReverseBetween(t *testing.T) {
	// 链表测试通常需要辅助函数
}

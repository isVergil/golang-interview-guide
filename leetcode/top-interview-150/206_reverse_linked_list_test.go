package topinterview150

import (
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

// 1 迭代
func reverseList(head *ListNode) *ListNode {
	var prev *ListNode // 初始化前驱为 nil
	curr := head       // 从头节点开始

	for curr != nil {
		//curr.Next, prev, curr = prev, curr, curr.Next

		next := curr.Next // 1. 临时保存下一个节点，防止断链
		curr.Next = prev  // 2. 反转：让当前节点指向前一个节点

		// 3. 指针整体向后移动一位
		prev = curr
		curr = next
	}

	// 最后 prev 指向的就是原链表的末尾，即新链表的头
	return prev
}

// 2 递归 先反转后面的链表，然后再处理当前节点。
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
	// 链表测试通常需要辅助函数
}

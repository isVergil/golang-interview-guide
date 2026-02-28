package topinterview150

import (
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
	// 1. 创建哨兵节点，指向 head
	dummy := &ListNode{Next: head}

	// 2. 初始化快慢指针都指向哨兵
	fast := dummy
	slow := dummy

	// 3. 快指针先走 n + 1 步
	// 为什么要走 n+1 步？为了让 slow 最后停在待删除节点的前一个位置
	for i := 0; i <= n; i++ {
		fast = fast.Next
	}

	// 4. 同时移动 fast 和 slow，直到 fast 走到链表末尾 (nil)
	for fast != nil {
		fast = fast.Next
		slow = slow.Next
	}

	// 5. 此时 slow.Next 就是倒数第 n 个节点，直接跳过它
	slow.Next = slow.Next.Next

	// 6. 返回哨兵节点的下一个节点（即真正的头节点）
	return dummy.Next
}

func TestRemoveNthFromEnd(t *testing.T) {
	// 链表测试通常需要辅助函数
}

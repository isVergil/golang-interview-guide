package topinterview150

import (
	"testing"
)

// 25. K 个一组翻转链表 (Reverse Nodes in k-Group)
//
// 题目描述:
// 给你链表的头节点 head ，每 k 个节点一组进行翻转，请你返回修改后的链表。
// k 是一个正整数，它的值小于或等于链表的长度。如果节点总数不是 k 的整数倍，那么请将最后剩余的节点保持原有顺序。
// 你不能只是单纯的改变节点内部的值，而是需要实际进行节点翻转。
//
// 示例 1：
// 输入：head = [1,2,3,4,5], k = 2
// 输出：[2,1,4,3,5]
//
// 示例 2：
// 输入：head = [1,2,3,4,5], k = 3
// 输出：[3,2,1,4,5]

func reverseKGroup(head *ListNode, k int) *ListNode {
	// 1. 检查长度是否满足 k 个
	cur := head
	for i := 0; i < k; i++ {
		if cur == nil {
			return head
		}
		cur = cur.Next
	}

	// 2. 局部反转（翻转前 k 个节点）
	var prev *ListNode
	cur = head
	for i := 0; i < k; i++ {
		tmp := cur.Next
		cur.Next = prev
		prev = cur
		cur = tmp
	}

	// 3. 递归连接
	// 翻转完后，原本的 head 变成了这组的“尾巴”
	// 此时的 curr 刚好指向下一组的开头
	head.Next = reverseKGroup(cur, k)

	// prev 指向的是这组翻转后的“新头”
	return prev

}

func TestReverseKGroup(t *testing.T) {
	// 链表测试通常需要辅助函数
}

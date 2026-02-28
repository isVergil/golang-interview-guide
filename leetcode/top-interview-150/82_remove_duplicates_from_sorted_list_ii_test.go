package topinterview150

import (
	"testing"
)

// 82. 删除排序链表中的重复元素 II (Remove Duplicates from Sorted List II)
//
// 题目描述:
// 给定一个已排序的链表的头 head ， 删除所有含有重复数字的节点，只保留原始链表中 没有重复出现 的数字 。返回 已排序的链表 。
//
// 示例 1：
// 输入：head = [1,2,3,3,4,4,5]
// 输出：[1,2,5]
//
// 示例 2：
// 输入：head = [1,1,1,2,3]
// 输出：[2,3]

func deleteDuplicates(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}

	dummyNode := &ListNode{Next: head}
	cur := dummyNode

	for cur.Next != nil && cur.Next.Next != nil {
		// 如果发现紧接着的两个节点值相等
		if cur.Next.Val == cur.Next.Next.Val {
			val := cur.Next.Val
			// 只要接下来的节点值等于 x，就不断跳过
			for cur.Next != nil && cur.Next.Val == val {
				cur.Next = cur.Next.Next
			}
		} else {
			// 没有重复，指针正常后移
			cur = cur.Next
		}
	}

	return dummyNode.Next

}

func TestDeleteDuplicates(t *testing.T) {
	// 链表测试通常需要辅助函数
}

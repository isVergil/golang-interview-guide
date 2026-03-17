package topinterview150

import (
	"testing"
)

// 23. 合并 K 个升序链表 (Merge k Sorted Lists)
//
// 题目描述:
// 给你一个链表数组，每个链表都已经按升序排列。
// 请你将所有链表合并到一个升序链表中，返回合并后的链表。
//
// 示例 1：
// 输入：lists = [[1,4,5],[1,3,4],[2,6]]
// 输出：[1,1,2,3,4,4,5,6]
//
// 示例 2：
// 输入：lists = []
// 输出：[]
//
// 示例 3：
// 输入：lists = [[]]
// 输出：[]

func mergeKLists(lists []*ListNode) *ListNode {
	if lists == nil || len(lists) == 0 {
		return nil
	}

	if len(lists) == 1 {
		return lists[0]
	}

	res := lists[0]
	for i := 1; i < len(lists); i++ {
		res = mergeTwoLists(res, lists[i])
	}

	return res
}

func TestMergeKLists(t *testing.T) {
	// 合并 K 个链表测试
}

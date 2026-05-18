package top100liked

import "testing"

// 160. 相交链表 (Intersection of Two Linked Lists)
//
// 题目描述:
// 给你两个单链表的头节点 headA 和 headB，请你找出并返回两个单链表相交的起始节点。
// 如果两个链表不存在相交节点，返回 nil。
// 要求时间复杂度 O(m+n)，空间复杂度 O(1)。
//
// 示例：
// 输入：listA = [4,1,8,4,5], listB = [5,6,1,8,4,5]，相交节点值为 8
// 输出：Reference of the node with value = 8
//
// 提示：双指针法，A走完走B，B走完走A，相遇点即为交点

func getIntersectionNode(headA, headB *ListNode) *ListNode {
	if headA == nil || headB == nil {
		return nil
	}
	pA, pB := headA, headB
	for pA != pB {
		if pA == nil {
			pA = headB
		} else {
			pA = pA.Next
		}
		if pB == nil {
			pB = headA
		} else {
			pB = pB.Next
		}
	}
	return pA
}

func TestGetIntersectionNode(t *testing.T) {
	// 构造相交链表
	common := &ListNode{Val: 8, Next: &ListNode{Val: 4, Next: &ListNode{Val: 5}}}
	headA := &ListNode{Val: 4, Next: &ListNode{Val: 1, Next: common}}
	headB := &ListNode{Val: 5, Next: &ListNode{Val: 6, Next: &ListNode{Val: 1, Next: common}}}

	got := getIntersectionNode(headA, headB)
	if got != common {
		t.Errorf("getIntersectionNode() = %v, want node with value 8", got)
	}

	// 不相交的情况
	listC := &ListNode{Val: 1, Next: &ListNode{Val: 2}}
	listD := &ListNode{Val: 3, Next: &ListNode{Val: 4}}
	got2 := getIntersectionNode(listC, listD)
	if got2 != nil {
		t.Errorf("getIntersectionNode() = %v, want nil", got2)
	}
}

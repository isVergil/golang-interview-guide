package topinterview150

import (
	"fmt"
	"reflect"
	"testing"
)

// 2. 两数相加 (Add Two Numbers)
//
// 题目描述:
// 给你两个 非空 的链表，表示两个非负的整数。它们每位数字都是按照 逆序 的方式存储的，并且每个节点只能存储 一位 数字。
// 请你将两个数相加，并以相同形式返回一个表示和的链表。
// 你可以假设除了数字 0 之外，这两个数都不会以 0 开头。
//
// 示例 1：
// 输入：l1 = [2,4,3], l2 = [5,6,4]
// 输出：[7,0,8]
// 解释：342 + 465 = 807.
//
// 示例 2：
// 输入：l1 = [0], l2 = [0]
// 输出：[0]
//
// 示例 3：
// 输入：l1 = [9,9,9,9,9,9,9], l2 = [9,9,9,9]
// 输出：[8,9,9,9,0,0,0,1]

/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */

func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
	head := &ListNode{}
	cur := head
	leap := 0
	for l1 != nil || l2 != nil || leap > 0 {
		curSum := leap
		if l1 != nil {
			curSum += l1.Val
			l1 = l1.Next
		}
		if l2 != nil {
			curSum += l2.Val
			l2 = l2.Next
		}
		leap = curSum / 10
		last := curSum % 10
		cur.Next = &ListNode{last, nil}
		cur = cur.Next
	}
	return head.Next
}

func TestAddTwoNumbers(t *testing.T) {
	// 辅助函数：切片转链表
	sliceToList := func(nums []int) *ListNode {
		dummy := &ListNode{}
		curr := dummy
		for _, num := range nums {
			fmt.Printf("%d,", num)
			curr.Next = &ListNode{Val: num}
			curr = curr.Next
		}
		return dummy.Next
	}
	//辅助函数：链表转切片
	listToSlice := func(head *ListNode) []int {
		res := []int{}
		for head != nil {
			res = append(res, head.Val)
			head = head.Next
		}
		return res
	}

	tests := []struct {
		name     string
		l1       []int
		l2       []int
		expected []int
	}{
		{
			name:     "Example 1",
			l1:       []int{2, 4, 3},
			l2:       []int{5, 6, 4},
			expected: []int{7, 0, 8},
		},
		{
			name:     "Example 2",
			l1:       []int{0},
			l2:       []int{0},
			expected: []int{0},
		},
		{
			name:     "Example 3",
			l1:       []int{9, 9, 9, 9, 9, 9, 9},
			l2:       []int{9, 9, 9, 9},
			expected: []int{8, 9, 9, 9, 0, 0, 0, 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: Uncomment when implementation is ready
			l1 := sliceToList(tt.l1)
			l2 := sliceToList(tt.l2)
			got := addTwoNumbers(l1, l2)
			gotSlice := listToSlice(got)
			if !reflect.DeepEqual(gotSlice, tt.expected) {
				t.Errorf("addTwoNumbers() = %v, want %v", gotSlice, tt.expected)
			}
		})
	}
}

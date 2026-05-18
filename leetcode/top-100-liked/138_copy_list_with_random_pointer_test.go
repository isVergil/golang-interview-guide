package top100liked

import (
	"reflect"
	"testing"
)

// 138. 随机链表的复制 (Copy List with Random Pointer)
//
// 题目描述:
// 给你一个长度为 n 的链表，每个节点包含一个额外增加的随机指针 random，
// 该指针可以指向链表中的任何节点或空节点。构造这个链表的深拷贝。
//
// 示例：
// 输入：head = [[7,null],[13,0],[11,4],[10,2],[1,0]]
// 输出：[[7,null],[13,0],[11,4],[10,2],[1,0]]
//
// 提示：
// 方法1：哈希表存原节点→新节点映射
// 方法2：在每个节点后面插入复制节点，再拆分

type RandomNode struct {
	Val    int
	Next   *RandomNode
	Random *RandomNode
}

func copyRandomList(head *RandomNode) *RandomNode {
	if head == nil {
		return nil
	}

	oldToNew := make(map[*RandomNode]*RandomNode)
	cur := head
	for cur != nil {
		oldToNew[cur] = &RandomNode{Val: cur.Val}
		cur = cur.Next
	}

	cur = head
	for cur != nil {
		oldToNew[cur].Next = oldToNew[cur.Next]
		oldToNew[cur].Random = oldToNew[cur.Random]
		cur = cur.Next
	}
	return oldToNew[head]
}

func TestCopyRandomList(t *testing.T) {
	// 构造 [7,null] -> [13,0] -> [11,4] -> [10,2] -> [1,0]
	n1 := &RandomNode{Val: 7}
	n2 := &RandomNode{Val: 13}
	n3 := &RandomNode{Val: 11}
	n4 := &RandomNode{Val: 10}
	n5 := &RandomNode{Val: 1}
	n1.Next = n2
	n2.Next = n3
	n3.Next = n4
	n4.Next = n5
	n2.Random = n1
	n3.Random = n5
	n4.Random = n3
	n5.Random = n1

	copied := copyRandomList(n1)

	// 验证值顺序
	vals := []int{}
	for node := copied; node != nil; node = node.Next {
		vals = append(vals, node.Val)
	}
	expected := []int{7, 13, 11, 10, 1}
	if !reflect.DeepEqual(vals, expected) {
		t.Errorf("copyRandomList values = %v, want %v", vals, expected)
	}

	// 验证是深拷贝（不是同一个节点）
	if copied == n1 {
		t.Errorf("copyRandomList returned same node, not a deep copy")
	}
}

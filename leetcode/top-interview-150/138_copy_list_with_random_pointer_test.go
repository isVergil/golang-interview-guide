package topinterview150

import (
	"testing"
)

// 138. 随机链表的复制 (Copy List with Random Pointer)
//
// 题目描述:
// 给你一个长度为 n 的链表，每个节点包含一个额外增加的随机指针 random ，该指针可以指向链表中的任何节点或空节点。
// 构造这个链表的 深拷贝。 深拷贝应该正好由 n 个 全新 节点组成，其中每个新节点的值都设为其对应的原节点的值。新节点的 next 指针和 random 指针也都应指向复制链表中的新节点，并使原链表和复制链表中的这些指针能够表示相同的链表状态。复制链表中的指针都不应指向原链表中的节点 。
// 例如，如果原链表中有 X 和 Y 两个节点，其中 X.random --> Y 。那么在复制链表中对应的两个节点 x 和 y ，同样有 x.random --> y 。
// 返回复制链表的头节点。
// 用一个由 n 个节点组成的链表来表示输入/输出中的链表。每个节点用一个 [val, random_index] 表示：
// val：一个表示 Node.val 的整数。
// random_index：随机指针指向的节点索引（范围从 0 到 n-1）；如果不指向任何节点，则为  null 。
// 你的代码 只 接受原链表的头节点 head 作为传入参数。
//
// 示例 1：
// 输入：head = [[7,null],[13,0],[11,4],[10,2],[1,0]]
// 输出：[[7,null],[13,0],[11,4],[10,2],[1,0]]
//
// 示例 2：
// 输入：head = [[1,1],[2,1]]
// 输出：[[1,1],[2,1]]
//
// 示例 3：
// 输入：head = [[3,null],[3,0],[3,null]]
// 输出：[[3,null],[3,0],[3,null]]

/**
 * Definition for a Node.
 * type Node struct {
 *     Val int
 *     Next *Node
 *     Random *Node
 * }
 */

func copyRandomList(head *Node) *Node {
	if head == nil {
		return nil
	}

	// 1 复制并交错插入
	// A -> B -> C 变成 A -> A' -> B -> B' -> C -> C'
	cur := head
	for cur != nil {
		newNode := &Node{Val: cur.Val, Next: cur.Next}
		cur.Next = newNode
		cur = newNode.Next
	}

	// 2 设置新节点的 Random 指针
	// 原 A.Random = C 则 A'.Random = A.Random.Next
	cur = head
	for cur != nil {
		if cur.Random != nil {
			// cur.Next 才是新节点
			cur.Next.Random = cur.Random.Next
		}
		cur = cur.Next.Next
	}

	// 3 拆分长链表
	newHead := head.Next
	cur = head
	for cur != nil {
		copedNode := cur.Next
		cur.Next = copedNode.Next
		if copedNode.Next != nil {
			copedNode.Next = copedNode.Next.Next
		}
		cur = cur.Next
	}
	return newHead

}

func TestCopyRandomList(t *testing.T) {
	// 测试深拷贝比较复杂，通常需要验证值、Next结构、Random结构以及内存地址不同
	// 这里仅提供框架
}

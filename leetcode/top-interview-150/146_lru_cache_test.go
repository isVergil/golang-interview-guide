package topinterview150

import (
	"testing"
)

// 146. LRU 缓存 (LRU Cache)
//
// 题目描述:
// 请你设计并实现一个满足  LRU (最近最少使用) 缓存 约束的数据结构。
// 实现 LRUCache 类：
// LRUCache(int capacity) 以 正整数 作为容量 capacity 初始化 LRU 缓存
// int get(int key) 如果关键字 key 存在于缓存中，则返回关键字的值，否则返回 -1 。
// void put(int key, int value) 如果关键字 key 已经存在，则变更其数据值 value ；如果不存在，则向缓存中插入该组 key-value 。如果插入操作导致关键字数量超过 capacity ，则应该 逐出 最久未使用的关键字。
// 函数 get 和 put 必须以 O(1) 的平均时间复杂度运行。
//
// 示例：
// 输入
// ["LRUCache", "put", "put", "get", "put", "get", "put", "get", "get", "get"]
// [[2], [1, 1], [2, 2], [1], [3, 3], [2], [4, 4], [1], [3], [4]]
// 输出
// [null, null, null, 1, null, -1, null, -1, 3, 4]
//
// 解释
// LRUCache lRUCache = new LRUCache(2);
// lRUCache.put(1, 1); // 缓存是 {1=1}
// lRUCache.put(2, 2); // 缓存是 {1=1, 2=2}
// lRUCache.get(1);    // 返回 1
// lRUCache.put(3, 3); // 该操作会使得关键字 2 作废，缓存是 {1=1, 3=3}
// lRUCache.get(2);    // 返回 -1 (未找到)
// lRUCache.put(4, 4); // 该操作会使得关键字 1 作废，缓存是 {4=4, 3=3}
// lRUCache.get(1);    // 返回 -1 (未找到)
// lRUCache.get(3);    // 返回 3
// lRUCache.get(4);    // 返回 4

// 双向链表
type lruNode struct {
	key, val   int
	prev, next *lruNode
}

type LRUCache struct {
	cap        int
	cache      map[int]*lruNode
	head, tail *lruNode
}

func ConstructorLRUCache(capacity int) LRUCache {
	l := LRUCache{
		cap:   capacity,
		cache: make(map[int]*lruNode),
		head:  &lruNode{},
		tail:  &lruNode{},
	}
	l.head.next = l.tail
	l.tail.prev = l.head
	return l
}

// get 操作：如果存在，移到头部，返回其值
func (this *LRUCache) Get(key int) int {
	if node, ok := this.cache[key]; ok {
		this.moveToHead(node)
		return node.val
	}
	return -1
}

// put 操作：如果存在则更新并移到头部；不存在则创建
func (this *LRUCache) Put(key int, value int) {
	if node, ok := this.cache[key]; ok {
		node.val = value
		this.moveToHead(node)
	} else {
		newNode := &lruNode{
			key: key,
			val: value,
		}
		this.cache[key] = newNode
		this.addToHead(newNode)
		if len(this.cache) > this.cap {
			removed := this.removeTail()
			delete(this.cache, removed.key)
		}
	}
}

// addToHead 操作：添加到 head 后面
func (this *LRUCache) addToHead(node *lruNode) {
	node.prev = this.head
	node.next = this.head.next
	this.head.next.prev = node
	this.head.next = node
}

// removeNode 操作：删除 node
func (this *LRUCache) removeNode(node *lruNode) {
	node.prev.next = node.next
	node.next.prev = node.prev
}

// moveToHead 操作：移动到 head
func (this *LRUCache) moveToHead(node *lruNode) {
	this.removeNode(node)
	this.addToHead(node)
}

// removeTail 操作：删除 tail
func (this *LRUCache) removeTail() *lruNode {
	node := this.tail.prev
	this.removeNode(node)
	return node
}

func TestLRUCache(t *testing.T) {
	// 测试 LRU 逻辑
}

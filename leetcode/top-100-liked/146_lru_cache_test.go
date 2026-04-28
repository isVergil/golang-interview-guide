package top100liked

import (
	"testing"
)

// 146. LRU 缓存 (LRU Cache)
//
// 题目描述:
// 请你设计并实现一个满足 LRU (最近最少使用) 缓存约束的数据结构。
// 实现 LRUCache 类：
// - LRUCache(int capacity) 以正整数作为容量 capacity 初始化 LRU 缓存
// - int Get(int key) 如果关键字 key 存在于缓存中，则返回关键字的值，否则返回 -1
// - void Put(int key, int value) 如果关键字 key 已经存在，则变更其数据值 value；
//   如果不存在，则向缓存中插入该组 key-value。如果插入操作导致关键字数量超过 capacity，
//   则应该逐出最久未使用的关键字。
//
// Get 和 Put 必须以 O(1) 的平均时间复杂度运行。
//
// 示例：
// 输入：
// ["LRUCache", "put", "put", "get", "put", "get", "put", "get", "get", "get"]
// [[2], [1, 1], [2, 2], [1], [3, 3], [2], [4, 4], [1], [3], [4]]
// 输出：
// [null, null, null, 1, null, -1, null, -1, 3, 4]

// 双向链表节点
type lruNode struct {
	key, val   int
	prev, next *lruNode
}

// LRUCache 哈希表 + 双向链表
// 哈希表 O(1) 查找节点，双向链表 O(1) 维护访问顺序
// head 后面是最近使用的，tail 前面是最久未使用的
type LRUCache struct {
	cap        int
	cache      map[int]*lruNode
	head, tail *lruNode // 哨兵节点，不存数据，省去边界判断
}

func NewLRUCache(capacity int) LRUCache {
	head, tail := &lruNode{}, &lruNode{}
	head.next = tail
	tail.prev = head
	return LRUCache{
		cap:   capacity,
		cache: make(map[int]*lruNode),
		head:  head,
		tail:  tail,
	}
}

// Get 查找 key，找到则提到最前面（标记为最近使用）
func (c *LRUCache) Get(key int) int {
	if n, ok := c.cache[key]; ok {
		c.remove(n)
		c.pushFront(n)
		return n.val
	}
	return -1
}

// Put 插入/更新 key-value，超容量时淘汰最久未使用的
func (c *LRUCache) Put(key int, value int) {
	// key 已存在，更新值，提到最前面
	if n, ok := c.cache[key]; ok {
		n.val = value
		c.remove(n)
		c.pushFront(n)
		return
	}

	// key 不存在，新建节点插到最前面
	n := &lruNode{key: key, val: value}
	c.cache[key] = n
	c.pushFront(n)

	// 超过容量，淘汰 tail 前面那个（最久未使用）
	if len(c.cache) > c.cap {
		last := c.tail.prev
		c.remove(last)
		delete(c.cache, last.key) // 节点里存 key 就是为了这一步能从 map 里删
	}
}

// remove 从链表中摘掉节点
func (c *LRUCache) remove(n *lruNode) {
	n.prev.next = n.next
	n.next.prev = n.prev
}

// pushFront 插到 head 后面（标记为最近使用）
func (c *LRUCache) pushFront(n *lruNode) {
	n.next = c.head.next
	n.prev = c.head
	c.head.next.prev = n
	c.head.next = n
}

func TestLRUCache(t *testing.T) {
	cache := NewLRUCache(2)

	cache.Put(1, 1)
	cache.Put(2, 2)
	if got := cache.Get(1); got != 1 {
		t.Errorf("Get(1) = %d, want 1", got)
	}

	cache.Put(3, 3) // 淘汰 key=2
	if got := cache.Get(2); got != -1 {
		t.Errorf("Get(2) = %d, want -1", got)
	}

	cache.Put(4, 4) // 淘汰 key=1
	if got := cache.Get(1); got != -1 {
		t.Errorf("Get(1) = %d, want -1", got)
	}
	if got := cache.Get(3); got != 3 {
		t.Errorf("Get(3) = %d, want 3", got)
	}
	if got := cache.Get(4); got != 4 {
		t.Errorf("Get(4) = %d, want 4", got)
	}
}

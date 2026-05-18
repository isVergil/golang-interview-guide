package top100liked

import "testing"

// 208. 实现 Trie 前缀树 (Implement Trie)
//
// 题目描述:
// 实现一个 Trie（前缀树），包含 insert、search 和 startsWith 三个操作。
// - insert(word) 向前缀树中插入字符串 word
// - search(word) 如果字符串 word 在前缀树中，返回 true
// - startsWith(prefix) 如果之前插入的字符串有前缀 prefix，返回 true
//
// 示例：
// 输入：["Trie","insert","search","search","startsWith","insert","search"]
//       [[],["apple"],["apple"],["app"],["app"],["app"],["app"]]
// 输出：[null,null,true,false,true,null,true]
//
// 提示：每个节点 26 个子节点指针 + isEnd 标记

type Trie struct {
	children [26]*Trie
	isEnd    bool
}

func NewTrie() Trie {
	return Trie{}
}

func (t *Trie) Insert(word string) {
	node := t
	for i := 0; i < len(word); i++ {
		idx := word[i] - 'a'
		if node.children[idx] == nil {
			node.children[idx] = &Trie{}
		}
		node = node.children[idx]
	}
	node.isEnd = true
}

func (t *Trie) Search(word string) bool {
	node := t.findNode(word)
	return node != nil && node.isEnd
}

func (t *Trie) StartsWith(prefix string) bool {
	return t.findNode(prefix) != nil
}

func (t *Trie) findNode(s string) *Trie {
	node := t
	for i := 0; i < len(s); i++ {
		idx := s[i] - 'a'
		if node.children[idx] == nil {
			return nil
		}
		node = node.children[idx]
	}
	return node
}

func TestTrie(t *testing.T) {
	trie := NewTrie()
	trie.Insert("apple")

	if got := trie.Search("apple"); !got {
		t.Errorf("Search(apple) = false, want true")
	}
	if got := trie.Search("app"); got {
		t.Errorf("Search(app) = true, want false")
	}
	if got := trie.StartsWith("app"); !got {
		t.Errorf("StartsWith(app) = false, want true")
	}

	trie.Insert("app")
	if got := trie.Search("app"); !got {
		t.Errorf("Search(app) = false, want true")
	}
}

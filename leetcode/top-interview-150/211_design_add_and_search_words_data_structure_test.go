package topinterview150

import (
	"testing"
)

// 211. 添加与搜索单词 - 数据结构设计 (Design Add and Search Words Data Structure)
//
// 题目描述:
// 请你设计一个数据结构，支持 添加新单词 和 查找字符串是否与任何先前添加的单词匹配 。
// 实现 WordDictionary 类：
// WordDictionary() 初始化词典对象
// void addWord(word) 将 word 添加到数据结构中，之后可以对它进行匹配
// bool search(word) 如果数据结构中存在字符串与 word 匹配，则返回 true ；否则，返回  false 。word 中可能包含一些 '.' ，每个 . 都可以表示任何一个字母。
//
// 示例：
// 输入：
// ["WordDictionary","addWord","addWord","addWord","search","search","search","search"]
// [[],["bad"],["dad"],["mad"],["pad"],["bad"],[".ad"],["b.."]]
// 输出：
// [null,null,null,null,false,true,true,true]

type WordDictionary struct {
	children [26]*WordDictionary
	isEnd    bool
}

func ConstructorWordDictionary() WordDictionary {
	return WordDictionary{}
}

func (this *WordDictionary) AddWord(word string) {
	node := this
	for i := 0; i < len(word); i++ {
		idx := word[i] - 'a'
		if node.children[idx] == nil {
			node.children[idx] = &WordDictionary{}
		}
		node = node.children[idx]
	}
	node.isEnd = true
}

func (this *WordDictionary) Search(word string) bool {
	return this.dfs(word, this)
}

func (this *WordDictionary) dfs(word string, node *WordDictionary) bool {
	for i := 0; i < len(word); i++ {
		ch := word[i]

		if ch == '.' {
			// 性能优化点：遇到通配符，遍历所有可能的子路径
			for _, child := range node.children {
				if child != nil && this.dfs(word[i+1:], child) {
					return true
				}
			}
			return false // 所有路径都走不通
		}

		// 普通字母处理
		idx := ch - 'a'
		if node.children[idx] == nil {
			return false
		}
		node = node.children[idx]
	}
	return node.isEnd
}

func TestWordDictionary(t *testing.T) {
	// 测试 Trie 逻辑
}

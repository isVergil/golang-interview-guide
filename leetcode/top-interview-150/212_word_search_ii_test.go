package topinterview150

import (
	"testing"
)

// 212. 单词搜索 II (Word Search II)
//
// 题目描述:
// 给定一个 m x n 二维字符网格 board 和一个单词（字符串）列表 words， 返回所有二维网格上的单词 。
// 单词必须按照字母顺序，通过相邻的单元格内的字母构成，其中“相邻”单元格是那些水平相邻或垂直相邻的单元格。同一个单元格内的字母在一个单词中不允许被重复使用。
//
// 示例 1：
// 输入：board = [["o","a","a","n"],["e","t","a","e"],["i","h","k","r"],["i","f","l","v"]], words = ["oath","pea","eat","rain"]
// 输出：["eat","oath"]
//
// 示例 2：
// 输入：board = [["a","b"],["c","d"]], words = ["abcb"]
// 输出：[]

type TrieNode struct {
	children [26]*TrieNode
	word     string // 优化点：直接存单词字符串，找到时无需再拼接
}

func findWords(board [][]byte, words []string) []string {
	// 1. 构建 Trie
	root := &TrieNode{}
	for _, w := range words {
		node := root
		for i := 0; i < len(w); i++ {
			idx := w[i] - 'a'
			if node.children[idx] == nil {
				node.children[idx] = &TrieNode{}
			}
			node = node.children[idx]
		}
		node.word = w
	}

	res := make([]string, 0)
	rows, cols := len(board), len(board[0])

	// 2. DFS 回溯
	var dfs func(r, c int, node *TrieNode)
	dfs = func(r, c int, node *TrieNode) {
		char := board[r][c]
		idx := char - 'a'
		child := node.children[idx]

		// 剪枝：字符不匹配或已访问过
		if child == nil {
			return
		}

		// 命中单词
		if child.word != "" {
			res = append(res, child.word)
			child.word = "" // 关键优化：去重，防止同一个单词被多次添加
		}

		// 标记已访问（原地修改，避免额外开销）
		board[r][c] = '#'

		// 四方向探索
		dirs := [4][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
		for _, d := range dirs {
			nr, nc := r+d[0], c+d[1]
			if nr >= 0 && nr < rows && nc >= 0 && nc < cols && board[nr][nc] != '#' {
				dfs(nr, nc, child)
			}
		}

		// 恢复现场（回溯）
		board[r][c] = char

		// 进阶优化：如果该节点已无子节点，将其从父节点中删除（剪枝掉已搜完的路径）
		// if isEmpty(child) { node.children[idx] = nil }
	}

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			dfs(r, c, root)
		}
	}

	return res
}

func TestFindWords(t *testing.T) {
	// 测试单词搜索逻辑
}

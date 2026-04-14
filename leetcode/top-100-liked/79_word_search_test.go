package top100liked

import (
	"testing"
)

// 79. 单词搜索 (Word Search)
//
// 题目描述:
// 给定一个 m x n 二维字符网格 board 和一个字符串单词 word 。
// 如果 word 存在于网格中，返回 true ；否则，返回 false 。
// 单词必须按照字母顺序，通过相邻的单元格内的字母构成，
// 其中"相邻"单元格是那些水平相邻或垂直相邻的单元格。同一个单元格内的字母不允许被重复使用。
//
// 示例 1：
// 输入：board = [["A","B","C","E"],["S","F","C","S"],["A","D","E","E"]], word = "ABCCED"
// 输出：true
//
// 示例 2：
// 输入：board = [["A","B","C","E"],["S","F","C","S"],["A","D","E","E"]], word = "SEE"
// 输出：true
//
// 示例 3：
// 输入：board = [["A","B","C","E"],["S","F","C","S"],["A","D","E","E"]], word = "ABCB"
// 输出：false

func exist(board [][]byte, word string) bool {
	m, n := len(board), len(board[0])

	// 优化1: 频次剪枝 — 字符不够直接返回
	var freq [128]int
	for _, row := range board {
		for _, c := range row {
			freq[c]++
		}
	}
	for _, c := range word {
		freq[c]--
		if freq[c] < 0 {
			return false
		}
	}

	// 优化2: 反转搜索 — 从稀少端开始，减少 DFS 入口和分支
	// 比如 word="AAAAAB"，棋盘里 A 很多 B 只有1个
	// 正向搜索入口多、分支多；反转后从 B 开始，入口只有1个
	if freq[word[0]] > freq[word[len(word)-1]] {
		b := []byte(word)
		for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
			b[i], b[j] = b[j], b[i]
		}
		word = string(b)
	}

	var backtrack func(int, int, int) bool
	backtrack = func(idx, i, j int) bool {
		if idx == len(word) {
			return true
		}
		if i < 0 || i >= m || j < 0 || j >= n || board[i][j] != word[idx] {
			return false
		}

		tmp := board[i][j]
		board[i][j] = '#'

		found := backtrack(idx+1, i-1, j) ||
			backtrack(idx+1, i+1, j) ||
			backtrack(idx+1, i, j-1) ||
			backtrack(idx+1, i, j+1)

		board[i][j] = tmp
		return found
	}

	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if backtrack(0, i, j) {
				return true
			}
		}
	}
	return false
}

func TestExist(t *testing.T) {
	board := [][]byte{
		{'A', 'B', 'C', 'E'},
		{'S', 'F', 'C', 'S'},
		{'A', 'D', 'E', 'E'},
	}

	tests := []struct {
		name     string
		board    [][]byte
		word     string
		expected bool
	}{
		{
			name:     "示例1-ABCCED",
			board:    board,
			word:     "ABCCED",
			expected: true,
		},
		{
			name:     "示例2-SEE",
			board:    board,
			word:     "SEE",
			expected: true,
		},
		{
			name:     "示例3-ABCB",
			board:    board,
			word:     "ABCB",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 每次测试需要复制 board，避免上一次测试修改影响
			boardCopy := make([][]byte, len(tt.board))
			for i := range tt.board {
				boardCopy[i] = make([]byte, len(tt.board[i]))
				copy(boardCopy[i], tt.board[i])
			}
			got := exist(boardCopy, tt.word)
			if got != tt.expected {
				t.Errorf("exist() = %v, want %v", got, tt.expected)
			}
		})
	}
}

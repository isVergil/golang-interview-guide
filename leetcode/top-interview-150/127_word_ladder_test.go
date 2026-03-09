package topinterview150

import (
	"testing"
)

// 127. 单词接龙 (Word Ladder)
//
// 题目描述:
// 字典 wordList 中从单词 beginWord 到 endWord 的 转换序列 是一个按下述规格形成的序列 beginWord -> s1 -> s2 -> ... -> sk：
// 每一对相邻的单词之间仅恰好相差一个字符。
// 对于 1 <= i <= k ，每个 si 都在 wordList 中。注意， beginWord 不必在 wordList 中。
// sk == endWord
// 给你两个单词 beginWord 和 endWord 和一个字典 wordList ，返回 从 beginWord 到 endWord 的 最短转换序列 中的 单词数目 。如果不存在这样的转换序列，返回 0 。
//
// 示例 1：
// 输入：beginWord = "hit", endWord = "cog", wordList = ["hot","dot","dog","lot","log","cog"]
// 输出：5
// 解释：一个最短转换序列是 "hit" -> "hot" -> "dot" -> "dog" -> "cog", 返回它的长度 5。
//
// 示例 2：
// 输入：beginWord = "hit", endWord = "cog", wordList = ["hot","dot","dog","lot","log"]
// 输出：0
// 解释：endWord "cog" 不在字典中，所以无法进行转换。

func ladderLength(beginWord string, endWord string, wordList []string) int {
	panic("not implemented")
}

func TestLadderLength(t *testing.T) {
	// 单词接龙测试
}

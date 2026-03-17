package topinterview150

import (
	"testing"
)

// 30. 串联所有单词的子串 (Substring with Concatenation of All Words)
//
// 题目描述:
// 给定一个字符串 s 和一个字符串数组 words。 words 中所有字符串 长度相同。
//  s 中的 串联子串 是指一个包含  words 中所有字符串以任意顺序排列连接起来的子串。
// 返回所有串联子串在 s 中的开始索引。你可以按 任意顺序 返回答案。
//
// 示例 1：
// 输入：s = "barfoothefoobarman", words = ["foo","bar"]
// 输出：[0,9]
//
// 示例 2：
// 输入：s = "wordgoodgoodgoodword", words = ["word","good","best","word"]
// 输出：[]

func findSubstring(s string, words []string) []int {
	if len(words) == 0 || len(s) == 0 {
		return nil
	}

	wordLen, wordNum, sLen := len(words[0]), len(words), len(s)
	// 统计目标单词频率 预设 capacity 减少 map 扩容
	counts := make(map[string]int, wordNum)
	for _, w := range words {
		counts[w]++
	}

	res := make([]int, 0)
	// 分组滑动窗口优化：起点只需从 0 到 wordLen - 1
	for i := 0; i < wordLen; i++ {
		left, right, count := i, i, 0
		// 记录当前窗口内匹配的单词总数
		currCounts := make(map[string]int)

		// 窗口右边界每次移动一个单词的长度
		for right+wordLen <= sLen {
			w := s[right : right+wordLen]
			right += wordLen

			if targetNum, ok := counts[w]; ok {
				currCounts[w]++
				count++
				// 如果当前单词数量超标，左边界收缩
				for currCounts[w] > targetNum {
					leftW := s[left : left+wordLen]
					currCounts[leftW]--
					count--
					left += wordLen
				}
				// 达到目标数量，记录起始索引
				if count == wordNum {
					res = append(res, left)
				}
			} else {
				// 遇到完全不在词典里的词，清空窗口，重新开始
				currCounts = make(map[string]int)
				count, left = 0, right
			}
		}
	}
	return res
}

func TestFindSubstring(t *testing.T) {
	// 串联子串测试
}

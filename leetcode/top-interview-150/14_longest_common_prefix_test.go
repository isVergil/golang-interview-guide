package topinterview150

import (
	"testing"
)

// 14. 最长公共前缀 (Longest Common Prefix)
//
// 题目描述:
// 编写一个函数来查找字符串数组中的最长公共前缀。
// 如果不存在公共前缀，返回空字符串 ""。
//
// 示例 1：
// 输入：strs = ["flower","flow","flight"]
// 输出："fl"

func longestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}

	if len(strs) == 1 {
		return strs[0]
	}

	// 以第一个字符串为基准进行纵向扫描
	baseStr := strs[0]
	for i := 0; i < len(baseStr); i++ {
		char := baseStr[i]
		// 检查后续每一个字符串的第 i 位
		for j := 1; j < len(strs); j++ {
			// 细节优化：只要当前索引 i 超过了某个字符串的长度，
			// 或者发现字符不匹配，立即截取并返回
			if i == len(strs[j]) || strs[j][i] != char {
				return baseStr[:i]
			}
		}
	}

	// 如果完整走完了第一个字符串，说明第一个字符串本身就是公共前缀
	return baseStr
}

func TestLongestCommonPrefix(t *testing.T) {
	// 最长公共前缀测试
}

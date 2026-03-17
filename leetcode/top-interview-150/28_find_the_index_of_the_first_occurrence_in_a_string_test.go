package topinterview150

import (
	"testing"
)

// 28. 找出字符串中第一个匹配项的下标 (Find the Index of the First Occurrence in a String)
//
// 题目描述:
// 给你两个字符串 haystack 和 needle ，请你在 haystack 字符串中找出 needle 字符串的第一个匹配项的下标（下标从 0 开始）。如果 needle 不是 haystack 的一部分，则返回 -1 。
//
// 示例 1：
// 输入：haystack = "sadbutsad", needle = "sad"
// 输出：0
// 解释："sad" 在下标 0 和 6 处匹配。第一个匹配项的下标是 0 。
//
// 示例 2：
// 输入：haystack = "leetcode", needle = "leeto"
// 输出：-1
// 解释："leeto" 没有在 "leetcode" 中出现，所以返回 -1 。

func strStr(haystack string, needle string) int {
	n, m := len(haystack), len(needle)
	if m == 0 {
		return 0
	}

	// 预处理 needle，计算 next 数组
	// j 指向前缀的末尾（同时也代表了当前相等前后缀的长度）
	// i 指向后缀的末尾
	// next[0] 永远是 0，因为一个字符没有前后缀
	next := make([]int, m)
	for i, j := 1, 0; i < m; i++ {
		// 【细节 1：不匹配时的回退】
		// 如果当前字符对不上，j 就像在主串匹配时一样，向后跳跃
		// 这是一个“套娃”逻辑：找一个更短的相等前后缀
		for j > 0 && needle[i] != needle[j] {
			j = next[j-1]
		}

		// 【细节 2：匹配成功】
		// 如果字符相等，说明相等前后缀的长度增加 1
		if needle[i] == needle[j] {
			j++
		}

		// 【细节 3：更新数组】
		// 将当前位置 i 对应的最长相等前后缀长度记入数组
		next[i] = j
	}

	// 开始匹配
	for i, j := 0, 0; i < n; i++ {
		// 当字符不匹配，j 寻找 needle 的前缀跳跃点
		for j > 0 && haystack[i] != needle[j] {
			j = next[j-1]
		}
		if haystack[i] == needle[j] {
			j++
		}

		// 如果 j 走到了 needle 的末尾，说明匹配成功
		if j == m {
			return i - m + 1
		}
	}
	return -1
}

func TestStrStr(t *testing.T) {
	// 字符串匹配测试
}

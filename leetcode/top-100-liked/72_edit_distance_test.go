package top100liked

import (
	"testing"
)

// 72. 编辑距离 (Edit Distance)
//
// 题目描述:
// 给你两个单词 word1 和 word2，请返回将 word1 转换成 word2 所使用的最少操作数。
// 你可以对一个单词进行如下三种操作：插入一个字符、删除一个字符、替换一个字符。
//
// 示例 1：
// 输入：word1 = "horse", word2 = "ros"
// 输出：3（horse → rorse → rose → ros）
//
// 示例 2：
// 输入：word1 = "intention", word2 = "execution"
// 输出：5
// 二维 dp
func minDistance(word1 string, word2 string) int {
	m, n := len(word1), len(word2)

	// dp[i][j] 表示 word1 前 i 个字符转换成 word2 前 j 个字符需要的最少操作数
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}

	// base case: 空串到长度为i/j的串需要i/j次插入
	for i := 0; i <= m; i++ {
		dp[i][0] = i
	}
	for j := 0; j <= n; j++ {
		dp[0][j] = j
	}

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if word1[i-1] == word2[j-1] {
				dp[i][j] = dp[i-1][j-1]
			} else {
				dp[i][j] = min(dp[i-1][j-1], min(dp[i-1][j], dp[i][j-1])) + 1
			}
		}
	}
	return dp[m][n]
}

// 空间优化版
func minDistance1(word1 string, word2 string) int {
	m, n := len(word1), len(word2)

	prev := make([]int, n+1)
	curr := make([]int, n+1)

	// base case: 空串到长度为i/j的串需要i/j次插入
	for i := 0; i <= n; i++ {
		prev[i] = i
	}
	for i := 1; i <= m; i++ {
		curr[0] = i
		for j := 1; j <= n; j++ {
			if word1[i-1] == word2[j-1] {
				curr[j] = prev[j-1]
			} else {
				curr[j] = min(
					prev[j-1],
					min(
						prev[j],
						curr[j-1],
					),
				) + 1
			}
		}
		prev, curr = curr, prev
	}
	return prev[n]
}

func TestMinDistance(t *testing.T) {
	tests := []struct {
		name         string
		word1, word2 string
		expected     int
	}{
		{name: "示例1", word1: "horse", word2: "ros", expected: 3},
		{name: "示例2", word1: "intention", word2: "execution", expected: 5},
		{name: "空串到空串", word1: "", word2: "", expected: 0},
		{name: "空串到非空", word1: "", word2: "abc", expected: 3},
		{name: "相同字符串", word1: "abc", word2: "abc", expected: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := minDistance(tt.word1, tt.word2)
			if got != tt.expected {
				t.Errorf("minDistance() = %v, want %v", got, tt.expected)
			}
		})

		t.Run(tt.name, func(t *testing.T) {
			got := minDistance1(tt.word1, tt.word2)
			if got != tt.expected {
				t.Errorf("minDistance1() = %v, want %v", got, tt.expected)
			}
		})
	}
}

package topinterview150

import (
	"testing"
)

// 205. 同构字符串 (Isomorphic Strings)
//
// 题目描述:
// 给定两个字符串 s 和 t ，判断它们是否是同构的。
// 如果 s 中的字符可以按某种映射关系替换得到 t ，那么这两个字符串是同构的。
// 每个出现的字符都应当映射到另一个字符，同时不改变字符的顺序。不同字符不能映射到同一个字符上，相同字符只能映射到同一个字符上，字符可以映射到自己本身。
//
// 示例 1:
// 输入：s = "egg", t = "add"
// 输出：true
//
// 示例 2：
// 输入：s = "foo", t = "bar"
// 输出：false
//
// 示例 3：
// 输入：s = "paper", t = "title"
// 输出：true

func isIsomorphic(s string, t string) bool {
	sPos, tPos := [256]int{}, [256]int{}

	for i := 0; i < len(s); i++ {
		// 如果它们上一次出现的位置不一样，说明不同构
		if sPos[s[i]] != tPos[t[i]] {
			return false
		}
		sPos[s[i]], tPos[t[i]] = i+1, i+1
	}
	return true
}

func TestIsIsomorphic(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		t        string
		expected bool
	}{
		{"Example 1", "egg", "add", true},
		{"Example 2", "foo", "bar", false},
		{"Example 3", "paper", "title", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := isIsomorphic(tt.s, tt.t); got != tt.expected {
			// 	t.Errorf("isIsomorphic() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

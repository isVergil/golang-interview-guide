package topinterview150

import "testing"

// 242. 有效的字母异位词 (Valid Anagram)
//
// 题目描述:
// 给定两个字符串 s 和 t ，编写一个函数来判断 t 是否是 s 的字母异位词。
// 注意：若 s 和 t 中每个字符出现的次数都相同，则称 s 和 t 互为字母异位词。
//
// 示例 1:
// 输入: s = "anagram", t = "nagaram"
// 输出: true
//
// 示例 2:
// 输入: s = "rat", t = "car"
// 输出: false

func isAnagram(s string, t string) bool {
	// 长度不等，肯定不是异位词
	if len(s) != len(t) {
		return false
	}

	// 因为题目假设只有小写字母，用长度 26 的数组代替 Map
	// 这种做法叫“简易哈希表”，性能比 map 高得多
	var record [26]int

	// 遍历 s，对应位置 +1
	// 'a' 的 ASCII 码是 97，s[i]-'a' 就能把字母映射到 0-25 的索引上
	for i := 0; i < len(s); i++ {
		record[s[i]-'a']++
	}

	// 遍历 t，对应位置 -1
	for i := 0; i < len(t); i++ {
		record[t[i]-'a']--
		// 如果减完后小于 0，说明 t 里的这个字母比 s 多，直接出局
		if record[t[i]-'a'] < 0 {
			return false
		}
	}

	return true
}

// 数组可以比较 直接对比
func isAnagram1(s string, t string) bool {
	record := [26]int{}
	for _, r := range s {
		record[r-rune('a')]++
	}
	for _, r := range t {
		record[r-rune('a')]--
		if record[r-rune('a')] < 0 {
			return false
		}
	}
	return record == [26]int{}
}

func TestIsAnagram(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		t        string
		expected bool
	}{
		{
			name:     "Example 1",
			s:        "anagram",
			t:        "nagaram",
			expected: true,
		},
		{
			name:     "Example 2",
			s:        "rat",
			t:        "car",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: Uncomment when implementation is ready
			// if got := isAnagram(tt.s, tt.t); got != tt.expected {
			// 	t.Errorf("isAnagram() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

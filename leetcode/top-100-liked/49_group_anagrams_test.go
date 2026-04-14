package top100liked

import (
	"reflect"
	"testing"
)

// 49. 字母异位词分组 (Group Anagrams)
//
// 题目描述:
// 给你一个字符串数组，请你将字母异位词组合在一起。可以按任意顺序返回结果列表。
// 字母异位词是由重新排列源单词的所有字母得到的一个新单词。
//
// 示例 1：
// 输入: strs = ["eat","tea","tan","ate","nat","bat"]
// 输出: [["bat"],["nat","tan"],["ate","eat","tea"]]
//
// 示例 2：
// 输入: strs = [""]
// 输出: [[""]]
//
// 示例 3：
// 输入: strs = ["a"]
// 输出: [["a"]]

func groupAnagrams(strs []string) [][]string {
	groups := make(map[[26]byte][]string)

	for _, s := range strs {
		// Q1: [26]byte 是数组(array)，是值类型，可比较，可以做 map key
		//     []byte 是切片(slice)，是引用类型，不可比较，不能做 map key
		//     Go 中数组和切片是两种不同类型：
		//       [N]T  → 值类型，赋值/传参会拷贝，支持 == 比较
		//       []T   → 引用类型，底层是 (ptr, len, cap)，不支持 == 比较
		var key [26]byte

		// Q2: range string 得到的 value 是 rune(int32)，用来做数组下标没问题
		//     因为 rune 和 byte 运算时会自动隐式转换为相同类型(都是整数类型)
		//     但对于纯 ASCII 场景，用 s[i] 直接取 byte 更高效，避免 rune 解码开销
		for i := 0; i < len(s); i++ {
			key[s[i]-'a']++
		}

		groups[key] = append(groups[key], s)
	}

	res := make([][]string, 0, len(groups))
	for _, g := range groups {
		res = append(res, g)
	}
	return res
}

func TestGroupAnagrams(t *testing.T) {
	tests := []struct {
		name     string
		strs     []string
		expected [][]string
	}{
		{
			name:     "示例1",
			strs:     []string{"eat", "tea", "tan", "ate", "nat", "bat"},
			expected: [][]string{{"ate", "eat", "tea"}, {"bat"}, {"nat", "tan"}},
		},
		{
			name:     "示例2",
			strs:     []string{""},
			expected: [][]string{{""}},
		},
		{
			name:     "示例3",
			strs:     []string{"a"},
			expected: [][]string{{"a"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := groupAnagrams(tt.strs)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("groupAnagrams() = %v, want %v", got, tt.expected)
			}
		})
	}
}

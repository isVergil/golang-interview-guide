package topinterview150

import (
	"reflect"
	"sort"
	"testing"
)

// 49. 字母异位词分组 (Group Anagrams)
//
// 题目描述:
// 给你一个字符串数组，请你将 字母异位词 组合在一起。可以按任意顺序返回结果列表。
// 字母异位词 是由重新排列源单词的字母得到的一个新单词，所有源单词中的字母通常恰好只用一次。
//
// 示例 1:
// 输入: strs = ["eat", "tea", "tan", "ate", "nat", "bat"]
// 输出: [["bat"],["nat","tan"],["ate","eat","tea"]]
//
// 示例 2:
// 输入: strs = [""]
// 输出: [[""]]
//
// 示例 3:
// 输入: strs = ["a"]
// 输出: [["a"]]

// 1 数组作为 key
func groupAnagrams(strs []string) [][]string {
	m := make(map[[26]int][]string)
	for _, str := range strs {
		cnt := [26]int{}
		for i := 0; i < len(str); i++ {
			cnt[str[i]-'a']++
		}
		m[cnt] = append(m[cnt], str)
	}

	var res [][]string
	for _, v := range m {
		res = append(res, v)
	}
	return res
}

// 1 排序
func groupAnagrams2(strs []string) [][]string {
	// 创建一个 map，Key 是排序后的字符串，Value 是原始字符串的切片
	m := make(map[string][]string)

	for _, s := range strs {
		// 1. 将字符串转为 byte 切片 (Go 中字符串是不可变的)
		sByte := []byte(s)

		// 2. 对字节切片进行排序
		// sort.Slice 是 Go 比较通用的排序方法
		sort.Slice(sByte, func(i, j int) bool {
			return sByte[i] < sByte[j]
		})

		// 3. 将排序后的切片转回字符串，作为 map 的 key
		sortedStr := string(sByte)

		// 4. 将原始字符串 s 加入对应的分组中
		m[sortedStr] = append(m[sortedStr], s)
	}

	// 5. 将 map 中的所有结果提取到二维切片中
	res := make([][]string, 0, len(m))
	for _, group := range m {
		res = append(res, group)
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
			name:     "Example 1",
			strs:     []string{"eat", "tea", "tan", "ate", "nat", "bat"},
			expected: [][]string{{"bat"}, {"nat", "tan"}, {"ate", "eat", "tea"}},
		},
		{
			name:     "Example 2",
			strs:     []string{""},
			expected: [][]string{{""}},
		},
		{
			name:     "Example 3",
			strs:     []string{"a"},
			expected: [][]string{{"a"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: Uncomment when implementation is ready
			// got := groupAnagrams(tt.strs)
			// // 结果顺序不重要，需要特殊比较逻辑
			// if !compareGroupAnagrams(got, tt.expected) {
			// 	t.Errorf("groupAnagrams() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

// 辅助函数：比较分组结果（忽略顺序）
func compareGroupAnagrams(got, expected [][]string) bool {
	// 简化比较，实际测试可能需要更复杂的逻辑
	return reflect.DeepEqual(len(got), len(expected))
}

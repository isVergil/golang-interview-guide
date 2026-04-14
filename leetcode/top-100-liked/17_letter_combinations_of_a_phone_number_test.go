package top100liked

import (
	"testing"
)

// 17. 电话号码的字母组合 (Letter Combinations of a Phone Number)
//
// 题目描述:
// 给定一个仅包含数字 2-9 的字符串，返回所有它能表示的字母组合。答案可以按任意顺序返回。
// 给出数字到字母的映射如下（与电话按键相同）。注意 1 不对应任何字母。
//
// 2 -> abc, 3 -> def, 4 -> ghi, 5 -> jkl
// 6 -> mno, 7 -> pqrs, 8 -> tuv, 9 -> wxyz
//
// 示例 1：
// 输入：digits = "23"
// 输出：["ad","ae","af","bd","be","bf","cd","ce","cf"]
//
// 示例 2：
// 输入：digits = ""
// 输出：[]
//
// 示例 3：
// 输入：digits = "2"
// 输出：["a","b","c"]
func letterCombinations(digits string) []string {
	// 特殊情况处理
	if len(digits) == 0 {
		return []string{}
	}

	var phoneMap = []string{
		"abc",
		"def",
		"ghi",
		"jkl",
		"mno",
		"pqrs",
		"tuv",
		"wxyz",
	}

	var res []string
	var path []byte // 用 byte 切片效率比 string 高

	var backTrack func(int)
	backTrack = func(idx int) {
		if idx == len(digits) {
			res = append(res, string(path))
			return
		}

		digit := digits[idx] - '2'
		letters := phoneMap[digit]

		for i := 0; i < len(letters); i++ {
			path = append(path, letters[i])

			backTrack(idx + 1)

			path = path[:len(path)-1]
		}
	}

	backTrack(0)

	return res
}

func TestLetterCombinations(t *testing.T) {
	tests := []struct {
		name        string
		digits      string
		expectedLen int
	}{
		{
			name:        "示例1",
			digits:      "23",
			expectedLen: 9,
		},
		{
			name:        "示例2-空字符串",
			digits:      "",
			expectedLen: 0,
		},
		{
			name:        "示例3-单个数字",
			digits:      "2",
			expectedLen: 3,
		},
		{
			name:        "包含7或9",
			digits:      "79",
			expectedLen: 16,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := letterCombinations(tt.digits)
			if len(got) != tt.expectedLen {
				t.Errorf("letterCombinations() returned %d results, want %d", len(got), tt.expectedLen)
			}
		})
	}
}

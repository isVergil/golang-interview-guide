package topinterview150

import (
	"testing"
)

// 17. 电话号码的字母组合 (Letter Combinations of a Phone Number)
//
// 题目描述:
// 给定一个仅包含数字 2-9 的字符串，返回所有它能表示的字母组合。答案可以按 任意顺序 返回。
// 给出数字到字母的映射如下（与电话按键相同）。注意 1 不对应任何字母。
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

var phoneMap = []string{
	"",
	"",
	"abc",
	"def",
	"ghi",
	"jkl",
	"mno",
	"pqrs",
	"tuv",
	"wxyz",
}

func letterCombinations(digits string) []string {
	if len(digits) == 0 {
		return nil
	}

	res := make([]string, 0)

	// 使用 byte 切片作为缓存，避免递归中产生大量临时字符串
	path := make([]byte, len(digits))

	var backtrack func(int)
	backtrack = func(idx int) {
		// 终止条件：路径长度达到 digits 长度
		if idx == len(digits) {
			res = append(res, string(path))
			return
		}

		// 获取当前数字对应的字母列表
		letters := phoneMap[digits[idx]-'0']
		for i := 0; i < len(letters); i++ {
			path[idx] = letters[i] // 选择当前字母
			backtrack(idx + 1)     // 递归下一层
		}
	}

	backtrack(0)
	return res

}

func TestLetterCombinations(t *testing.T) {
	// 电话号码字母组合测试
}

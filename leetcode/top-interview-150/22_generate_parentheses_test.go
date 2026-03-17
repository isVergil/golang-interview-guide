package topinterview150

import (
	"testing"
)

// 22. 括号生成 (Generate Parentheses)
//
// 题目描述:
// 数字 n 代表生成括号的对数，请你设计一个函数，用于能够生成所有可能的并且 有效的 括号组合。
//
// 示例 1：
// 输入：n = 3
// 输出：["((()))","(()())","(())()","()(())","()()()"]
//
// 示例 2：
// 输入：n = 1
// 输出：["()"]
// 回溯
func generateParenthesis(n int) []string {
	res := make([]string, 0)
	// 结果字符串的总长度固定为 2*n
	path := make([]byte, 2*n)

	// leftCount: 已使用的左括号数
	// rightCount: 已使用的右括号数
	var backtrack func(idx, leftCount, rCount int)
	backtrack = func(idx, leftCount, rCount int) {
		// 终止条件：填满了 2*n 个位置
		if idx == 2*n {
			res = append(res, string(path))
			return
		}

		// 决策 1：尝试放左括号
		// 只要左括号还没用完就可以放
		if leftCount < n {
			path[idx] = '('
			backtrack(idx+1, leftCount+1, rCount)
		}

		// 决策 2：尝试放右括号
		// 只有当前右括号数小于左括号数时，放右括号才合法
		if rCount < leftCount {
			path[idx] = ')'
			backtrack(idx+1, leftCount, rCount+1)
		}
	}

	backtrack(0, 0, 0)
	return res
}

func TestGenerateParenthesis(t *testing.T) {
	// 括号生成测试
}

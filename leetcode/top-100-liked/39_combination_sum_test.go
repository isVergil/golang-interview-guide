package top100liked

import (
	"sort"
	"testing"
)

// 39. 组合总和 (Combination Sum)
//
// 题目描述:
// 给你一个无重复元素的整数数组 candidates 和一个目标整数 target ，
// 找出 candidates 中可以使数字和为目标数 target 的所有不同组合，
// 并以列表形式返回。你可以按任意顺序返回这些组合。
// candidates 中的同一个数字可以无限制重复被选取。如果至少一个数字的被选数量不同，则两种组合是不同的。
//
// 示例 1：
// 输入：candidates = [2,3,6,7], target = 7
// 输出：[[2,2,3],[7]]
// 解释：2 和 3 可以形成一组候选，2 + 2 + 3 = 7 。注意 2 可以使用多次。7 也是一个候选，7 = 7 。仅有这两种组合。
//
// 示例 2：
// 输入：candidates = [2,3,5], target = 8
// 输出：[[2,2,2,2],[2,3,3],[3,5]]
//
// 示例 3：
// 输入：candidates = [2], target = 1
// 输出：[]

func combinationSum(candidates []int, target int) [][]int {
	var res [][]int
	var path []int

	sort.Ints(candidates)
	var backtrack func(sum int, start int)
	backtrack = func(sum int, start int) {
		if sum == target {
			// 方式 A
			// tmp := make([]int, len(path))
			// copy(tmp, path)
			// 方式 B (更简洁)
			res = append(res, append(make([]int, 0), path...))
			return
		}
		for i := start; i < len(candidates); i++ {
			if sum+candidates[i] > target {
				break
			}
			path = append(path, candidates[i])
			backtrack(sum+candidates[i], i)
			path = path[:len(path)-1]
		}
	}

	backtrack(0, 0)

	return res
}

func TestCombinationSum(t *testing.T) {
	tests := []struct {
		name        string
		candidates  []int
		target      int
		expectedLen int
	}{
		{
			name:        "示例1",
			candidates:  []int{2, 3, 6, 7},
			target:      7,
			expectedLen: 2,
		},
		{
			name:        "示例2",
			candidates:  []int{2, 3, 5},
			target:      8,
			expectedLen: 3,
		},
		{
			name:        "示例3-无解",
			candidates:  []int{2},
			target:      1,
			expectedLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := combinationSum(tt.candidates, tt.target)
			if len(got) != tt.expectedLen {
				t.Errorf("combinationSum() returned %d combinations, want %d", len(got), tt.expectedLen)
			}
		})
	}
}

package top100liked

import (
	"testing"
)

// 78. 子集 (Subsets)
//
// 题目描述:
// 给你一个整数数组 nums ，数组中的元素互不相同。返回该数组所有可能的子集（幂集）。
// 解集不能包含重复的子集。你可以按任意顺序返回解集。
//
// 示例 1：
// 输入：nums = [1,2,3]
// 输出：[[],[1],[2],[1,2],[3],[1,3],[2,3],[1,2,3]]
//
// 示例 2：
// 输入：nums = [0]
// 输出：[[],[0]]

func subsets(nums []int) [][]int {
	var res [][]int
	var path []int

	var backtrack func(int)
	backtrack = func(idx int) {
		if idx == len(nums) {
			return
		}

		res = append(res, append([]int(nil), path...))

		for i := idx; i < len(nums); i++ {
			path = append(path, nums[i])
			backtrack(i + 1)
			path = path[:len(path)-1]
		}
	}

	backtrack(0)

	return res

}

func TestSubsets(t *testing.T) {
	tests := []struct {
		name        string
		nums        []int
		expectedLen int
	}{
		{
			name:        "示例1",
			nums:        []int{1, 2, 3},
			expectedLen: 8,
		},
		{
			name:        "示例2",
			nums:        []int{0},
			expectedLen: 2,
		},
		{
			name:        "两个元素",
			nums:        []int{1, 2},
			expectedLen: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := subsets(tt.nums)
			if len(got) != tt.expectedLen {
				t.Errorf("subsets() returned %d subsets, want %d", len(got), tt.expectedLen)
			}
		})
	}
}

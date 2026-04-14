package top100liked

import (
	"fmt"
	"testing"
)

// 46. 全排列 (Permutations)
//
// 题目描述:
// 给定一个不含重复数字的数组 nums ，返回其所有可能的全排列。你可以按任意顺序返回答案。
//
// 示例 1：
// 输入：nums = [1,2,3]
// 输出：[[1,2,3],[1,3,2],[2,1,3],[2,3,1],[3,1,2],[3,2,1]]
//
// 示例 2：
// 输入：nums = [0,1]
// 输出：[[0,1],[1,0]]
//
// 示例 3：
// 输入：nums = [1]
// 输出：[[1]]

func permute(nums []int) [][]int {
	var res [][]int
	path := []int{}

	// used 数组标记当前数字是否在 path 中，避免重复使用
	used := make([]bool, len(nums))

	var backtrack func()
	backtrack = func() {
		// 路径长度等于数组长度，代表找到了
		if len(path) == len(nums) {
			temp := make([]int, len(path))
			copy(temp, path)
			res = append(res, temp)
			return
		}

		for i := 0; i < len(nums); i++ {
			// 如果这个数字已经在路径里了，跳过
			if used[i] {
				continue
			}

			// 【进】：做选择
			path = append(path, nums[i])
			used[i] = true
			fmt.Printf("第 %d 层 path: %v\n", i+1, path)
			fmt.Printf("第 %d 层 used: %v\n", i+1, used)

			// 递归：进入下一层决策树
			backtrack()

			// 【退】：撤销选择（回溯的核心）
			// 这一行把最后加进去的元素踢出去，让 path 恢复到进入这一层循环前的样子
			path = path[:len(path)-1]
			used[i] = false
		}
	}

	backtrack()
	return res
}

func TestPermute(t *testing.T) {
	tests := []struct {
		name        string
		nums        []int
		expectedLen int
	}{
		{
			name:        "示例1",
			nums:        []int{1, 2, 3},
			expectedLen: 6,
		},
		// {
		// 	name:        "示例2",
		// 	nums:        []int{0, 1},
		// 	expectedLen: 2,
		// },
		// {
		// 	name:        "示例3",
		// 	nums:        []int{1},
		// 	expectedLen: 1,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := permute(tt.nums)
			if len(got) != tt.expectedLen {
				t.Errorf("permute() returned %d permutations, want %d", len(got), tt.expectedLen)
			}
		})
	}
}

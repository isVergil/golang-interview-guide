package top100liked

import (
	"reflect"
	"sort"
	"testing"
)

// 15. 三数之和 (3Sum)
//
// 题目描述:
// 给你一个整数数组 nums ，判断是否存在三元组 [nums[i], nums[j], nums[k]]
// 满足 i != j、i != k 且 j != k ，同时还满足 nums[i] + nums[j] + nums[k] == 0。
// 请你返回所有和为 0 且不重复的三元组。注意：答案中不可以包含重复的三元组。
//
// 示例 1：
// 输入：nums = [-1,0,1,2,-1,-4]
// 输出：[[-1,-1,2],[-1,0,1]]
//
// 示例 2：
// 输入：nums = [0,1,1]
// 输出：[]
//
// 示例 3：
// 输入：nums = [0,0,0]
// 输出：[[0,0,0]]

func threeSum(nums []int) [][]int {
	// 先排序
	sort.Ints(nums)
	res := make([][]int, 0)
	for i := 0; i < len(nums)-2; i++ {
		if nums[i] > 0 {
			break
		}

		if i > 0 && nums[i] == nums[i-1] {
			continue
		}

		// 固定 i， l 和 r 变换
		l, r := i+1, len(nums)-1

		for l < r {
			sum := nums[i] + nums[l] + nums[r]
			if sum == 0 {
				res = append(res, []int{nums[i], nums[l], nums[r]})

				// 重复情况跳过
				for l < r && nums[l] == nums[l+1] {
					l++
				}

				// 重复情况跳过
				for l < r && nums[r] == nums[r-1] {
					r--
				}
				l++
				r--
			} else if sum > 0 {
				r--
			} else {
				l++
			}
		}
	}
	return res
}

func TestThreeSum(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected [][]int
	}{
		{
			name:     "示例1",
			nums:     []int{-1, 0, 1, 2, -1, -4},
			expected: [][]int{{-1, -1, 2}, {-1, 0, 1}},
		},
		{
			name:     "示例2",
			nums:     []int{0, 1, 1},
			expected: nil,
		},
		{
			name:     "示例3",
			nums:     []int{0, 0, 0},
			expected: [][]int{{0, 0, 0}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := threeSum(tt.nums)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("threeSum() = %v, want %v", got, tt.expected)
			}
		})
	}
}

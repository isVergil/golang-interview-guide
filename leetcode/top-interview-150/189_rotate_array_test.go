package topinterview150

import (
	"testing"
)

// 189. 轮转数组 (Rotate Array)
//
// 题目描述:
// 给定一个整数数组 nums，将数组中的元素向右轮转 k 个位置，其中 k 是非负数。
//
// 示例 1:
// 输入: nums = [1,2,3,4,5,6,7], k = 3
// 输出: [5,6,7,1,2,3,4]
// 解释:
// 向右轮转 1 步: [7,1,2,3,4,5,6]
// 向右轮转 2 步: [6,7,1,2,3,4,5]
// 向右轮转 3 步: [5,6,7,1,2,3,4]
//
// 示例 2:
// 输入: nums = [-1,-100,3,99], k = 2
// 输出: [3,99,-1,-100]
// 解释:
// 向右轮转 1 步: [99,-1,-100,3]
// 向右轮转 2 步: [3,99,-1,-100]

func rotate(nums []int, k int) {
	res := make([]int, len(nums))
	for i, num := range nums {
		res[(i+k)%len(nums)] = num
	}
	copy(nums, res)
}

func TestRotate(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		k        int
		expected []int
	}{
		{
			name:     "Example 1",
			nums:     []int{1, 2, 3, 4, 5, 6, 7},
			k:        3,
			expected: []int{5, 6, 7, 1, 2, 3, 4},
		},
		{
			name:     "Example 2",
			nums:     []int{-1, -100, 3, 99},
			k:        2,
			expected: []int{3, 99, -1, -100},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: Uncomment when implementation is ready
			// nums := make([]int, len(tt.nums))
			// copy(nums, tt.nums)
			//
			// rotate(nums, tt.k)
			//
			// if !reflect.DeepEqual(nums, tt.expected) {
			// 	t.Errorf("rotate() resulted in %v, want %v", nums, tt.expected)
			// }
		})
	}
}

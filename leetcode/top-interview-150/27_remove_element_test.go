package topinterview150

import (
	"sort"
	"testing"
)

// 27. 移除元素 (Remove Element)
//
// 题目描述:
// 给你一个数组 nums 和一个值 val，你需要 "原地" 移除所有数值等于 val 的元素，并返回移除后数组的新长度。
// 不要使用额外的数组空间，你必须仅使用 O(1) 额外空间并 "原地" 修改输入数组。
// 元素的顺序可以改变。你不需要考虑数组中超出新长度后面的元素。
//
// 示例 1：
// 输入：nums = [3,2,2,3], val = 3
// 输出：2, nums = [2,2,_,_]
//
// 示例 2：
// 输入：nums = [0,1,2,2,3,0,4,2], val = 2
// 输出：5, nums = [0,1,4,0,3,_,_,_]

func removeElement(nums []int, val int) int {
	idx := 0
	for _, num := range nums {
		if num != val {
			nums[idx] = num
			idx++
		}
	}
	return idx
}

func TestRemoveElement(t *testing.T) {
	tests := []struct {
		name         string
		nums         []int
		val          int
		expectedK    int
		expectedNums []int
	}{
		{
			name:         "Example 1",
			nums:         []int{3, 2, 2, 3},
			val:          3,
			expectedK:    2,
			expectedNums: []int{2, 2},
		},
		{
			name:         "Example 2",
			nums:         []int{0, 1, 2, 2, 3, 0, 4, 2},
			val:          2,
			expectedK:    5,
			expectedNums: []int{0, 1, 3, 0, 4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nums := make([]int, len(tt.nums))
			copy(nums, tt.nums)

			k := removeElement(nums, tt.val)

			if k != tt.expectedK {
				t.Errorf("removeElement() = %v, want %v", k, tt.expectedK)
			}

			actualSlice := nums[:k]
			sort.Ints(actualSlice)
			sort.Ints(tt.expectedNums)
			for i := 0; i < k; i++ {
				if actualSlice[i] != tt.expectedNums[i] {
					t.Errorf("element at index %d = %v, want %v", i, actualSlice[i], tt.expectedNums[i])
				}
			}
		})
	}
}

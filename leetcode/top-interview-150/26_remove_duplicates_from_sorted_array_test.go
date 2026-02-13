package topinterview150

import (
	"reflect"
	"testing"
)

// 26. 删除有序数组中的重复项 (Remove Duplicates from Sorted Array)
//
// 题目描述:
// 给你一个 "非严格递增排列" 的数组 nums ，请你 "原地" 删除重复出现的元素，使每个元素 "只出现一次" ，返回删除后数组的新长度。
// 元素的 "相对顺序" 应该保持 一致 。然后返回 nums 中唯一元素的个数。
//
// 示例 1：
// 输入：nums = [1,1,2]
// 输出：2, nums = [1,2,_]
//
// 示例 2：
// 输入：nums = [0,0,1,1,1,2,2,3,3,4]
// 输出：5, nums = [0,1,2,3,4]
func removeDuplicates(nums []int) int {
	// 在这里写入你的代码
	// 提示：数组是有序的，重复元素一定相邻。使用双指针维护不重复序列的尾部。
	idx, lastNum := 1, nums[0]
	for _, num := range nums {
		if num != lastNum {
			nums[idx] = num
			idx++
			lastNum = num
		}
	}
	return idx
}

func TestRemoveDuplicates(t *testing.T) {
	tests := []struct {
		name         string
		nums         []int
		expected     int
		expectedNums []int
	}{
		{
			name:         "Example 1",
			nums:         []int{1, 1, 2},
			expected:     2,
			expectedNums: []int{1, 2},
		},
		{
			name:         "Example 2",
			nums:         []int{0, 0, 1, 1, 1, 2, 2, 3, 3, 4},
			expected:     5,
			expectedNums: []int{0, 1, 2, 3, 4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nums := make([]int, len(tt.nums))
			copy(nums, tt.nums)

			k := removeDuplicates(nums)

			if k != tt.expected {
				t.Errorf("removeDuplicates() = %v, want %v", k, tt.expected)
			}
			if !reflect.DeepEqual(nums[:k], tt.expectedNums) {
				t.Errorf("nums[:k] = %v, want %v", nums[:k], tt.expectedNums)
			}
		})
	}
}

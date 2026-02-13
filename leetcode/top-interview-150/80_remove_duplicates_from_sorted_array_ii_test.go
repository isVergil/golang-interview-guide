package topinterview150

import (
	"testing"
)

// 80. 删除有序数组中的重复项 II (Remove Duplicates from Sorted Array II)
//
// 题目描述:
// 给你一个有序数组 nums ，请你 "原地" 删除重复出现的元素，使得出现次数超过两次的元素只出现两次 ，返回删除后数组的新长度。
// 不要使用额外的数组空间，你必须在 "原地" 修改输入数组 并在使用 O(1) 额外空间的条件下完成。
//
// 示例 1：
// 输入：nums = [1,1,1,2,2,3]
// 输出：5, nums = [1,1,2,2,3]
// 解释：函数应返回新长度 length = 5, 并且原数组的前五个元素被修改为 1, 1, 2, 2, 3。 不需要考虑数组中超出新长度后面的元素。
//
// 示例 2：
// 输入：nums = [0,0,1,1,1,1,2,3,3]
// 输出：7, nums = [0,0,1,1,2,3,3]
// 解释：函数应返回新长度 length = 7, 并且原数组的前七个元素被修改为 0, 0, 1, 1, 2, 3, 3。 不需要考虑数组中超出新长度后面的元素。

func removeDuplicatesII(nums []int) int {
	n := len(nums)

	// 如果长度小于等于 2，直接返回 n，因为最多允许重复两次
	if n <= 2 {
		return n
	}

	// slow 从 2 开始，因为前两个元素无论如何都会被保留
	slow := 2

	for fast := 2; fast < n; fast++ {
		if nums[fast] != nums[slow-2] {
			nums[slow] = nums[fast]
			slow++
		}
	}
	return slow

}

func TestRemoveDuplicatesII(t *testing.T) {
	tests := []struct {
		name         string
		nums         []int
		expected     int
		expectedNums []int
	}{
		{
			name:         "Example 1",
			nums:         []int{1, 1, 1, 2, 2, 3},
			expected:     5,
			expectedNums: []int{1, 1, 2, 2, 3},
		},
		{
			name:         "Example 2",
			nums:         []int{0, 0, 1, 1, 1, 1, 2, 3, 3},
			expected:     7,
			expectedNums: []int{0, 0, 1, 1, 2, 3, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: Uncomment when implementation is ready
			// nums := make([]int, len(tt.nums))
			// copy(nums, tt.nums)
			//
			// k := removeDuplicatesII(nums)
			//
			// if k != tt.expected {
			// 	t.Errorf("removeDuplicatesII() = %v, want %v", k, tt.expected)
			// }
			// if !reflect.DeepEqual(nums[:k], tt.expectedNums) {
			// 	t.Errorf("nums[:k] = %v, want %v", nums[:k], tt.expectedNums)
			// }
		})
	}
}

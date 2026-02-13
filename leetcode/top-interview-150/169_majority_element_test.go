package topinterview150

import "testing"

// 169. 多数元素 (Majority Element)
//
// 题目描述:
// 给定一个大小为 n 的数组 nums ，返回其中的多数元素。
// 多数元素是指在数组中出现次数 "大于" ⌊ n/2 ⌋ 的元素。
//
// 示例 1：
// 输入：nums = [3,2,3]
// 输出：3
//
// 示例 2：
// 输入：nums = [2,2,1,1,1,2,2]
// 输出：2

func majorityElement(nums []int) int {
	// 法1 map 计数 o(n)
	// counts := make(map[int]int)
	// limit := len(nums) / 2
	// for _, num := range nums {
	// 	counts[num]++
	// 	if counts[num] > limit {
	// 		return num
	// 	}
	// }
	// return -1

	// 法2 摩尔投票法 o(1)
	counts, cur := 0, nums[0]
	for _, num := range nums {
		if counts == 0 {
			cur = num
		}
		if cur != num {
			counts--
		} else {
			counts++
		}
	}
	return cur
}

func TestMajorityElement(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected int
	}{
		{
			name:     "Example 1",
			nums:     []int{3, 2, 3},
			expected: 3,
		},
		{
			name:     "Example 2",
			nums:     []int{2, 2, 1, 1, 1, 2, 2},
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: Uncomment when implementation is ready
			if got := majorityElement(tt.nums); got != tt.expected {
				t.Errorf("majorityElement() = %v, want %v", got, tt.expected)
			}
		})
	}
}

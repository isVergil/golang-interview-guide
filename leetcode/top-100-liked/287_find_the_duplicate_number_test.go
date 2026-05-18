package top100liked

import "testing"

// 287. 寻找重复数 (Find the Duplicate Number)
//
// 题目描述:
// 给定一个包含 n + 1 个整数的数组 nums，其数字都在 [1, n] 范围内（包括 1 和 n），
// 可知至少存在一个重复的整数。假设 nums 只有一个重复的整数，返回这个重复的数。
// 要求：不修改数组，空间复杂度 O(1)。
//
// 示例 1：
// 输入：nums = [1,3,4,2,2]
// 输出：2
//
// 示例 2：
// 输入：nums = [3,1,3,4,2]
// 输出：3
//
// 提示：把数组看成链表（index→value），用快慢指针找环入口

func findDuplicate(nums []int) int {
	slow, fast := nums[0], nums[nums[0]]
	for slow != fast {
		slow = nums[slow]
		fast = nums[nums[fast]]
	}

	slow = 0
	for slow != fast {
		slow = nums[slow]
		fast = nums[fast]
	}
	return slow
}

func TestFindDuplicate(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected int
	}{
		{name: "示例1", nums: []int{1, 3, 4, 2, 2}, expected: 2},
		{name: "示例2", nums: []int{3, 1, 3, 4, 2}, expected: 3},
		{name: "重复多次", nums: []int{2, 2, 2, 2, 2}, expected: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findDuplicate(tt.nums)
			if got != tt.expected {
				t.Errorf("findDuplicate() = %v, want %v", got, tt.expected)
			}
		})
	}
}

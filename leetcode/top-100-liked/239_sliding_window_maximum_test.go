package top100liked

import (
	"reflect"
	"testing"
)

// 239. 滑动窗口最大值 (Sliding Window Maximum)
//
// 题目描述:
// 给你一个整数数组 nums，有一个大小为 k 的滑动窗口从数组的最左侧移动到最右侧。
// 你只可以看到在滑动窗口内的 k 个数字。滑动窗口每次只向右移动一位。
// 返回滑动窗口中的最大值。
//
// 示例：
// 输入：nums = [1,3,-1,-3,5,3,6,7], k = 3
// 输出：[3,3,5,5,6,7]
//
// 提示：单调递减双端队列，队头始终是窗口内最大值的下标

func maxSlidingWindow(nums []int, k int) []int {
	n := len(nums)
	if n == 0 || k == 0 {
		return nil
	}
	res := make([]int, 0, n-k+1)
	dq := make([]int, 0, k)
	for i := 0; i < n; i++ {
		// 出窗口了
		if len(dq) > 0 && dq[0] <= i-k {
			dq = dq[1:]
		}

		// 维持单调递减队列
		for len(dq) > 0 && nums[dq[len(dq)-1]] <= nums[i] {
			dq = dq[:len(dq)-1]
		}
		dq = append(dq, i)
		if i >= k-1 {
			res = append(res, nums[dq[0]])
		}
	}
	return res
}

func TestMaxSlidingWindow(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		k        int
		expected []int
	}{
		{
			name:     "示例",
			nums:     []int{1, 3, -1, -3, 5, 3, 6, 7},
			k:        3,
			expected: []int{3, 3, 5, 5, 6, 7},
		},
		{
			name:     "k=1",
			nums:     []int{1, -1},
			k:        1,
			expected: []int{1, -1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := maxSlidingWindow(tt.nums, tt.k)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("maxSlidingWindow() = %v, want %v", got, tt.expected)
			}
		})
	}
}

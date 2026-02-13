package topinterview150

import (
	"fmt"
	"testing"
)

// 228. 汇总区间 (Summary Ranges)
//
// 题目描述:
// 给定一个  无重复元素 的 有序 整数数组 nums 。
// 返回 恰好覆盖数组中所有数字 的 最小有序 区间范围列表 。也就是说，nums 的每个元素都恰好被某个区间范围所覆盖，并且不存在属于某个范围但不属于 nums 的数字 x 。
// 列表中的每个区间范围 [a,b] 应该按如下格式输出：
// "a->b" ，如果 a != b
// "a" ，如果 a == b
//
// 示例 1：
// 输入：nums = [0,1,2,4,5,7]
// 输出：["0->2","4->5","7"]
// 解释：区间范围是：
// [0,2] --> "0->2"
// [4,5] --> "4->5"
// [7,7] --> "7"
//
// 示例 2：
// 输入：nums = [0,2,3,4,6,8,9]
// 输出：["0","2->4","6","8->9"]
// 解释：区间范围是：
// [0,0] --> "0"
// [2,4] --> "2->4"
// [6,6] --> "6"
// [8,9] --> "8->9"

func summaryRanges(nums []int) []string {
	res := []string{}

	// 双指针
	for idx := 0; idx < len(nums); {
		start := idx
		for idx+1 < len(nums) && nums[idx+1] == nums[idx]+1 {
			idx++
		}
		if idx == start {
			res = append(res, fmt.Sprintf("%d", nums[idx]))
		} else {
			res = append(res, fmt.Sprintf("%d->%d", nums[start], nums[idx]))
		}
		idx++
	}
	return res
}

func TestSummaryRanges(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected []string
	}{
		{
			name:     "Example 1",
			nums:     []int{0, 1, 2, 4, 5, 7},
			expected: []string{"0->2", "4->5", "7"},
		},
		{
			name:     "Example 2",
			nums:     []int{0, 2, 3, 4, 6, 8, 9},
			expected: []string{"0", "2->4", "6", "8->9"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: Uncomment when implementation is ready
			// if got := summaryRanges(tt.nums); !reflect.DeepEqual(got, tt.expected) {
			// 	t.Errorf("summaryRanges() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

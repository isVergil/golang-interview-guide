package topinterview150

import "testing"

// 219. 存在重复元素 II (Contains Duplicate II)
//
// 题目描述:
// 给你一个整数数组 nums 和一个整数 k ，判断数组中是否存在两个 不同的索引 i 和 j ，满足 nums[i] == nums[j] 且 abs(i - j) <= k 。
// 如果存在，返回 true ；否则，返回 false 。
//
// 示例 1：
// 输入：nums = [1,2,3,1], k = 3
// 输出：true
//
// 示例 2：
// 输入：nums = [1,0,1,1], k = 1
// 输出：true
//
// 示例 3：
// 输入：nums = [1,2,3,1,2,3], k = 2
// 输出：false

// 1 哈希表
func containsNearbyDuplicate(nums []int, k int) bool {
	idxMap := make(map[int]int)

	for idx, num := range nums {
		if prevIdx, ok := idxMap[num]; ok {
			if idx-prevIdx <= k {
				return true
			}
		}
		idxMap[num] = idx
	}
	return false
}

// 2 滑动窗口 维护固定长度 k 的map 如果遍历元素包含在 map 中则为 true
func containsNearbyDuplicate2(nums []int, k int) bool {
	// 篮子：存储当前窗口内的数字
	// 使用 struct{} 是为了节省内存，它不占空间
	window := make(map[int]struct{})

	for i := 0; i < len(nums); i++ {
		// 1. 检查：当前数字是否已经在最近的 k 个数里了？
		if _, ok := window[nums[i]]; ok {
			return true
		}

		// 2. 放入：把当前数字加入窗口
		window[nums[i]] = struct{}{}

		// 3. 维护：如果窗口大小超过了 k
		// 也就是当 i 达到了 k 的时候，下一次循环前必须删掉最左边的数
		if len(window) > k {
			// 删掉下标为 i-k 的那个数字，它是窗口里最老的
			delete(window, nums[i-k])
		}
	}

	return false
}

func TestContainsNearbyDuplicate(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		k        int
		expected bool
	}{
		{
			name:     "Example 1",
			nums:     []int{1, 2, 3, 1},
			k:        3,
			expected: true,
		},
		{
			name:     "Example 2",
			nums:     []int{1, 0, 1, 1},
			k:        1,
			expected: true,
		},
		{
			name:     "Example 3",
			nums:     []int{1, 2, 3, 1, 2, 3},
			k:        2,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: Uncomment when implementation is ready
			// if got := containsNearbyDuplicate(tt.nums, tt.k); got != tt.expected {
			// 	t.Errorf("containsNearbyDuplicate() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

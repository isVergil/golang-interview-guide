package topinterview150

import "testing"

// 45. 跳跃游戏 II (Jump Game II)
//
// 题目描述:
// 给定一个长度为 n 的 0 索引整数数组 nums。初始位置为 nums[0]。
// 每个元素 nums[i] 表示从索引 i 向前跳转的最大长度。
// 换句话说，如果你在 nums[i] 处，你可以跳转到任意 nums[i + j] 处:
// 0 <= j <= nums[i]
// i + j < n
// 返回到达 nums[n - 1] 的最小跳跃次数。生成的测试用例可以到达 nums[n - 1]。
//
// 示例 1:
// 输入: nums = [2,3,1,1,4]
// 输出: 2
// 解释: 跳到最后一个位置的最小跳跃数是 2。
//      从下标为 0 跳到下标为 1 的位置，跳 1 步，然后跳 3 步到达数组的最后一个位置。
//
// 示例 2:
// 输入: nums = [2,3,0,1,4]
// 输出: 2

func jump(nums []int) int {
	// 边界情况
	n := len(nums)
	if n <= 1 {
		return 0
	}

	// 步数、探测的最远、当前步数能达到的最远
	step, maxLeap, end := 0, 0, 0
	for idx := 0; idx < n-1; idx++ {
		if nums[idx]+idx > maxLeap {
			maxLeap = nums[idx] + idx
		}

		// 遍历超过了 探测的最远 就没法到达了
		if idx > maxLeap {
			return -1
		}

		// 到达当前步数能到达的最远 不得不跳了
		if idx == end {
			step++
			end = maxLeap
			// 当前步数能达到最后 提前结束
			if end >= n-1 {
				return step
			}
		}
	}
	return step
}

func TestJump(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected int
	}{
		{
			name:     "Example 1",
			nums:     []int{2, 3, 1, 1, 4},
			expected: 2,
		},
		{
			name:     "Example 2",
			nums:     []int{2, 3, 0, 1, 4},
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: Uncomment when implementation is ready
			// if got := jump(tt.nums); got != tt.expected {
			// 	t.Errorf("jump() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

package top100liked

import (
	"testing"
)

// 45. 跳跃游戏 II (Jump Game II)
//
// 题目描述:
// 给定一个长度为 n 的 0 索引整数数组 nums。初始位置为 nums[0]。
// 每个元素 nums[i] 表示从索引 i 向前跳转的最大长度。返回到达 nums[n-1] 的最小跳跃次数。
// 生成的测试用例可以到达 nums[n-1]。
//
// 示例 1：
// 输入：nums = [2,3,1,1,4]
// 输出：2（跳到索引 1，然后跳到最后）
//
// 示例 2：
// 输入：nums = [2,3,0,1,4]
// 输出：2

// 每一跳贪心地选"这一层能到的最远距离"作为下一层的边界
// 这就像 BFS 的层序遍历：
//   - 每一"层"是当前跳跃次数能到的所有位置
//   - curEnd 是当前层的右边界
//   - farthest 是下一层的右边界
//   - 走完一层（i == curEnd），jumps++，进入下一层
func jump(nums []int) int {
	n := len(nums)
	jumps := 0
	curEnd := 0
	farthest := 0

	for i := 0; i < n-2; i++ {
		if i+nums[i] > farthest {
			farthest = i + nums[i]
		}
		if i == curEnd {
			jumps++
			curEnd = farthest
		}
	}

	return jumps
}

func TestJump(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected int
	}{
		{name: "示例1", nums: []int{2, 3, 1, 1, 4}, expected: 2},
		{name: "示例2", nums: []int{2, 3, 0, 1, 4}, expected: 2},
		{name: "单元素", nums: []int{0}, expected: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := jump(tt.nums)
			if got != tt.expected {
				t.Errorf("jump() = %v, want %v", got, tt.expected)
			}
		})
	}
}

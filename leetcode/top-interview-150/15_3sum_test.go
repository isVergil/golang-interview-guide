package topinterview150

import (
	"sort"
	"testing"
)

// 15. 三数之和 (3Sum)
//
// 题目描述:
// 给你一个包含 n 个整数的数组 nums，判断 nums 中是否存在三个元素 a，b，c ，使得 a + b + c = 0 ？请你找出所有和为 0 且不重复的三元组。
// 注意：答案中不可以包含重复的三元组。
//
// 示例 1：
// 输入：nums = [-1,0,1,2,-1,-4]
// 输出：[[-1,-1,2],[-1,0,1]]

func threeSum(nums []int) [][]int {
	n := len(nums)

	// 必须先排序
	sort.Ints(nums)
	res := make([][]int, 0)

	for i := 0; i < n-2; i++ {
		// 由于数组已排序，如果 nums[i] > 0，后面三个数之和肯定也大于 0
		if nums[i] > 0 {
			break
		}

		// 去重。如果当前数字和上一个数字相同，直接跳过，避免重复组合
		if i > 0 && nums[i] == nums[i-1] {
			continue
		}

		l := i + 1
		r := n - 1
		for l < r {
			sum := nums[i] + nums[l] + nums[r]
			if sum == 0 {
				res = append(res, []int{nums[i], nums[l], nums[r]})

				// 左指针去重。跳过所有与当前 nums[left] 相同的数
				for l < r && nums[l] == nums[l+1] {
					l++
				}

				// 右指针去重。跳过所有与当前 nums[right] 相同的数
				for l < r && nums[r] == nums[r-1] {
					r--
				}

				l++
				r--
			} else if sum < 0 {
				// 总和太小，左指针右移以增大数值
				l++
			} else {
				// 总和太大，右指针左移以减小数值
				r--
			}
		}
	}
	return res
}

func TestThreeSum(t *testing.T) {
	// 三数之和测试
}

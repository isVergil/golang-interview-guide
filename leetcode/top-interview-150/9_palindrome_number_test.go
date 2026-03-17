package topinterview150

import (
	"testing"
)

// 9. 回文数 (Palindrome Number)
//
// 题目描述:
// 给你一个整数 x ，如果 x 是一个回文整数，返回 true ；否则，返回 false 。
// 回文数是指正序（从左向右）和倒序（从右向左）读都是一样的整数。
// 例如，121 是回文，而 123 不是。
//
// 示例 1：
// 输入：x = 121
// 输出：true
//
// 示例 2：
// 输入：x = -121
// 输出：false
// 解释：从左向右读, 为 -121 。 从右向左读, 为 121- 。因此它不是一个回文数。
//
// 示例 3：
// 输入：x = 10
// 输出：false
// 解释：从右向左读, 为 01 。因此它不是一个回文数。

func isPalindromeNumber(x int) bool {
	// 边界条件优化 (细节把控)：
	// 1. 负数绝不可能是回文数 (例如 -121 -> 121-)
	// 2. 最后一位是 0 的数字，只有 0 本身是回文数 (例如 10 -> 01，不成立)
	// 提前拦截可以省去后续大量的算力
	if x < 0 || (x%10 == 0 && x != 0) {
		return false
	}

	revertedHalf := 0
	// 核心逻辑：只反转数字的后半部分
	// 当原数字 x 小于或等于反转后的数字 revertedHalf 时，说明我们已经处理了一半以上的位数
	for x > revertedHalf {
		revertedHalf = revertedHalf*10 + x%10
		x /= 10
	}

	// 此时有两种情况：
	// 1. 偶数长度：比如 1221，此时 x == 12，revertedHalf == 12。直接判断 x == revertedHalf
	// 2. 奇数长度：比如 12321，此时 x == 12，revertedHalf == 123。中间的 3 对于回文没有影响，通过 revertedHalf/10 去掉尾数后比较
	return x == revertedHalf || x == revertedHalf/10
}

func TestIsPalindromeNumber(t *testing.T) {
	tests := []struct {
		name     string
		x        int
		expected bool
	}{
		{"Example 1", 121, true},
		{"Example 2", -121, false},
		{"Example 3", 10, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := isPalindromeNumber(tt.x); got != tt.expected {
			// 	t.Errorf("isPalindromeNumber() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

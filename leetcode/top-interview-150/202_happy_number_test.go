package topinterview150

import "testing"

// 202. 快乐数 (Happy Number)
//
// 题目描述:
// 编写一个算法来判断一个数 n 是不是快乐数。
// 「快乐数」 定义为：
// 对于一个正整数，每一次将该数替换为它每个位置上的数字的平方和。
// 然后重复这个过程直到这个数变为 1，也可能是 无限循环 但始终变不到 1。
// 如果这个过程 结果为 1，那么这个数就是快乐数。
// 如果 n 是 快乐数 就返回 true ；不是，则返回 false 。
//
// 示例 1：
// 输入：n = 19
// 输出：true
// 解释：
// 1^2 + 9^2 = 82
// 8^2 + 2^2 = 68
// 6^2 + 8^2 = 100
// 1^2 + 0^2 + 0^2 = 1
//
// 示例 2：
// 输入：n = 2
// 输出：false

// 1 快慢指针
func isHappy(n int) bool {
	slow, fast := n, getNext(n)

	// 如果快慢指针不相等，且快指针还没到达 1
	for fast != 1 && slow != fast {
		slow = getNext(slow)          // 走一步
		fast = getNext(getNext(fast)) // 走两步
	}

	// 如果是因为 fast == 1 出来的，说明是快乐数
	return fast == 1
}

// 2 哈希表
func isHappy1(n int) bool {
	// 使用 map 作为哈希表，记录出现过的数字
	m := make(map[int]bool)

	for n != 1 {
		if m[n] {
			return false
		}
		m[n] = true
		n = getNext(n)
	}

	// 如果是因为 fast == 1 出来的，说明是快乐数
	return true
}

// 辅助函数：计算各位数字的平方和
func getNext(n int) int {
	sum := 0
	for n > 0 {
		digit := n % 10
		sum += digit * digit
		n = n / 10
	}
	return sum
}

func TestIsHappy(t *testing.T) {
	tests := []struct {
		name     string
		n        int
		expected bool
	}{
		{
			name:     "Example 1",
			n:        19,
			expected: true,
		},
		{
			name:     "Example 2",
			n:        2,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: Uncomment when implementation is ready
			// if got := isHappy(tt.n); got != tt.expected {
			// 	t.Errorf("isHappy() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

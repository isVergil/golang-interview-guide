package top100liked

import (
	"reflect"
	"testing"
)

// 347. 前 K 个高频元素 (Top K Frequent Elements)
//
// 题目描述:
// 给你一个整数数组 nums 和一个整数 k，请你返回其中出现频率前 k 高的元素。
// 你可以按任意顺序返回答案，题目保证答案唯一。
//
// 示例 1：
// 输入：nums = [1,1,1,2,2,3], k = 2
// 输出：[1,2]
//
// 示例 2：
// 输入：nums = [1], k = 1
// 输出：[1]

func topKFrequent(nums []int, k int) []int {
	// 统计频率
	freq := make(map[int]int)
	for _, num := range nums {
		freq[num]++
	}

	// 桶排序：下标是频率，值是元素列表
	n := len(nums)
	buckets := make([][]int, n+1)
	for num, cnt := range freq {
		buckets[cnt] = append(buckets[cnt], num)
	}

	// 从高频到低频取 k 个
	res := make([]int, 0, k)
	for i := n; i >= 0 && len(res) < k; i-- {
		res = append(res, buckets[i]...)
	}
	return res
}

func TestTopKFrequent(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		k        int
		expected []int
	}{
		{name: "示例1", nums: []int{1, 1, 1, 2, 2, 3}, k: 2, expected: []int{1, 2}},
		{name: "示例2", nums: []int{1}, k: 1, expected: []int{1}},
		{name: "全部相同", nums: []int{3, 3, 3}, k: 1, expected: []int{3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := topKFrequent(tt.nums, tt.k)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("topKFrequent() = %v, want %v", got, tt.expected)
			}
		})
	}
}

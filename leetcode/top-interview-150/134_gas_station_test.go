package topinterview150

import (
	"testing"
)

// 134. 加油站 (Gas Station)
//
// 题目描述:
// 在一条环路上有 n 个加油站，其中第 i 个加油站有汽油 gas[i] 升。
// 你有一辆油箱容量无限的汽车，从第 i 个加油站开往第 i+1 个加油站需要消耗汽油 cost[i] 升。你从其中的一个加油站出发，开始时油箱为空。
// 给定两个整数数组 gas 和 cost ，如果你可以按顺时针方向绕环路行驶一周，则返回出发时的加油站编号，否则返回 -1 。如果存在解，则 保证 它是 唯一的。
//
// 示例 1:
// 输入: gas = [1,2,3,4,5], cost = [3,4,5,1,2]
// 输出: 3
//
// 示例 2:
// 输入: gas = [2,3,4], cost = [3,4,3]
// 输出: -1

func canCompleteCircuit(gas []int, cost []int) int {
	totalSum, curSum, start := 0, 0, 0

	for i := 0; i < len(gas); i++ {
		diff := gas[i] - cost[i]
		totalSum += diff
		curSum += diff

		if curSum < 0 {
			start = i + 1
			curSum = 0
		}
	}

	if totalSum < 0 {
		return -1
	}

	return start
}

func TestCanCompleteCircuit(t *testing.T) {
	tests := []struct {
		name     string
		gas      []int
		cost     []int
		expected int
	}{
		{"Example 1", []int{1, 2, 3, 4, 5}, []int{3, 4, 5, 1, 2}, 3},
		{"Example 2", []int{2, 3, 4}, []int{3, 4, 3}, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := canCompleteCircuit(tt.gas, tt.cost); got != tt.expected {
			// 	t.Errorf("canCompleteCircuit() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

package top100liked

import (
	"testing"
)

// 739. 每日温度 (Daily Temperatures)
//
// 题目描述:
// 给定一个整数数组 temperatures ，表示每天的温度，返回一个数组 answer ，其中 answer[i] 是指对于第 i 天，下一个更高温度出现在几天后。如果气温在这之后都不会升高，请在该位置用 0 来代替。
//
// 示例 1:
// 输入: temperatures = [73,74,75,71,69,72,76,73]
// 输出: [1,1,4,2,1,1,0,0]
//
// 示例 2:
// 输入: temperatures = [30,40,50,60]
// 输出: [1,1,1,0]
//
// 示例 3:
// 输入: temperatures = [30,60,90]
// 输出: [1,1,0]

func dailyTemperatures(temperatures []int) []int {
	panic("not implemented")
}

func TestDailyTemperatures(t *testing.T) {
	tests := []struct {
		name         string
		temperatures []int
		expected     []int
	}{
		{"Example 1", []int{73, 74, 75, 71, 69, 72, 76, 73}, []int{1, 1, 4, 2, 1, 1, 0, 0}},
		{"Example 2", []int{30, 40, 50, 60}, []int{1, 1, 1, 0}},
		{"Example 3", []int{30, 60, 90}, []int{1, 1, 0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := dailyTemperatures(tt.temperatures); !reflect.DeepEqual(got, tt.expected) {
			// 	t.Errorf("dailyTemperatures() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

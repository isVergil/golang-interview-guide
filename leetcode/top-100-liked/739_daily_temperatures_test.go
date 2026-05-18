package top100liked

import (
	"reflect"
	"testing"
)

// 739. 每日温度 (Daily Temperatures)
//
// 题目描述:
// 给定一个整数数组 temperatures，表示每天的温度，返回一个数组 answer，
// 其中 answer[i] 表示对于第 i 天，下一个更高温度出现在几天后。
// 如果气温在这之后都不会升高，用 0 代替。
//
// 示例 1：
// 输入：temperatures = [73,74,75,71,69,72,76,73]
// 输出：[1,1,4,2,1,1,0,0]
//
// 示例 2：
// 输入：temperatures = [30,40,50,60]
// 输出：[1,1,1,0]
//
// 提示：单调递减栈，当前温度 > 栈顶时弹出计算天数差

func dailyTemperatures(temperatures []int) []int {
	n := len(temperatures)
	res := make([]int, n)
	stack := []int{}
	for i := 0; i < n; i++ {
		for len(stack) > 0 && temperatures[i] > temperatures[stack[len(stack)-1]] {
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			res[top] = i - top
		}
		stack = append(stack, i)
	}
	return res
}

func TestDailyTemperatures(t *testing.T) {
	tests := []struct {
		name         string
		temperatures []int
		expected     []int
	}{
		{
			name:         "示例1",
			temperatures: []int{73, 74, 75, 71, 69, 72, 76, 73},
			expected:     []int{1, 1, 4, 2, 1, 1, 0, 0},
		},
		{
			name:         "递增",
			temperatures: []int{30, 40, 50, 60},
			expected:     []int{1, 1, 1, 0},
		},
		{
			name:         "递减",
			temperatures: []int{60, 50, 40, 30},
			expected:     []int{0, 0, 0, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := dailyTemperatures(tt.temperatures)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("dailyTemperatures() = %v, want %v", got, tt.expected)
			}
		})
	}
}

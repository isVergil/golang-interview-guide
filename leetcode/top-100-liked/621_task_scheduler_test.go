package top100liked

import (
	"reflect"
	"testing"
)

// 621. 任务调度器 (Task Scheduler)
//
// 题目描述:
// 给你一个用字符数组 tasks 表示的 CPU 需要执行的任务列表，每个字母表示一种不同种类的任务。
// 任务可以以任意顺序执行，并且每个任务都可以在 1 个单位时间内执行完。
// 在任何一个单位时间，CPU 可以完成一个任务，或者处于待命状态。
// 两个相同种类的任务之间必须有长度为整数 n 的冷却时间。
// 返回完成所有任务所需要的最短时间。
//
// 示例 1：
// 输入：tasks = ["A","A","A","B","B","B"], n = 2
// 输出：8
// 解释：A -> B -> idle -> A -> B -> idle -> A -> B
//
// 示例 2：
// 输入：tasks = ["A","A","A","B","B","B"], n = 0
// 输出：6
//
// 提示：贪心，按出现次数最多的任务排布，计算空闲槽位
//
//	  (maxFreq - 1) * (n + 1) + maxCount
//	   ↑               ↑         ↑
//	   前面几轮        每轮长度    最后一轮的任务数
//		       ┌─── n+1 ───┐
//		轮 1:  │ A  x  x  x │
//		轮 2:  │ A  x  x  x │    maxFreq-1 = 2 轮，每轮 n+1 个位置
//		       └────────────┘
//		轮 3:  │ A │              最后一轮只放"最高频"的那些任务
//		       └──┘
//		        ↑ maxCount 个
//
// x 可以是其他任务，也可以是 idle。其他任务不够填就只能空转（idle），够填就没有 idle，这时 len(tasks) 更大。
func leastInterval(tasks []byte, n int) int {
	var freq [26]int
	for _, t := range tasks {
		freq[t-'A']++
	}

	// 找最高频次
	maxFreq := 0
	for _, v := range freq {
		if v > maxFreq {
			maxFreq = v
		}
	}

	maxCount := 0
	for _, f := range freq {
		if f == maxFreq {
			maxCount++
		}
	}

	res := (maxFreq-1)*(n+1) + maxCount
	if len(tasks) > res {
		return len(tasks)
	}
	return res
}

func TestLeastInterval(t *testing.T) {
	tests := []struct {
		name     string
		tasks    []byte
		n        int
		expected int
	}{
		{
			name:     "示例1",
			tasks:    []byte{'A', 'A', 'A', 'B', 'B', 'B'},
			n:        2,
			expected: 8,
		},
		{
			name:     "无冷却",
			tasks:    []byte{'A', 'A', 'A', 'B', 'B', 'B'},
			n:        0,
			expected: 6,
		},
		{
			name:     "多种任务",
			tasks:    []byte{'A', 'A', 'A', 'A', 'A', 'A', 'B', 'C', 'D', 'E', 'F', 'G'},
			n:        2,
			expected: 16,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := leastInterval(tt.tasks, tt.n)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("leastInterval() = %v, want %v", got, tt.expected)
			}
		})
	}
}

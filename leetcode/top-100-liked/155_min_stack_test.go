package top100liked

import "testing"

// 155. 最小栈 (Min Stack)
//
// 题目描述:
// 设计一个支持 push、pop、top 操作，并能在常数时间内检索到最小元素的栈。
// - push(val) 将元素 val 推入栈中
// - pop()     删除栈顶的元素
// - top()     获取栈顶元素
// - getMin()  获取栈中的最小元素
// 所有操作要求 O(1) 时间复杂度。
//
// 示例：
// 输入：["MinStack","push","push","push","getMin","pop","top","getMin"]
//       [[],[-2],[0],[-3],[],[],[],[]]
// 输出：[null,null,null,null,-3,null,0,-2]
//
// 思路（双栈法）：
//   主栈 stack    ：正常保存所有元素
//   辅助栈 minStack：栈顶永远是「主栈当前范围内的最小值」
// 同步维护两个栈即可在 O(1) 内拿到最小值。
//
// 关键点：
//  1) Push 时用 val <= minTop 判断（必须是 <=，重复最小值要全部入栈，
//     否则 Pop 掉一个相等值后 GetMin 会得到错误结果）。
//  2) Pop 时只有当弹出的值等于 minStack 栈顶，才同步弹出 minStack。

// MinStack 双栈实现的最小栈
type MinStack struct {
	stack    []int // 主栈
	minStack []int // 辅助栈，栈顶 = 主栈当前最小值
}

// NewMinStack 构造函数，预分配容量减少扩容
func NewMinStack() MinStack {
	return MinStack{
		stack:    make([]int, 0, 16),
		minStack: make([]int, 0, 16),
	}
}

// Push 将 val 压入主栈；若 val <= 当前最小值，同步压入辅助栈
func (s *MinStack) Push(val int) {
	s.stack = append(s.stack, val)
	if len(s.minStack) == 0 || val <= s.minStack[len(s.minStack)-1] {
		s.minStack = append(s.minStack, val)
	}
}

// Pop 弹出主栈栈顶；若它正好是当前最小值，同步弹出辅助栈
func (s *MinStack) Pop() {
	top := s.stack[len(s.stack)-1]
	s.stack = s.stack[:len(s.stack)-1]
	if top == s.minStack[len(s.minStack)-1] {
		s.minStack = s.minStack[:len(s.minStack)-1]
	}
}

// Top 返回主栈栈顶
func (s *MinStack) Top() int {
	return s.stack[len(s.stack)-1]
}

// GetMin 返回当前最小值（辅助栈栈顶）
func (s *MinStack) GetMin() int {
	return s.minStack[len(s.minStack)-1]
}

func TestMinStack(t *testing.T) {
	s := NewMinStack()
	s.Push(-2)
	s.Push(0)
	s.Push(-3)

	if got := s.GetMin(); got != -3 {
		t.Errorf("GetMin() = %v, want -3", got)
	}

	s.Pop()

	if got := s.Top(); got != 0 {
		t.Errorf("Top() = %v, want 0", got)
	}

	if got := s.GetMin(); got != -2 {
		t.Errorf("GetMin() = %v, want -2", got)
	}
}

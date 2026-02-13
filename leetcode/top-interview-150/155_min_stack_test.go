package topinterview150

import "testing"

// 155. 最小栈 (Min Stack)
//
// 题目描述:
// 设计一个支持 push ，pop ，top 操作，并能在常数时间内检索到最小元素的栈。
// 实现 MinStack 类:
// MinStack() 初始化堆栈对象。
// void push(int val) 将元素val推入堆栈。
// void pop() 删除堆栈顶部的元素。
// int top() 获取堆栈顶部的元素。
// int getMin() 获取堆栈中的最小元素。
//
// 示例 1:
// 输入：
// ["MinStack","push","push","push","getMin","pop","top","getMin"]
// [[],[-2],[0],[-3],[],[],[],[]]
//
// 输出：
// [null,null,null,null,-3,null,0,-2]
//
// 解释：
// MinStack minStack = new MinStack();
// minStack.push(-2);
// minStack.push(0);
// minStack.push(-3);
// minStack.getMin();   --> 返回 -3.
// minStack.pop();
// minStack.top();      --> 返回 0.
// minStack.getMin();   --> 返回 -2.

type MinStack struct {
	stack    []int
	minStack []int
}

func ConstructorMinStack() MinStack {
	return MinStack{
		stack:    []int{},
		minStack: []int{},
	}
}

func (this *MinStack) Push(val int) {
	this.stack = append(this.stack, val)

	// 如果最小栈为空，或者新值 <= 当前最小值，则更新最小栈
	if len(this.minStack) == 0 || val <= this.minStack[len(this.minStack)-1] {
		this.minStack = append(this.minStack, val)
	} else {
		this.minStack = append(this.minStack, this.minStack[len(this.minStack)-1])
	}

}

func (this *MinStack) Pop() {
	if len(this.stack) > 0 {
		this.stack = this.stack[:len(this.stack)-1]
		this.minStack = this.minStack[:len(this.minStack)-1]
	}
}

func (this *MinStack) Top() int {
	return this.stack[len(this.stack)-1]
}

func (this *MinStack) GetMin() int {
	return this.minStack[len(this.minStack)-1]
}

func TestMinStack(t *testing.T) {
	// 简单的测试逻辑
	// minStack := ConstructorMinStack()
	// minStack.Push(-2)
	// minStack.Push(0)
	// minStack.Push(-3)
	// if got := minStack.GetMin(); got != -3 {
	// 	t.Errorf("GetMin() = %v, want %v", got, -3)
	// }
	// minStack.Pop()
	// if got := minStack.Top(); got != 0 {
	// 	t.Errorf("Top() = %v, want %v", got, 0)
	// }
	// if got := minStack.GetMin(); got != -2 {
	// 	t.Errorf("GetMin() = %v, want %v", got, -2)
	// }
}

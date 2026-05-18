package top100liked

import (
	"reflect"
	"testing"
)

// 199. 二叉树的右视图 (Binary Tree Right Side View)
//
// 题目描述:
// 给定一个二叉树的根节点 root，想象自己站在它的右侧，
// 按照从顶部到底部的顺序，返回从右侧所能看到的节点值。
//
// 示例 1：
// 输入：root = [1,2,3,null,5,null,4]
// 输出：[1,3,4]
//
// 示例 2：
// 输入：root = [1,null,3]
// 输出：[1,3]
//
// 提示：BFS 层序遍历，取每层最后一个节点；或 DFS 先访问右子树
// BFS
func rightSideView(root *TreeNode) []int {
	if root == nil {
		return nil
	}
	res := make([]int, 0)
	queue := make([]*TreeNode, 0)
	queue = append(queue, root)
	for len(queue) > 0 {
		size := len(queue)
		res = append(res, queue[size-1].Val)
		for i := 0; i < size; i++ {
			if queue[i].Left != nil {
				queue = append(queue, queue[i].Left)
			}
			if queue[i].Right != nil {
				queue = append(queue, queue[i].Right)
			}
		}
		queue = queue[size:]
	}
	return res
}

// DFS
func rightSideView1(root *TreeNode) []int {
	var res []int
	var dfs func(*TreeNode, int)
	dfs = func(node *TreeNode, depth int) {
		if node == nil {
			return
		}
		if len(res) == depth {
			res = append(res, node.Val)
		}
		dfs(node.Right, depth+1)
		dfs(node.Left, depth+1)
	}
	dfs(root, 0)
	return res
}

func TestRightSideView(t *testing.T) {
	tests := []struct {
		name     string
		root     *TreeNode
		expected []int
	}{
		{
			name: "示例1",
			root: &TreeNode{Val: 1,
				Left:  &TreeNode{Val: 2, Right: &TreeNode{Val: 5}},
				Right: &TreeNode{Val: 3, Right: &TreeNode{Val: 4}},
			},
			expected: []int{1, 3, 4},
		},
		{
			name:     "只有右子树",
			root:     &TreeNode{Val: 1, Right: &TreeNode{Val: 3}},
			expected: []int{1, 3},
		},
		{
			name:     "空树",
			root:     nil,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rightSideView(tt.root)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("rightSideView() = %v, want %v", got, tt.expected)
			}
		})

		t.Run(tt.name, func(t *testing.T) {
			got := rightSideView1(tt.root)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("rightSideView1() = %v, want %v", got, tt.expected)
			}
		})
	}
}

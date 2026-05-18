package top100liked

import "testing"

// 337. 打家劫舍 III (House Robber III)
//
// 题目描述:
// 小偷发现了一个新的可行窃的地区，这个地区的所有房屋的排列类似于一棵二叉树。
// 如果两个直接相连的房子在同一天晚上被打劫，房屋将自动报警。
// 计算在不触动警报的情况下，小偷一晚能够盗取的最高金额。
//
// 示例 1：
// 输入：root = [3,2,3,null,3,null,1]
// 输出：7
// 解释：偷 3 + 3 + 1 = 7
//
// 示例 2：
// 输入：root = [3,4,5,1,3,null,1]
// 输出：9
// 解释：偷 4 + 5 = 9
//
// 提示：树形 DP，每个节点返回 [不偷当前节点的最大值, 偷当前节点的最大值]

func rob3(root *TreeNode) int {
	a, b := dfsRob3(root)
	return max(a, b)
}

func dfsRob3(node *TreeNode) (int, int) {
	if node == nil {
		return 0, 0
	}
	leftNoRob, leftRob := dfsRob3(node.Left)
	rightNoRob, rightRob := dfsRob3(node.Right)

	// 偷当前节点 = 不偷左 + 不偷右 + 偷当前
	curRob := leftNoRob + rightNoRob + node.Val
	// 不偷当前节点 = max(偷左) + max(偷右)
	curNoRob := max(leftNoRob, leftRob) + max(rightNoRob, rightRob)

	return curNoRob, curRob
}

func TestRob3(t *testing.T) {
	tests := []struct {
		name     string
		root     *TreeNode
		expected int
	}{
		{
			name: "示例1",
			root: &TreeNode{Val: 3,
				Left:  &TreeNode{Val: 2, Right: &TreeNode{Val: 3}},
				Right: &TreeNode{Val: 3, Right: &TreeNode{Val: 1}},
			},
			expected: 7,
		},
		{
			name: "示例2",
			root: &TreeNode{Val: 3,
				Left:  &TreeNode{Val: 4, Left: &TreeNode{Val: 1}, Right: &TreeNode{Val: 3}},
				Right: &TreeNode{Val: 5, Right: &TreeNode{Val: 1}},
			},
			expected: 9,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rob3(tt.root)
			if got != tt.expected {
				t.Errorf("rob3() = %v, want %v", got, tt.expected)
			}
		})
	}
}

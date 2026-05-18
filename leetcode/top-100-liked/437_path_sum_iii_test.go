package top100liked

import "testing"

// 437. 路径总和 III (Path Sum III)
//
// 题目描述:
// 给定一个二叉树的根节点 root 和一个整数 targetSum，
// 求该二叉树里节点值之和等于 targetSum 的路径的数目。
// 路径不需要从根节点开始，也不需要在叶子节点结束，但是路径方向必须是向下的。
//
// 示例：
// 输入：root = [10,5,-3,3,2,null,11,3,-2,null,1], targetSum = 8
// 输出：3
// 解释：和等于 8 的路径有 3 条：5->3, 5->2->1, -3->11
//
// 提示：前缀和 + 哈希表，类似"和为K的子数组"在树上的应用

func pathSum(root *TreeNode, targetSum int) int {
	prefix := map[int64]int{0: 1}
	var res int
	var dfs func(node *TreeNode, cur int64)
	dfs = func(node *TreeNode, cur int64) {
		if node == nil {
			return
		}
		cur += int64(node.Val)
		res += prefix[cur-int64(targetSum)]
		prefix[cur]++
		dfs(node.Left, cur)
		dfs(node.Right, cur)
		prefix[cur]--
	}
	dfs(root, 0)
	return res
}

func TestPathSum(t *testing.T) {
	tests := []struct {
		name      string
		root      *TreeNode
		targetSum int
		expected  int
	}{
		{
			name: "示例",
			root: &TreeNode{Val: 10,
				Left: &TreeNode{Val: 5,
					Left:  &TreeNode{Val: 3, Left: &TreeNode{Val: 3}, Right: &TreeNode{Val: -2}},
					Right: &TreeNode{Val: 2, Right: &TreeNode{Val: 1}},
				},
				Right: &TreeNode{Val: -3, Right: &TreeNode{Val: 11}},
			},
			targetSum: 8,
			expected:  3,
		},
		{
			name:      "单节点匹配",
			root:      &TreeNode{Val: 5},
			targetSum: 5,
			expected:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pathSum(tt.root, tt.targetSum)
			if got != tt.expected {
				t.Errorf("pathSum() = %v, want %v", got, tt.expected)
			}
		})
	}
}

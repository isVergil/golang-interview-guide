package top100liked

import "testing"

// 538. 把二叉搜索树转换为累加树 (Convert BST to Greater Tree)
//
// 题目描述:
// 给出二叉搜索树的根节点，该树的节点值各不相同，
// 请你将其转换为累加树（Greater Sum Tree），使每个节点 node 的新值等于
// 原树中大于或等于 node.val 的值之和。
//
// 示例：
// 输入：root = [4,1,6,0,2,5,7,null,null,null,3,null,null,null,8]
// 输出：[30,36,21,36,35,26,15,...]
//
// 提示：反向中序遍历（右→根→左），用累加变量记录已遍历节点的和

func convertBST(root *TreeNode) *TreeNode {
	sum := 0
	var dfs func(node *TreeNode)
	dfs = func(node *TreeNode) {
		if node == nil {
			return
		}
		dfs(node.Right) // 先访问比当前节点大的所有节点
		sum += node.Val // 累加自身
		node.Val = sum  // 更新为累加值
		dfs(node.Left)  // 再访问比当前节点小的所有节点
	}
	dfs(root)
	return root
}

func TestConvertBST(t *testing.T) {
	// 输入：[4, 1, 6, 0, 2, 5, 7]
	root := &TreeNode{Val: 4,
		Left:  &TreeNode{Val: 1, Left: &TreeNode{Val: 0}, Right: &TreeNode{Val: 2}},
		Right: &TreeNode{Val: 6, Left: &TreeNode{Val: 5}, Right: &TreeNode{Val: 7}},
	}

	result := convertBST(root)

	// 验证根节点：4 + 5 + 6 + 7 = 22
	if result != nil && result.Val != 22 {
		t.Errorf("root.Val = %v, want 22", result.Val)
	}
	// 验证右子节点：6 + 7 = 13
	if result != nil && result.Right != nil && result.Right.Val != 13 {
		t.Errorf("root.Right.Val = %v, want 13", result.Right.Val)
	}
}

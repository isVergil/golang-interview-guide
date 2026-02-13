package topinterview150

import (
	"strings"
	"testing"
)

// 71. 简化路径 (Simplify Path)
//
// 题目描述:
// 给你一个字符串 path ，表示指向某一文件或目录的 Unix 风格 绝对路径 （以 '/' 开头），请你将其转化为更加简洁的规范路径。
// 在 Unix 风格的文件系统中，一个点（.）表示当前目录本身；此外，两个点 （..） 表示将目录切换到上一级（指向父目录）；两者都可以是复杂相对路径的组成部分。任意多个连续的斜杠（即，'//'）都被视为单个斜杠 '/' 。 对于此问题，任何其他格式的点（例如，'...'）均被视为文件/目录名称。
// 请注意，返回的 规范路径 必须遵循下述格式：
// 1. 始终以斜杠 '/' 开头。
// 2. 两个目录名之间必须只有一个斜杠 '/' 。
// 3. 最后一个目录名（如果存在）不能 以 '/' 结尾。
// 4. 此外，路径仅包含从根目录到目标文件或目录的路径上的目录（即，不含 '.' 或 '..'）。
// 返回简化后得到的 规范路径 。
//
// 示例 1：
// 输入：path = "/home/"
// 输出："/home"
// 解释：注意，最后一个目录名后面没有斜杠。
//
// 示例 2：
// 输入：path = "/../"
// 输出："/"
// 解释：从根目录向上一级是不可行的，因为根目录是你可以到达的最高级。
//
// 示例 3：
// 输入：path = "/home//foo/"
// 输出："/home/foo"
// 解释：在规范路径中，多个连续斜杠需要用一个斜杠替换。

func simplifyPath(path string) string {
	parts := strings.Split(path, "/")
	stack := []string{}
	for _, part := range parts {
		if part == "" || part == "." {
			continue
		}
		if part == ".." {
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}
		} else {
			stack = append(stack, part)
		}
	}
	return "/" + strings.Join(stack, "/")
}

func TestSimplifyPath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "Example 1",
			path:     "/home/",
			expected: "/home",
		},
		{
			name:     "Example 2",
			path:     "/../",
			expected: "/",
		},
		{
			name:     "Example 3",
			path:     "/home//foo/",
			expected: "/home/foo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: Uncomment when implementation is ready
			// if got := simplifyPath(tt.path); got != tt.expected {
			// 	t.Errorf("simplifyPath() = %v, want %v", got, tt.expected)
			// }
		})
	}
}

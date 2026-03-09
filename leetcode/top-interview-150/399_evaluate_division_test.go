package topinterview150

import (
	"testing"
)

// 399. 除法评估 (Evaluate Division)
//
// 题目描述:
// 给你一个变量对数组 equations 和一个实数值数组 values 作为已知条件，其中 equations[i] = [Ai, Bi] 和 values[i] 共同表示等式 Ai / Bi = values[i] 。
// 每个 Ai 或 Bi 是一个表示单个变量的字符串。
// 另有一些查询 queries ，其中 queries[j] = [Cj, Dj] 表示第 j 个查询，你需要根据已知条件找出 Cj / Dj = ? 的结果返回。
// 返回 所有查询的结果 。如果无法确定结果，则返回 -1.0 。
//
// 示例 1：
// 输入：equations = [["a","b"],["b","c"]], values = [2.0,3.0], queries = [["a","c"],["b","a"],["a","e"],["a","a"],["x","x"]]
// 输出：[6.00000,0.50000,-1.00000,1.00000,-1.00000]
//
// 示例 2：
// 输入：equations = [["a","b"],["b","c"],["bc","cd"]], values = [1.5,2.5,5.0], queries = [["a","c"],["c","b"],["bc","e"],["a","a"],["x","x"]]
// 输出：[3.75000,0.40000,-1.00000,1.00000,-1.00000]

func calcEquation(equations [][]string, values []float64, queries [][]string) []float64 {
	panic("not implemented")
}

func TestCalcEquation(t *testing.T) {
	// 除法评估测试
}

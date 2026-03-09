package topinterview150

import (
	"testing"
)

// 909. 蛇梯棋 (Snakes and Ladders)
//
// 题目描述:
// 给你一个大小为 n x n 的整数矩阵 board ，其中 board[i][j] 表示棋盘格子的值。
// 棋盘格子的编号从 1 到 n^2 ，按 蛇形方式 编号，从左下角开始（即 board[n-1][0]），每一行交替方向。
// 玩家从格子 1 （位于 board[n-1][0]）开始。每一回合，玩家可以移动 1 到 6 个格子。
// 如果玩家到达的格子是一个蛇或梯子（即 board[i][j] != -1），则玩家必须移动到该蛇或梯子的目的地。
// 返回达到格子 n^2 所需的最少移动次数。如果无法到达，返回 -1 。
//
// 示例 1：
// 输入：board = [[-1,-1,-1,-1,-1,-1],[-1,-1,-1,-1,-1,-1],[-1,-1,-1,-1,-1,-1],[-1,35,-1,-1,13,-1],[-1,-1,-1,-1,-1,-1],[-1,15,-1,-1,-1,-1]]
// 输出：4

func snakesAndLadders(board [][]int) int {
	panic("not implemented")
}

func TestSnakesAndLadders(t *testing.T) {
	// 蛇梯棋测试
}

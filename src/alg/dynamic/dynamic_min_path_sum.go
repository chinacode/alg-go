package dynamic

import (
	"fmt"
	"util"
)

func MinPathSum2(grid [][]int) int {
	l := len(grid)
	if l < 1 {
		return 0
	}
	dp := make([][]int, l)
	for i, arr := range grid {
		dp[i] = make([]int, len(arr))
	}
	dp[0][0] = grid[0][0]
	for i := 0; i < l; i++ {
		for j := 0; j < len(grid[i]); j++ {
			if i == 0 && j != 0 {
				dp[i][j] = dp[i][j-1] + grid[i][j]
			} else if j == 0 && i != 0 {
				dp[i][j] = dp[i-1][j] + grid[i][j]
			} else if i != 0 && j != 0 {
				dp[i][j] = util.Min(dp[i-1][j], dp[i][j-1]) + grid[i][j]
			}
		}
	}
	return dp[l-1][len(dp[l-1])-1]
}

func MinPathSum(longArr [][]int) int {
	level := len(longArr)
	if level == 1 {
		return longArr[0][0]
	}

	dp := make([][]int, level)
	for i, arr := range longArr {
		dp[i] = make([]int, len(arr))
	}

	//fmt.Println(dp)
	for i := 0; i < len(longArr[0]); i++ {
		if i == 0 {
			dp[0][0] = longArr[0][0]
		} else {
			dp[0][i] = dp[0][i-1] + longArr[0][i]
		}
	}

	for l := 1; l < level; l++ {
		levelLen := len(longArr[l])
		for i := 0; i < levelLen; i++ {
			//fmt.Println(l, i)
			if i == 0 {
				dp[l][i] = util.Min(dp[l-1][i], dp[l-1][len(dp[l-1])-1])
			} else {
				dp[l][i] = util.Min(dp[l-1][i], dp[l][i-1])
			}

			dp[l][i] = dp[l][i] + longArr[l][i]
		}
	}
	//util.PrintLevelArray(dp)

	return dp[level-1][len(dp[level-1])-1]
}

/**
 */
func TestMinPathSum() {
	//给定一个三角形，找出自顶向下的最小路径和。每一步只能移动到下一行中相邻的结点上。
	/**
	输入:
	[
	  [1,3,1],
	  [1,5,1],
	  [4,2,1]
	]
	输出: 7
	解释: 因为路径 1→3→1→1→1 的总和最小。
	*/
	longArr := util.InitRandLevelArrayRange(3, 4, 1, 10)
	util.PrintLevelArray(longArr)

	var resultData int
	loopCount := 1
	loopCount = 5000000
	util.Start("first", "")
	for i := 0; i < loopCount; i++ {
		resultData = MinPathSum(longArr)
	}
	fmt.Println(resultData)
	util.Cut("first", "")

	util.Start("second", "")
	for i := 0; i < loopCount; i++ {
		resultData = MinPathSum2(longArr)
	}
	fmt.Println(resultData)
	util.Cut("second", "")
}

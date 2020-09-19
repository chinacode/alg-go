package dynamic

import (
	"fmt"
	"util"
)

func MinTriangleSum2(triangle [][]int) int {
	if len(triangle) < 1 {
		return 0
	}
	if len(triangle) == 1 {
		return triangle[0][0]
	}
	dp := make([][]int, len(triangle))
	for i, arr := range triangle {
		dp[i] = make([]int, len(arr))
	}
	result := 1<<31 - 1
	dp[0][0] = triangle[0][0]
	dp[1][1] = triangle[1][1] + triangle[0][0]
	dp[1][0] = triangle[1][0] + triangle[0][0]

	for i := 2; i < len(triangle); i++ {
		for j := 0; j < len(triangle[i]); j++ {
			if j == 0 {
				dp[i][j] = dp[i-1][j] + triangle[i][j]
			} else if j == (len(triangle[i]) - 1) {
				dp[i][j] = dp[i-1][j-1] + triangle[i][j]
			} else {
				dp[i][j] = min(dp[i-1][j-1], dp[i-1][j]) + triangle[i][j]
			}
		}
	}
	for _, k := range dp[len(dp)-1] {
		result = min(result, k)
	}
	return result
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

/**
求取最小值
*/
func MinTriangleSum(longArr [][]int) int {
	if len(longArr) == 1 {
		return longArr[0][0]
	}
	sum := 0
	for _, levelArr := range longArr {
		min := levelArr[0]
		for i := 1; i < len(levelArr); i++ {
			if min > levelArr[i] {
				min = levelArr[i]
			}
		}
		sum += min
	}
	return sum
}

func initTriangleArr(level int) [][]int {
	levelArray := make([][]int, level)
	for i := 0; i < level; i++ {
		levelArray[i] = util.InitNoRepeatRandArrayRange(i+1, 1, 10)
	}
	return levelArray
}

/**
 */
func TestMinTriangleSum() {
	//给定一个三角形，找出自顶向下的最小路径和。每一步只能移动到下一行中相邻的结点上。
	/**
	[
	     [2],
	    [3,4],
	   [6,5,7],
	  [4,1,8,3]
	]
	*/
	//longArr := util.InitRandArrayRange(15, 0, 10)
	longArr := initTriangleArr(5)
	fmt.Println(longArr)

	var resultData int
	loopCount := 1
	//loopCount = 10000000
	util.Start("first", "")
	for i := 0; i < loopCount; i++ {
		resultData = MinTriangleSum(longArr)
	}
	fmt.Println(resultData)
	util.Cut("first", "")

	util.Start("second", "")
	for i := 0; i < loopCount; i++ {
		resultData = MinTriangleSum2(longArr)
	}
	fmt.Println(resultData)
	util.Cut("second", "")

}

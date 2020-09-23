package dynamic

import (
	"fmt"
	"util"
)

func MinPathSum2(triangle [][]int) int {
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
	//fmt.Println(dp)
	for _, k := range dp[len(dp)-1] {
		result = util.Min(result, k)
	}
	return result
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
	sum := 1<<31 - 1
	for _, v := range dp[level-1] {
		if v < sum {
			sum = v
		}
	}
	return sum
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
	longArr := util.InitRandLevelArrayRange(3, 3, 1, 10)
	util.PrintLevelArray(longArr)

	var resultData int
	loopCount := 1
	//loopCount = 5000000
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

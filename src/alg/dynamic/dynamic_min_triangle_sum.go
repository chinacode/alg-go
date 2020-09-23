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
	dp[1][0] = triangle[1][0] + triangle[0][0]
	dp[1][1] = triangle[1][1] + triangle[0][0]

	for i := 2; i < len(triangle); i++ {
		for j := 0; j < len(triangle[i]); j++ {
			if j == 0 {
				dp[i][j] = dp[i-1][j] + triangle[i][j]
			} else if j == (len(triangle[i]) - 1) {
				dp[i][j] = dp[i-1][j-1] + triangle[i][j]
			} else {
				dp[i][j] = util.Min(dp[i-1][j-1], dp[i-1][j]) + triangle[i][j]
			}
		}
	}
	//fmt.Println(dp)
	for _, k := range dp[len(dp)-1] {
		result = util.Min(result, k)
	}
	return result
}

/**
求取最小值
*/
func MinTriangleSum3(longArr [][]int) int {
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

func MinTriangleSum(longArr [][]int) int {
	level := len(longArr)
	if level == 1 {
		return longArr[0][0]
	}

	dp := make([][]int, level)
	for i, arr := range longArr {
		dp[i] = make([]int, len(arr))
	}

	dp[0][0] = longArr[0][0]
	dp[1][0] = longArr[0][0] + longArr[1][0]
	dp[1][1] = longArr[0][0] + longArr[1][1]

	for l := 2; l < level; l++ {
		levelLen := len(longArr[l])
		for i := 0; i < levelLen; i++ {
			if 0 == i {
				dp[l][i] = dp[l-1][i] + longArr[l][i]
			} else if i == levelLen-1 {
				dp[l][i] = dp[l-1][i-1] + longArr[l][i]
			} else {
				dp[l][i] = util.Min(dp[l-1][i-1], dp[l-1][i]) + longArr[l][i]
			}
		}
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
	相邻结点定义 2 [3,4] 5 [1,8]
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
	loopCount = 5000000
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

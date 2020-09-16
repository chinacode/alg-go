package dynamic

import (
	"fmt"
	"util"
)

/**
状态转移方程
dp[n]=dp[n-1]+dp[n-2]
*/
func ArrayMaxSubArray(n int) int {
	if 1 == n {
		return 1
	}
	pre1 := 1
	pre2 := 2
	stepN := 0
	for i := 3; i <= n; i++ {
		stepN = pre1 + pre2
		pre1 = pre2
		pre2 = stepN
	}
	return stepN
}

/**
 */
func ArrayMaxSubArray2(n int) int {
	if n == 1 {
		return 1
	}
	dp := make([]int, n+1)
	dp[1] = 1
	dp[2] = 2
	for i := 3; i <= n; i++ {
		dp[i] = dp[i-1] + dp[i-2]
	}
	return dp[n]
}

/**
 */
func TestArrayMaxSubArray() {
	//给定一个整数数组 nums ，找到一个具有最大和的连续子数组（子数组最少包含一个元素），返回其最大和。

	longArr := util.InitRandArray(10, 0, 50)
	var resultData int
	loopCount := 1
	stairsN := 10
	//loopCount = 50000000
	util.Start("first", "")
	for i := 0; i < loopCount; i++ {
		resultData = ArrayMaxSubArray(stairsN)
	}
	fmt.Println(resultData)
	util.Cut("first", "")

	util.Start("second", "")
	for i := 0; i < loopCount; i++ {
		resultData = ArrayMaxSubArray2(stairsN)
	}
	fmt.Println(resultData)
	util.Cut("second", "")

}

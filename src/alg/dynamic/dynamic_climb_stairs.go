package dynamic

import (
	"fmt"
	"util"
)

/**
状态转移方程
dp[n]=dp[n-1]+dp[n-2]
*/
func climbStairs(n int) int {
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
func climbStairs2(n int) int {
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
func TestClimbStairs() {
	//假设你正在爬楼梯。需要 n 阶你才能到达楼顶。每次你可以爬 1 或 2 个台阶。你有多少种不同的方法可以爬到楼顶呢？ **注意：**给定 n 是一个正整数。
	var resultData int
	loopCount := 1
	stairsN := 10
	loopCount = 50000000
	util.Start("first", "")
	for i := 0; i < loopCount; i++ {
		resultData = climbStairs(stairsN)
	}
	fmt.Println(resultData)
	util.Cut("first", "")

	util.Start("second", "")
	for i := 0; i < loopCount; i++ {
		resultData = climbStairs2(stairsN)
	}
	fmt.Println(resultData)
	util.Cut("second", "")

}

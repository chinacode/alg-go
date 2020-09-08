package alg

import (
	"fmt"
	"util"
)

func theBestTimeForStocks2(prices []int) int {
	if len(prices) < 2 {
		return 0
	}
	dp := make([][2]int, len(prices))
	dp[0][0] = 0
	dp[0][1] = -prices[0]
	for i := 1; i < len(prices); i++ {
		dp[i][0] = max(dp[i-1][0], dp[i-1][1]+prices[i])
		dp[i][1] = max(dp[i-1][0]-prices[i], dp[i-1][1])
	}
	return dp[len(prices)-1][0]
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func theBestTimeForStocks(tmpList []int) int {
	total := 0
	if len(tmpList) <= 1 {
		return total
	}

	lastDay := 0
	for _, price := range tmpList {
		if price > lastDay && lastDay != 0 {
			total += price - lastDay
		}

		lastDay = price
	}
	return total
}

func TestTheBestTimeForStocks() {
	tmpList := []int{
		//7, 1, 5, 3, 2, 4,
		//7, 6, 4, 3, 1,
		1, 2, 3, 4, 5,
	}

	resultData := 0
	loopCount := 1
	loopCount = 50000000
	util.Start("first", "")
	for i := 0; i < loopCount; i++ {
		resultData = theBestTimeForStocks(tmpList)
	}
	fmt.Println(resultData)
	util.Cut("first", "")

	util.Start("second", "")
	for i := 0; i < loopCount; i++ {
		resultData = theBestTimeForStocks2(tmpList)
	}
	fmt.Println(resultData)
	util.Cut("second", "")

}

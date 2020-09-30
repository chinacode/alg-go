package dynamic

import (
	"fmt"
	"util"
)

func maxMultipe(sum int, num int) int {
	return sum / num
}

func MinCoinChange2(nums []int, sum int) []int {
	if len(nums) < 1 {
		return nums
	}

	return nums
}

func MinCoinChange(nums []int, sum int) []int {
	if len(nums) < 1 {
		return nums
	}
	dp := make([][]int, len(nums))
	for i := range nums {
		dp[i] = make([]int, sum+1)
	}

	for i, coin := range nums {
		for j := 1; j <= sum; j++ {
			if i == 0 {
				dp[i][j] = j / coin
				continue
			}

			if j < coin {
				dp[i][j] = dp[i-1][j]
			} else {
				dp[i][j] = util.Min(dp[i-1][j], 1+dp[i][j-coin])
			}
		}
	}
	util.PrintLevelArray(dp)

	return nums
}

func MinCoinChangeError(nums []int, sum int) []int {
	if len(nums) < 1 {
		return nums
	}
	for i := len(nums) - 1; i >= 0; i-- {
		if nums[i] < 0 || nums[i] > 0 {
			continue
		}
		nums[i] = maxMultipe(sum, i)
		sum = sum - nums[i]*i
		if sum == 0 {
			break
		}
		//return MinCoinChange(nums, sum)
	}

	total := 0
	for _, num := range nums {
		if num <= 0 {
			continue
		}
		total = total + num
	}
	fmt.Println(total)
	return nums
}

/**
 */
func TestMinCoinChange() {
	/**
	给定4种面额的硬币1分，2分，5分，6分，如果要找11分的零钱，怎么做才能使得找的硬币数量总和最少。
	*/
	sum := 11
	longArr := []int{1, 2, 5, 6}
	//fmt.Println(longArr)

	var resultData []int
	loopCount := 1
	//loopCount = 5000000
	util.Start("first", "")
	for i := 0; i < loopCount; i++ {
		resultData = MinCoinChange(longArr, sum)
	}
	fmt.Println(resultData)
	util.Cut("first", "")

	util.Start("second", "")
	for i := 0; i < loopCount; i++ {
		resultData = MinCoinChange2(longArr, sum)
	}
	fmt.Println(resultData)
	util.Cut("second", "")

}

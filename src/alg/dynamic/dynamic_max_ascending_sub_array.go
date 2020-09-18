package dynamic

import (
	"fmt"
	"util"
)

func MaxAscendingSubArray2(nums []int) ([]int, int) {
	if len(nums) < 1 {
		return nums, 0
	}
	dp := make([]int, len(nums))
	result := 1
	for i := 0; i < len(nums); i++ {
		dp[i] = 1
		for j := 0; j < i; j++ {
			if nums[j] < nums[i] {
				dp[i] = max(dp[j]+1, dp[i])
			}
		}
		result = max(result, dp[i])
	}

	return nums, result
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

/**
 */
func MaxAscendingSubArray(longArr []int) ([]int, int) {
	maxI := 0
	maxLength := 0
	tmpMaxLength := 0
	for i := 1; i < len(longArr); i++ {
		if longArr[i-1] < longArr[i] {
			tmpMaxLength++
		} else {
			tmpMaxLength = 1
		}
		if maxLength < tmpMaxLength {
			maxI = i
			maxLength = tmpMaxLength
		}
	}
	return longArr[maxI-maxLength+1 : maxI+1], maxLength
}

/**
 */
func TestMaxAscendingSubArray() {
	//给定一个无序的整数数组，找到其中最长上升子序列的长度。(不要求连续)
	//longArr := util.InitRandArrayRange(15, 0, 10)
	longArr := []int{10, 9, 2, 5, 3, 7, 101, 18}
	fmt.Println(longArr)

	var resultData int
	var resultDataArr []int
	loopCount := 1
	//loopCount = 10000000
	util.Start("first", "")
	for i := 0; i < loopCount; i++ {
		resultDataArr, resultData = MaxAscendingSubArray(longArr)
	}
	fmt.Println(resultData)
	fmt.Println(resultDataArr)
	util.Cut("first", "")

	util.Start("second", "")
	for i := 0; i < loopCount; i++ {
		resultDataArr, resultData = MaxAscendingSubArray2(longArr)
	}
	fmt.Println(resultData)
	fmt.Println(resultDataArr)
	util.Cut("second", "")

}

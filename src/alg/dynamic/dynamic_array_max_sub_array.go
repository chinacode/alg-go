package dynamic

import (
	"fmt"
	"util"
)

func ArrayMaxSubArray2(nums []int) ([]int, int) {
	if len(nums) < 1 {
		return nums, nums[0]
	}
	dp := make([]int, len(nums))
	//设置初始化值
	dp[0] = nums[0]
	for i := 1; i < len(nums); i++ {
		//处理 dp[i-1] < 0 的情况
		if dp[i-1] < 0 {
			dp[i] = nums[i]
		} else {
			dp[i] = dp[i-1] + nums[i]
		}
	}

	index := 0
	result := -1 << 31 //复制最小整数
	for i, k := range dp {
		if result < k {
			index = i
			result = k
		}
		//result = max(result, k)
	}
	starIndex := 0
	for i := index; i >= 0; i-- {
		if dp[i] < 0 {
			starIndex = i
			break
		}
	}
	if nums[starIndex] < 0 {
		starIndex++
	}

	return nums[starIndex : index+1], result
}

/**
 */
func ArrayMaxSubArray(longArr []int) ([]int, int) {
	dp := make([]int, len(longArr))

	maxI := 0
	maxValue := -1 << 31
	dp[0] = longArr[0]
	for i := 1; i < len(longArr); i++ {
		if dp[i-1] < 0 {
			dp[i] = longArr[i]
		} else {
			dp[i] = dp[i-1] + longArr[i]
		}
		if dp[i] >= maxValue {
			maxI = i
			maxValue = dp[i]
		}
	}
	starIndex := 0
	for i := maxI; i >= 0; i-- {
		if dp[i] < 0 {
			starIndex = i
			break
		}
	}
	if longArr[starIndex] < 0 {
		starIndex++
	}
	return longArr[starIndex : maxI+1], maxValue
}

/**
 */
func TestArrayMaxSubArray() {
	//给定一个整数数组 nums ，找到一个具有最大和的连续子数组（子数组最少包含一个元素），返回其最大和。

	longArr := util.InitRandArrayRange(15, -5, 5)
	fmt.Println(longArr)

	var resultData int
	var resultDataArr []int
	loopCount := 1
	loopCount = 10000000
	util.Start("first", "")
	for i := 0; i < loopCount; i++ {
		resultDataArr, resultData = ArrayMaxSubArray(longArr)
	}
	fmt.Println(resultData)
	fmt.Println(resultDataArr)
	util.Cut("first", "")

	util.Start("second", "")
	for i := 0; i < loopCount; i++ {
		resultDataArr, resultData = ArrayMaxSubArray2(longArr)
	}
	fmt.Println(resultData)
	fmt.Println(resultDataArr)
	util.Cut("second", "")

}

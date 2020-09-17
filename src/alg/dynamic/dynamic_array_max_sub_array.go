package dynamic

import (
	"fmt"
	"util"
)

/**
 */
func ArrayMaxSubArray(longArr []int, subArr []int) int {
	maxSum := 0

	return maxSum
}

/**
 */
func ArrayMaxSubArray2(longArr []int, subArr []int) int {
	maxSum := 0

	return maxSum
}

/**
 */
func TestArrayMaxSubArray() {
	//给定一个整数数组 nums ，找到一个具有最大和的连续子数组（子数组最少包含一个元素），返回其最大和。

	longArr := util.InitRandArray(10)
	subArr := longArr[2:6]

	var resultData int
	loopCount := 1
	//loopCount = 50000000
	util.Start("first", "")
	for i := 0; i < loopCount; i++ {
		resultData = ArrayMaxSubArray(longArr, subArr)
	}
	fmt.Println(resultData)
	util.Cut("first", "")

	util.Start("second", "")
	for i := 0; i < loopCount; i++ {
		resultData = ArrayMaxSubArray2(longArr, subArr)
	}
	fmt.Println(resultData)
	util.Cut("second", "")

}

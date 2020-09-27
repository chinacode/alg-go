package dynamic

import (
	"fmt"
	"util"
)

func TheftMaxMoney2(nums []int) int {
	if len(nums) < 1 {
		return 0
	}
	if len(nums) == 1 {
		return nums[0]
	}
	if len(nums) == 2 {
		return util.Max(nums[0], nums[1])
	}
	nums[1] = util.Max(nums[0], nums[1])
	for i := 2; i < len(nums); i++ {
		nums[i] = util.Max(nums[i-2]+nums[i], nums[i-1])
	}
	return nums[len(nums)-1]
}

func TheftMaxMoney(longArr []int) int {
	size := len(longArr)
	if 1 == size {
		return longArr[0]
	}
	if 2 == size {
		return max(longArr[0], longArr[1])
	}

	dp := make([]int, size)
	dp[0] = longArr[0]
	dp[1] = util.Max(longArr[0], longArr[1])
	for i := 2; i < size; i++ {
		dp[i] = util.Max(dp[i-1], dp[i-2]+longArr[i])
	}

	//fmt.Println(dp)
	return dp[size-1]
}

func TheftMaxMoneySelf(longArr []int) int {
	size := len(longArr)
	if 1 == size {
		return longArr[0]
	}
	if 2 == size {
		return max(longArr[0], longArr[1])
	}

	longArr[1] = util.Max(longArr[0], longArr[1])
	for i := 2; i < size; i++ {
		longArr[i] = util.Max(longArr[i-1], longArr[i-2]+longArr[i])
	}

	//fmt.Println(dp)
	return longArr[size-1]
}

/**
 */
func TestTheftMaxMoney() {
	/**
	你是一个专业的小偷，计划偷窃沿街的房屋。每间房内都藏有一定的现金，
	影响你偷窃的唯一制约因素就是相邻的房屋装有相互连通的防盗系统，
	如果两间相邻的房屋在同一晚上被小偷闯入，系统会自动报警。

	输入: [2,7,9,3,1]
	输出: 12
	解释: 偷窃 1 号房屋 (金额 = 2), 偷窃 3 号房屋 (金额 = 9)，接着偷窃 5 号房屋 (金额 = 1)。
	     偷窃到的最高金额 = 2 + 9 + 1 = 12 。
	*/
	longArr := util.InitRandArrayRange(5, 0, 10)
	fmt.Println(longArr)

	var resultData int
	loopCount := 1
	//loopCount = 5000000
	util.Start("first", "")
	for i := 0; i < loopCount; i++ {
		resultData = TheftMaxMoney(longArr)
	}
	fmt.Println(resultData)
	util.Cut("first", "")

	longArr1 := longArr
	util.Start("second", "")
	for i := 0; i < loopCount; i++ {
		resultData = TheftMaxMoney2(longArr1)
	}
	fmt.Println(resultData)
	util.Cut("second", "")

	longArr2 := longArr
	util.Start("third", "")
	for i := 0; i < loopCount; i++ {
		resultData = TheftMaxMoneySelf(longArr2)
	}
	fmt.Println(resultData)
	util.Cut("third", "")

}

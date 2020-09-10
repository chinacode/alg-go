package string

import (
	"fmt"
	"util"
)

func getStringZChange2(nums []int, target int) []int {
	result := []int{}
	m := make(map[int]int)
	for i, k := range nums {
		if value, exist := m[target-k]; exist {
			result = append(result, value)
			result = append(result, i)
		}
		m[k] = i
	}
	return result
}

func getStringZChange(tmpStr string, lineNum int) string {
	newStr := ""

	return newStr
}

/**
L     D     R
E   O E   I I
E C   I H   N
T     S     G
*/
func TestGetStringZChange() {
	tmpStr := "LEETCODEISHIRING"

	var resultData string
	//ThreeSum := 0
	lineNum := 4
	loopCount := 1
	loopCount = 3000000
	util.Start("first", "")
	for i := 0; i < loopCount; i++ {
		resultData = getStringZChange(tmpStr, lineNum)
	}
	fmt.Println(resultData)
	util.Cut("first", "")

	util.Start("second", "")
	for i := 0; i < loopCount; i++ {
		//resultData = getStringZChange2(tmpList, ThreeSum)
	}
	fmt.Println(resultData)
	util.Cut("second", "")

}

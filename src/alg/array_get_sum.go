package alg

import (
	"fmt"
	"util"
)

func getSumArray(tmpList []int, sum int) []int {
	newList := []int{0, 0}
	valIndexMap := map[int]int{}

	for i := 0; i < len(tmpList); i++ {
		valIndexMap[tmpList[i]] = i
	}

	for _, v := range tmpList {
		diff := sum - v
		if valIndexMap[diff] != 0 {
			newList = []int{valIndexMap[v], valIndexMap[diff]}
			break
		}
	}

	return newList
}

func TestGetSumArray() {
	tmpList := []int{
		1, 2, 6, 3,
		//7, 6, 4, 3, 1,
		//1, 2, 3, 4, 5,
	}

	var resultData []int
	sum := 5
	loopCount := 1
	loopCount = 10000000
	util.Start("first", "")
	for i := 0; i < loopCount; i++ {
		resultData = getSumArray(tmpList, sum)
	}
	fmt.Println(resultData)
	util.Cut("first", "")

	util.Start("second", "")
	for i := 0; i < loopCount; i++ {
		//resultData = getSumArray2(tmpList)
	}
	fmt.Println(resultData)
	util.Cut("second", "")

}

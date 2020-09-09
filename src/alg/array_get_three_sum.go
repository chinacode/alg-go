package alg

import (
	"fmt"
	"util"
)

func getThreeSumArray2(nums []int, target int) []int {
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

func getThreeSumArray(tmpList []int, threeSum int) [][]int {
	var newList [][]int
	//valIndexMap := map[int]int{}

	//for i := 0; i < len(tmpList); i++ {
	//	//	valIndexMap[tmpList[i]] = i
	//	//}

	//for i, v := range tmpList {
	//	value, exists := valIndexMap[threeSum-v]
	//	if exists {
	//		newList = []int{i, value}
	//	}
	//	valIndexMap[v] = i
	//}

	return newList
}

func TestGetThreeSumArray() {
	tmpList := []int{
		1, 2, 6, 3,
		//7, 6, 4, 3, 1,
		//1, 2, 3, 4, 5,
	}

	var resultData [][]int
	ThreeSum := 7
	loopCount := 1
	loopCount = 3000000
	util.Start("first", "")
	for i := 0; i < loopCount; i++ {
		resultData = getThreeSumArray(tmpList, ThreeSum)
	}
	fmt.Println(resultData)
	util.Cut("first", "")

	util.Start("second", "")
	for i := 0; i < loopCount; i++ {
		//resultData = getThreeSumArray2(tmpList, ThreeSum)
	}
	fmt.Println(resultData)
	util.Cut("second", "")

}

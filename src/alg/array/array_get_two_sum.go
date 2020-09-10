package array

import (
	"fmt"
	"util"
)

func getSumArray2(nums []int, target int) []int {
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

func getSumArray(tmpList []int, sum int) []int {
	newList := []int{0, 0}
	valIndexMap := map[int]int{}

	//for i := 0; i < len(tmpList); i++ {
	//	valIndexMap[tmpList[i]] = i
	//}

	for i, v := range tmpList {
		value, exists := valIndexMap[sum-v]
		if exists {
			newList = []int{i, value}
		}
		valIndexMap[v] = i
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
	sum := 7
	loopCount := 1
	loopCount = 3000000
	util.Start("first", "")
	for i := 0; i < loopCount; i++ {
		resultData = getSumArray(tmpList, sum)
	}
	fmt.Println(resultData)
	util.Cut("first", "")

	util.Start("second", "")
	for i := 0; i < loopCount; i++ {
		resultData = getSumArray2(tmpList, sum)
	}
	fmt.Println(resultData)
	util.Cut("second", "")

}

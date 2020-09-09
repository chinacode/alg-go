package alg

import (
	"fmt"
	"util"
)

func addOneArray2(tmpList []int) []int {
	newList := make([]int, len(tmpList))
	copy(newList, tmpList)

	var result []int
	addon := 0
	for i := len(newList) - 1; i >= 0; i-- {
		newList[i] += addon
		addon = 0
		if i == len(newList)-1 {
			newList[i]++
		}
		if newList[i] == 10 {
			addon = 1
			newList[i] = newList[i] % 10
		}
	}
	if addon == 1 {
		result = make([]int, 1)
		result[0] = 1
		result = append(result, newList...)
	} else {
		result = newList
	}
	return result
}

func addOneArray(tmpList []int) []int {
	listLen := len(tmpList)
	newList := make([]int, listLen)
	copy(newList, tmpList)

	carry := 0
	for i := listLen - 1; i >= 0; i-- {
		if i == listLen-1 {
			newList[i] = newList[i] + 1
		}
		if newList[i]+carry == 10 {
			newList[i] = 0
			carry = 1
		} else {
			newList[i] = newList[i] + carry
			carry = 0
		}
	}
	return newList
}

func TestAddOneArray() {
	tmpList := []int{
		1, 2, 9, 9,
		//7, 6, 4, 3, 1,
		//1, 2, 3, 4, 5,
	}

	var resultData []int
	loopCount := 1
	loopCount = 10000000
	util.Start("first", "")
	for i := 0; i < loopCount; i++ {
		resultData = addOneArray(tmpList)
	}
	fmt.Println(resultData)
	util.Cut("first", "")

	util.Start("second", "")
	for i := 0; i < loopCount; i++ {
		resultData = addOneArray2(tmpList)
	}
	fmt.Println(resultData)
	util.Cut("second", "")

}

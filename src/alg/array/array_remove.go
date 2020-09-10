package array

import (
	"fmt"
	"util"
)

func removeArray(tmpList []int, k int) []int {
	newList := make([]int, len(tmpList))
	copy(newList, tmpList)
	for i := 0; i < len(newList); i++ {
		if newList[i] == k {
			tmpList = append(newList[:i], newList[i+1:]...)
			i--
		}
	}
	return newList
}

func removeSortArrayRepeatItem(tmpList []int) []int {
	last := -1
	for i := 0; i < len(tmpList); i++ {
		if tmpList[i] == last && last >= 0 {
			tmpList = append(tmpList[:i], tmpList[i+1:]...)
			i--
		}

		last = tmpList[i]
	}
	return tmpList
}

func TestRemoveArray() {
	tmpList := []int{
		0, 0, 1, 1, 1, 2, 2, 3, 3, 4,
		//7, 6, 4, 3, 1,
		//1, 2, 3, 4, 5,
	}

	var resultData []int
	removeK := 3
	loopCount := 1
	loopCount = 10000000
	util.Start("first", "")
	for i := 0; i < loopCount; i++ {
		resultData = removeArray(tmpList, removeK)
	}
	fmt.Println(resultData)
	util.Cut("first", "")

	util.Start("second", "")
	for i := 0; i < loopCount; i++ {
		resultData = removeSortArrayRepeatItem(tmpList)
	}
	fmt.Println(resultData)
	util.Cut("second", "")

}

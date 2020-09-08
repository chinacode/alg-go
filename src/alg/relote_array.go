package alg

import (
	"fmt"
	"util"
)

func rotateArray2(nums []int, k int) []int {
	newList := make([]int, len(nums))
	copy(newList, nums)
	reverse(newList)
	reverse(newList[:k%len(newList)])
	reverse(newList[k%len(newList):])
	return newList
}

func reverse(arr []int) {
	for i := 0; i < len(arr)/2; i++ {
		arr[i], arr[len(arr)-i-1] = arr[len(arr)-i-1], arr[i]
	}
}

func rotateArray3(tmpList []int, k int) []int {
	if len(tmpList) <= 1 {
		return tmpList
	}

	listLength := len(tmpList)
	newList := make([]int, listLength)
	for i := 0; i < listLength; i++ {

	}
	return newList
}

func rotateArray(tmpList []int, k int) []int {
	if len(tmpList) <= 1 {
		return tmpList
	}

	listLength := len(tmpList)
	newList := make([]int, listLength)
	for i := 0; i < listLength; i++ {
		trueIndex := listLength - k + i
		if trueIndex > listLength-1 {
			trueIndex = trueIndex - listLength
		}
		//println(trueIndex, tmpList[trueIndex])
		//newList = append(newList, tmpList[trueIndex])
		newList[i] = tmpList[trueIndex]
	}
	return newList
}

func TestRotateArray() {
	tmpList := []int{
		1, 2, 3, 4, 5, 6, 7,
		//7, 6, 4, 3, 1,
		//1, 2, 3, 4, 5,
	}

	var resultData []int
	rotateK := 3
	loopCount := 1
	loopCount = 5000000
	util.Start("first", "")
	for i := 0; i < loopCount; i++ {
		resultData = rotateArray(tmpList, rotateK)
	}
	fmt.Println(resultData)
	util.Cut("first", "")

	util.Start("second", "")
	for i := 0; i < loopCount; i++ {
		resultData = rotateArray2(tmpList, rotateK)
	}
	fmt.Println(resultData)
	util.Cut("second", "")

}

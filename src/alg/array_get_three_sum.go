package alg

import (
	"fmt"
	"sort"
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

func getThreeSumArray(tmpList []int) [][]int {
	var newList [][]int

	sort.Ints(tmpList)
	listLen := len(tmpList)
	for i, v := range tmpList {
		if v > 0 {
			break
		}
		l := i + 1
		r := listLen - 1

		for l < r {
			if v+tmpList[l]+tmpList[r] == 0 {
				for v+tmpList[l]+tmpList[r] < 0 {
					r++
				}
			} else if v+tmpList[l]+tmpList[r] < 0 {
				for v+tmpList[l]+tmpList[r] > 0 {
					r++
				}
			} else {

			}
		}

		newList = append(newList, []int{v, tmpList[l], tmpList[r]})
	}

	return newList
}

func TestGetThreeSumArray() {
	tmpList := []int{
		-1, 0, 1, 2, -1, -4,
		//7, 6, 4, 3, 1,
		//1, 2, 3, 4, 5,
	}

	var resultData [][]int
	//ThreeSum := 0
	loopCount := 1
	//loopCount = 3000000
	util.Start("first", "")
	for i := 0; i < loopCount; i++ {
		resultData = getThreeSumArray(tmpList)
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

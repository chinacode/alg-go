package string

import (
	"fmt"
	"util"
)

func FirstNoRepeatCharacter(tmpStr string) int {
	var arr [26]int
	for i, k := range tmpStr {
		arr[k-'a'] = i
	}
	var strBytes = []byte(tmpStr)
	for i, v := range strBytes {
		if arr[v-'a'] == i {
			return i
		} else {
			arr[v-'a'] = -1
		}
	}

	return -1
}

func FirstNoRepeatCharacter2(s string) int {
	var arr [26]int
	for i, k := range s {
		arr[k-'a'] = i
	}
	for i, k := range s {
		if i == arr[k-'a'] {
			return i
		} else {
			arr[k-'a'] = -1
		}
	}
	return -1
}

/**
 */
func TestFirstNoRepeatCharacter() {
	//给定一个字符串，找到它的第一个不重复的字符，并返回它的索引。如果不存在，则返回 -1 。
	//tmpStr := util.RandStringBytesMask(11)
	//tmpStr := "leetcode"
	tmpStr := "eetcode"

	//fmt.Println(tmpStr)
	var resultData int
	loopCount := 1
	//loopCount = 20000000
	util.Start("first", "")
	for i := 0; i < loopCount; i++ {
		resultData = FirstNoRepeatCharacter(tmpStr)
	}
	fmt.Println(resultData)
	util.Cut("first", "")

	util.Start("second", "")
	for i := 0; i < loopCount; i++ {
		resultData = FirstNoRepeatCharacter2(tmpStr)
	}
	fmt.Println(resultData)
	util.Cut("second", "")

}

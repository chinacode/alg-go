package string

import (
	"fmt"
	"util"
)

func StringReverse(tmpStr string) string {
	var strBytes = []byte(tmpStr)
	stringLen := len(strBytes)
	//var source byte
	//var dist byte
	for i := 0; i < stringLen/2; i++ {
		//source = strBytes[stringLen-i-1]
		//dist = strBytes[i]
		//strBytes[i] = source
		//strBytes[stringLen-i-1] = dist

		strBytes[i], strBytes[stringLen-i-1] = strBytes[stringLen-i-1], strBytes[i]
	}

	return string(strBytes)

}

func StringReverse2(tmpStr string) string {
	var s = []byte(tmpStr)
	left := 0
	right := len(s) - 1
	for left < right {
		s[left], s[right] = s[right], s[left]
		left++
		right--
	}
	return string(s)
}

/**
 */
func TestStringReverse() {
	tmpStr := util.RandStringBytesMask(11)

	fmt.Println(tmpStr)
	var resultData string
	loopCount := 1
	loopCount = 20000000
	util.Start("first", "")
	for i := 0; i < loopCount; i++ {
		resultData = StringReverse(tmpStr)
	}
	fmt.Println(resultData)
	util.Cut("first", "")

	util.Start("second", "")
	for i := 0; i < loopCount; i++ {
		resultData = StringReverse2(tmpStr)
	}
	fmt.Println(resultData)
	util.Cut("second", "")

}

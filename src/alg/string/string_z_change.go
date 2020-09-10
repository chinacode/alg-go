package string

import (
	"fmt"
	"strings"
	"util"
)

/**
 * N 变化 先找规律 每一行 都是按照规则进行取模
 */
func getStringZChange2(tmpStr string, rowNum int) string {
	if rowNum == 1 {
		return tmpStr
	}
	var b = []rune(tmpStr)
	var res = make([]string, rowNum)
	var length = len(b)
	var period = rowNum*2 - 2
	for i := 0; i < length; i++ {
		var mod = i % period
		var str = string(b[i])
		if mod < rowNum {
			res[mod] += str
		} else {
			res[period-mod] += str
		}
	}
	return strings.Join(res, "")
}

func getStringZChange2Arr(tmpStr string, rowNum int) [][]string {
	resultArr := make([][]string, rowNum)
	if rowNum == 1 {
		return resultArr
	}
	var length = len(tmpStr)
	var period = rowNum*2 - 2
	for i := 0; i < length; i++ {
		var mod = i % period
		var str = tmpStr[i : i+1]
		var index = period - mod
		if mod < rowNum {
			index = mod
		}
		resultArr[index] = append(resultArr[index], str)
	}
	return resultArr
}

func getStringZChangeArr(tmpStr string, rowNum int) [][]string {
	resultArr := [][]string{}
	if rowNum == 1 {
		return resultArr
	}
	stringLen := len(tmpStr)

	column := 0
	tmpIndex := rowNum
	period := rowNum - 1
	columnArr := []string{}
	for i := 0; i < stringLen; i++ {
		selStr := tmpStr[i : i+1]
		if column == 0 || column%period == 0 {
			columnArr = append(columnArr, selStr)
		} else {
			columnArr = make([]string, rowNum)
			columnArr[tmpIndex-1] = selStr
		}
		if len(columnArr) == rowNum {
			resultArr = append(resultArr, columnArr)
			columnArr = []string{}
			column++
			tmpIndex--
			if tmpIndex <= 1 {
				tmpIndex = rowNum
			}
		}
	}

	//fmt.Println(resultArr)
	//return newStr
	return resultArr
}

/**
L     D     R
E   O E   I I
E C   I H   N
T     S     G
*/
func TestGetStringZChange() {
	tmpStr := "LEETCODEISHIRINGALSJDLASJDLASJLDJSLDJSA"

	var resultData string
	var resultDataArr [][]string
	//ThreeSum := 0
	rowNum := 8
	loopCount := 1
	loopCount = 200000
	util.Start("first", "")
	for i := 0; i < loopCount; i++ {
		resultData = getStringZChange2(tmpStr, rowNum)
	}
	fmt.Println(resultData)
	util.Cut("first", "")

	util.Start("second", "")
	for i := 0; i < loopCount; i++ {
		//resultData = getStringZChange2(tmpList, ThreeSum)
	}
	fmt.Println(resultData)
	util.Cut("second", "")

	util.Start("first arr", "")
	for i := 0; i < loopCount; i++ {
		resultDataArr = getStringZChangeArr(tmpStr, rowNum)
	}
	fmt.Println(resultDataArr)
	util.Cut("first arr", "")

	util.Start("second arr", "")
	for i := 0; i < loopCount; i++ {
		resultDataArr = getStringZChange2Arr(tmpStr, rowNum)
	}
	fmt.Println(resultDataArr)
	util.Cut("second arr", "")

}

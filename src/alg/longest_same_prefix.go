package alg

import (
	"fmt"
	"strings"
	"util"
)

func longestSamePrefix(strList []string) string {
	if len(strList) <= 1 {
		return ""
	}
	i := 0
	var tmpStr string
	for i = 0; i < len(strList[0]); i++ {
		for _, str := range strList {
			if str[i:i+1] != tmpStr && tmpStr != "" {
				return strList[0][:i]
			}
			tmpStr = str[i : i+1]
		}
		tmpStr = ""
	}
	return strList[0][:i]
}

func longestCommonPrefix(strs []string) string {
	if len(strs) < 1 {
		return ""
	}
	prefix := strs[0]
	for _, k := range strs {
		for strings.Index(k, prefix) != 0 {
			if len(prefix) == 0 {
				return ""
			}
			prefix = prefix[:len(prefix)-1]
		}
	}
	return prefix
}

func TestLongestSamePrefix() {
	strList := []string{
		"ashjhdbkajkjkja",
		"ashasdsasdfsfsf",
		"asdasdsasdfsfsf",
		"assasdsasdfsfsf",
	}

	samePrefix := ""

	util.Start("first", "")
	for i := 0; i < 50000000; i++ {
		samePrefix = longestSamePrefix(strList)
	}
	fmt.Println(samePrefix)
	util.Cut("first", "")

	util.Start("second", "")
	for i := 0; i < 50000000; i++ {
		samePrefix = longestCommonPrefix(strList)
	}
	fmt.Println(samePrefix)
	util.Cut("second", "")

}

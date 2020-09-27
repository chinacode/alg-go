package util

import (
	"config"
	"fmt"
	"math/rand"
	"sort"
)

func Min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func PrintLevelArray(levelArr [][]int) {
	for level := 0; level < len(levelArr); level++ {
		if level == 0 {
			fmt.Println("[")
		}

		fmt.Println("	", levelArr[level])

		if level == len(levelArr)-1 {
			fmt.Println("]")
		}
	}
}

func InitListNodes(linkNum int) *config.ListNode {
	headNode := &config.ListNode{}
	var preNode *config.ListNode
	for i := 1; i <= linkNum; i++ {
		curNote := config.ListNode{Val: i, Next: nil}
		if i == 1 {
			headNode = &curNote
		} else {
			preNode.Next = &curNote
		}
		preNode = &curNote
	}
	return headNode
}

func GenerateRangeNum(min, max int) int {
	//rand.Seed(time.Now().Unix())
	randNum := rand.Intn(max-min) + min
	return randNum
}

func InitRangeLinkList(linkNum int, min int, max int) *config.ListNode {
	headNode := &config.ListNode{}
	if min >= max {
		return headNode
	}
	var preNode *config.ListNode
	for i := 1; i <= linkNum; i++ {
		v := GenerateRangeNum(min, max)

		curNote := config.ListNode{Val: v, Next: nil}
		if i == 1 {
			headNode = &curNote
		} else {
			preNode.Next = &curNote
		}
		preNode = &curNote
	}
	return headNode
}

func InitListRandNodes(linkNum int) *config.ListNode {
	return InitRangeLinkList(linkNum, 10, 50)
}

func InitListSortRandNodes(linkNum int) *config.ListNode {
	headNode := &config.ListNode{}
	var preNode *config.ListNode

	list := make([]int, linkNum)
	for i := 0; i < linkNum; i++ {
		v := GenerateRangeNum(10, 50)
		list[i] = v
	}
	sort.Ints(list)
	for i, v := range list {
		curNote := config.ListNode{Val: v, Next: nil}
		if i == 0 {
			headNode = &curNote
		} else {
			preNode.Next = &curNote
		}
		preNode = &curNote
	}
	return headNode
}

func InitRandCycleLinkedList(linkNum int, pos int) *config.ListNode {
	headNode := &config.ListNode{}
	var preNode *config.ListNode
	var posNode *config.ListNode

	if linkNum < pos {
		return headNode
	}

	for i := 1; i <= linkNum; i++ {
		v := GenerateRangeNum(10, 50)

		curNote := config.ListNode{Val: v, Next: nil}
		if i == 1 {
			headNode = &curNote
		} else {
			preNode.Next = &curNote
		}

		if i == pos {
			posNode = &curNote
		}
		if i == linkNum {
			curNote.Next = posNode
		}
		preNode = &curNote
	}
	return headNode
}

func ChangeList2ListNode(list []int) *config.ListNode {
	headNode := &config.ListNode{}
	var preNode *config.ListNode

	for i, v := range list {
		curNode := config.ListNode{Val: v, Next: nil}
		if i == 0 {
			headNode = &curNode
		} else {
			preNode.Next = &curNode
		}
		preNode = &curNode
	}

	return headNode
}

func LinkedListPrint(head *config.ListNode) []int {
	if head == nil {
		return nil
	}
	var res []int
	for head != nil {
		res = append(res, head.Val)
		head = head.Next
	}
	//for i, j := 0, len(res)-1; i < j; {
	//	res[i], res[j] = res[j], res[i]
	//	i++
	//	j--
	//}
	return res
}

func InitRandArray(num int) []int {
	return InitRandArrayRange(num, 0, 50)
}

func InitRandArrayRange(num int, min int, max int) []int {
	arr := make([]int, num)
	for i := 0; i < num; i++ {
		v := GenerateRangeNum(min, max)
		arr[i] = v
	}
	return arr
}

func InitNoRepeatRandArrayRange(num int, min int, max int) []int {
	arr := make([]int, num)
	exists := map[int]int{}
	for i := 0; i < num; i++ {
		v := GenerateRangeNum(min, max)
		for exists[v] > 0 {
			v = GenerateRangeNum(min, max)
		}

		exists[v] = 1
		arr[i] = v
	}
	return arr
}

func InitRandLevelArrayRange(level int, num int, min int, max int) [][]int {
	levelArray := make([][]int, level)
	for i := 0; i < level; i++ {
		levelArray[i] = InitNoRepeatRandArrayRange(num, min, max)
	}
	return levelArray
}

///---------------------------------String--------------------------------///
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandStringBytesMask(n int) string {
	b := make([]byte, n)
	for i := 0; i < n; {
		if idx := int(rand.Int63() & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i++
		}
	}
	return string(b)
}

func RandStringBytesMaskImprSrc(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

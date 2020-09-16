package util

import (
	"alg/linked"
	"math/rand"
	"sort"
)

func InitListNodes(linkNum int) *linked.ListNode {
	headNode := &linked.ListNode{}
	var preNode *linked.ListNode
	for i := 1; i <= linkNum; i++ {
		curNote := linked.ListNode{Val: i, Next: nil}
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

func InitRangeLinkList(linkNum int, min int, max int) *linked.ListNode {
	headNode := &linked.ListNode{}
	if min >= max {
		return headNode
	}
	var preNode *linked.ListNode
	for i := 1; i <= linkNum; i++ {
		v := GenerateRangeNum(min, max)

		curNote := linked.ListNode{Val: v, Next: nil}
		if i == 1 {
			headNode = &curNote
		} else {
			preNode.Next = &curNote
		}
		preNode = &curNote
	}
	return headNode
}

func InitListRandNodes(linkNum int) *linked.ListNode {
	return InitRangeLinkList(linkNum, 10, 50)
}

func InitListSortRandNodes(linkNum int) *linked.ListNode {
	headNode := &linked.ListNode{}
	var preNode *linked.ListNode

	list := make([]int, linkNum)
	for i := 0; i < linkNum; i++ {
		v := GenerateRangeNum(10, 50)
		list[i] = v
	}
	sort.Ints(list)
	for i, v := range list {
		curNote := linked.ListNode{Val: v, Next: nil}
		if i == 0 {
			headNode = &curNote
		} else {
			preNode.Next = &curNote
		}
		preNode = &curNote
	}
	return headNode
}

func InitRandCycleLinkedList(linkNum int, pos int) *linked.ListNode {
	headNode := &linked.ListNode{}
	var preNode *linked.ListNode
	var posNode *linked.ListNode

	if linkNum < pos {
		return headNode
	}

	for i := 1; i <= linkNum; i++ {
		v := GenerateRangeNum(10, 50)

		curNote := linked.ListNode{Val: v, Next: nil}
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

func ChangeList2ListNode(list []int) *linked.ListNode {
	headNode := &linked.ListNode{}
	var preNode *linked.ListNode

	for i, v := range list {
		curNode := linked.ListNode{Val: v, Next: nil}
		if i == 0 {
			headNode = &curNode
		} else {
			preNode.Next = &curNode
		}
		preNode = &curNode
	}

	return headNode
}

func LinkedListPrint(head *linked.ListNode) []int {
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

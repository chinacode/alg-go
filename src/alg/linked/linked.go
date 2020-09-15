package linked

import (
	"math/rand"
	"sort"
)

type ListNode struct {
	Next *ListNode
	Val  int
}

func initListNodes(linkNum int) *ListNode {
	headNode := &ListNode{}
	var preNode *ListNode
	for i := 1; i <= linkNum; i++ {
		curNote := ListNode{Val: i, Next: nil}
		if i == 1 {
			headNode = &curNote
		} else {
			preNode.Next = &curNote
		}
		preNode = &curNote
	}
	return headNode
}

func generateRangeNum(min, max int) int {
	//rand.Seed(time.Now().Unix())
	randNum := rand.Intn(max-min) + min
	return randNum
}

func initListRandNodes(linkNum int) *ListNode {
	headNode := &ListNode{}
	var preNode *ListNode
	for i := 1; i <= linkNum; i++ {
		v := generateRangeNum(10, 50)

		curNote := ListNode{Val: v, Next: nil}
		if i == 1 {
			headNode = &curNote
		} else {
			preNode.Next = &curNote
		}
		preNode = &curNote
	}
	return headNode
}

func initListSortRandNodes(linkNum int) *ListNode {
	headNode := &ListNode{}
	var preNode *ListNode

	list := make([]int, linkNum)
	for i := 0; i < linkNum; i++ {
		v := generateRangeNum(10, 50)
		list[i] = v
	}
	sort.Ints(list)
	for i, v := range list {
		curNote := ListNode{Val: v, Next: nil}
		if i == 0 {
			headNode = &curNote
		} else {
			preNode.Next = &curNote
		}
		preNode = &curNote
	}
	return headNode
}

func initRandCycleLinkedList(linkNum int, pos int) *ListNode {
	headNode := &ListNode{}
	var preNode *ListNode
	var posNode *ListNode

	if linkNum < pos {
		return headNode
	}

	for i := 1; i <= linkNum; i++ {
		v := generateRangeNum(10, 50)

		curNote := ListNode{Val: v, Next: nil}
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

func changeList2ListNode(list []int) *ListNode {
	headNode := &ListNode{}
	var preNode *ListNode

	for i, v := range list {
		curNode := ListNode{Val: v, Next: nil}
		if i == 0 {
			headNode = &curNode
		} else {
			preNode.Next = &curNode
		}
		preNode = &curNode
	}

	return headNode
}

func linkedListPrint(head *ListNode) []int {
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

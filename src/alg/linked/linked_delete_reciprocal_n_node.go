package linked

import (
	"fmt"
	"util"
)

func deleteLinkedReciprocalNode2(head *ListNode, n int) *ListNode {
	result := &ListNode{}
	result.Next = head
	var pre *ListNode
	cur := result
	i := 1
	for head != nil {
		if i >= n {
			pre = cur
			cur = cur.Next
			//fmt.Println(i, n, pre.Val)
		}
		head = head.Next
		i++
	}
	pre.Next = pre.Next.Next
	return result.Next
}

/**
 * 刪除链表倒数的数值 ， 两个数值的遍历（一个先开始，一个后开始）；另外就是两次遍历先计算结构后遍历判断
 */
func deleteLinkedReciprocalNode(tmpListNodes *ListNode, reciprocalN int) *ListNode {
	index := 1
	var preNode *ListNode

	result := &ListNode{}
	result.Next = tmpListNodes
	curNode := result
	for nil != tmpListNodes {
		if index >= reciprocalN {
			//第二个循环后开始
			preNode = curNode
			curNode = curNode.Next
			fmt.Println(index, reciprocalN, curNode.Val, preNode.Val)
		}

		//第一个优先先开始
		tmpListNodes = tmpListNodes.Next
		index++
	}
	preNode.Next = preNode.Next.Next
	return result.Next
}

/**
 */
func TestDeleteLinkedReciprocalNode() {
	headNode := util.InitListNodes(10)

	var resultData *ListNode
	reciprocalN := 1
	loopCount := 1
	//loopCount = 200000
	util.Start("first", "")
	for i := 0; i < loopCount; i++ {
		resultData = deleteLinkedReciprocalNode(headNode, reciprocalN)
	}
	fmt.Println(util.LinkedListPrint(resultData))
	util.Cut("first", "")

	headNode = util.InitListNodes(10)
	util.Start("second", "")
	for i := 0; i < loopCount; i++ {
		resultData = deleteLinkedReciprocalNode2(headNode, reciprocalN)
	}
	fmt.Println(util.LinkedListPrint(resultData))
	util.Cut("second", "")

}

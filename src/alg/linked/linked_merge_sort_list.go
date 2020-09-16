package linked

import (
	"fmt"
	"util"
)

/**
 * 刪除链表倒数的数值 ， 两个数值的遍历（一个先开始，一个后开始）；另外就是两次遍历先计算结构后遍历判断
 */
func MergeSortList(node1 *ListNode, node2 *ListNode) *ListNode {
	if node1.Next == nil && node2.Next == nil {
		return nil
	}
	var tmpNode *ListNode
	var preNode *ListNode
	result := &ListNode{}

	for node1 != nil || node2 != nil {
		if nil == node1 {
			tmpNode = node2
			node2 = node2.Next
		} else if nil == node2 {
			tmpNode = node1
			node1 = node1.Next
		} else if node1.Val <= node2.Val {
			tmpNode = node1
			node1 = node1.Next
		} else {
			tmpNode = node2
			node2 = node2.Next
		}

		if nil == result.Next {
			result.Next = tmpNode
			preNode = tmpNode
		} else {
			preNode.Next = tmpNode
			preNode = tmpNode
		}
	}

	return result.Next
}

/**
先比较小的排序，后面看剩余部分直接连接即可
*/
func MergeSortList2(l1 *ListNode, l2 *ListNode) *ListNode {
	preHead := &ListNode{}
	result := preHead
	for l1 != nil && l2 != nil {
		if l1.Val < l2.Val {
			preHead.Next = l1
			l1 = l1.Next
		} else {
			preHead.Next = l2
			l2 = l2.Next
		}
		preHead = preHead.Next
	}
	if l1 != nil {
		preHead.Next = l1
	}
	if l2 != nil {
		preHead.Next = l2
	}
	return result.Next
}

/**
 */
func TestMergeSortList() {
	//headNode1 := changeList2ListNode([]int{18, 24, 33, 33, 40})
	//headNode2 := changeList2ListNode([]int{11, 28, 34, 36, 38})
	headNode1 := util.InitListSortRandNodes(10)
	headNode2 := util.InitListSortRandNodes(10)
	fmt.Println(util.LinkedListPrint(headNode1), util.LinkedListPrint(headNode2))

	var resultData *ListNode
	loopCount := 1
	//loopCount = 200000
	//util.Start("first", "")
	//for i := 0; i < loopCount; i++ {
	//	resultData = MergeSortList(headNode1, headNode2)
	//}
	//fmt.Println(linkedListPrint(resultData))
	//util.Cut("first", "")

	util.Start("second", "")
	for i := 0; i < loopCount; i++ {
		resultData = MergeSortList2(headNode1, headNode2)
	}
	fmt.Println(util.LinkedListPrint(resultData))
	util.Cut("second", "")

}

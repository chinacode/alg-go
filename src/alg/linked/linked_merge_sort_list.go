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
		if nil != node1 && (nil == node2.Next || node1.Val <= node2.Val) {
			tmpNode = node1
			node1 = node1.Next
		} else {
			tmpNode = node2
			node2 = node2.Next
		}
		if nil == result.Next {
			result.Next = preNode
			preNode = tmpNode
		} else {
			preNode.Next = tmpNode
			preNode = tmpNode
		}
	}

	return result.Next
}

/**
 */
func TestMergeSortList() {
	headNode1 := initListSortRandNodes(5)
	headNode2 := initListSortRandNodes(5)
	fmt.Println(linkedListPrint(headNode1), linkedListPrint(headNode2))

	var resultData *ListNode
	loopCount := 1
	//loopCount = 200000
	util.Start("first", "")
	for i := 0; i < loopCount; i++ {
		resultData = MergeSortList(headNode1, headNode2)
	}
	fmt.Println(linkedListPrint(resultData))
	util.Cut("first", "")

	//headNode = initListSortRandNodes(5)
	//util.Start("second", "")
	//for i := 0; i < loopCount; i++ {
	//	//resultData = MergeSortList2(headNode, reciprocalN)
	//}
	//fmt.Println(linkedListPrint(resultData))
	//util.Cut("second", "")

}

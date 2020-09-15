package linked

import (
	"fmt"
	"util"
)

/**
通过双方不同不长进行处理
*/
func linkedTwoNumberAdd(nodes *ListNode) bool {
	tmpNode1 := nodes
	tmpNode2 := nodes

	for nil != tmpNode1 && nil != tmpNode1.Next && nil != tmpNode2 {
		if tmpNode1 == tmpNode2 {
			return true
		}
		tmpNode1 = tmpNode1.Next.Next
		tmpNode2 = tmpNode2.Next
	}

	return false
}

/**
 */
func linkedTwoNumberAdd2(nodes *ListNode) bool {
	if nodes == nil {
		return false
	}
	fast := nodes.Next // 快指针，每次走两步
	for fast != nil && nodes != nil && fast.Next != nil {
		if fast == nodes { // 快慢指针相遇，表示有环
			return true
		}
		fast = fast.Next.Next
		nodes = nodes.Next // 慢指针，每次走一步
	}
	return false
}

/**
基础简单用法
*/
func linkedTwoNumberAdd3(head *ListNode) bool {
	m := make(map[*ListNode]int)
	for head != nil {
		if _, exist := m[head]; exist {
			return true
		}
		m[head] = 1
		head = head.Next
	}
	return false
}

/**
 */
func TestLinkedTwoNumberAdd() {
	headNodes := initRandCycleLinkedList(5, 2)

	var resultData bool
	loopCount := 1
	loopCount = 5000000
	util.Start("first", "")
	for i := 0; i < loopCount; i++ {
		resultData = linkedTwoNumberAdd(headNodes)
	}
	fmt.Println(resultData)
	util.Cut("first", "")

	util.Start("second", "")
	for i := 0; i < loopCount; i++ {
		resultData = linkedTwoNumberAdd2(headNodes)
	}
	fmt.Println(resultData)
	util.Cut("second", "")

	util.Start("third", "")
	for i := 0; i < loopCount; i++ {
		resultData = linkedTwoNumberAdd3(headNodes)
	}
	fmt.Println(resultData)
	util.Cut("third", "")

}

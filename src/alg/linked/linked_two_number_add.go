package linked

import (
	"config"
	"fmt"
	"util"
)

/**
通过双方不同不长进行处理
*/
func linkedTwoNumberAdd(nodes1 *config.ListNode, nodes2 *config.ListNode) *config.ListNode {
	var carry int
	var preNode *config.ListNode
	resultNodes := &config.ListNode{}
	for nil != nodes1 && nil != nodes2 {
		num := nodes1.Val + nodes2.Val + carry
		if num >= 10 {
			num = num - 10
			carry = 1
		} else {
			carry = 0
		}

		currentNode := &config.ListNode{Val: num, Next: nil}

		if nil == resultNodes.Next {
			resultNodes.Next = currentNode
		} else {
			preNode.Next = currentNode
		}

		preNode = currentNode

		nodes1 = nodes1.Next
		nodes2 = nodes2.Next
	}

	var extra *config.ListNode
	if nil != nodes1 {
		extra = nodes1
	}
	if nil != nodes2 {
		extra = nodes2
	}
	if carry == 0 {
		preNode.Next = extra
	} else {
		for nil != extra {
			num := extra.Val + carry
			if num >= 10 {
				num = num - 10
				carry = 1
			} else {
				preNode.Next = extra
				break
			}
			extra.Val = num

			preNode.Next = extra
			preNode = extra

			extra = extra.Next
		}
	}

	return resultNodes.Next
}

/**
 */
func linkedTwoNumberAdd2(l1 *config.ListNode, l2 *config.ListNode) *config.ListNode {
	list := &config.ListNode{Val: 0, Next: nil}
	//这里用一个result，只是为了后面返回节点方便，并无他用
	result := list
	tmp := 0
	for l1 != nil || l2 != nil || tmp != 0 {
		if l1 != nil {
			tmp += l1.Val
			l1 = l1.Next
		}
		if l2 != nil {
			tmp += l2.Val
			l2 = l2.Next
		}
		list.Next = &config.ListNode{nil, tmp % 10}
		tmp = tmp / 10
		list = list.Next
	}
	return result.Next
}

/**
 */
func TestLinkedTwoNumberAdd() {
	headNodes1 := util.InitRangeLinkList(5, 0, 9)
	headNodes2 := util.InitRangeLinkList(10, 0, 9)

	fmt.Println(util.LinkedListPrint(headNodes1))
	fmt.Println(util.LinkedListPrint(headNodes2))

	var resultData *config.ListNode
	loopCount := 1
	loopCount = 10000000
	util.Start("first", "")
	for i := 0; i < loopCount; i++ {
		resultData = linkedTwoNumberAdd(headNodes1, headNodes2)
	}
	fmt.Println(util.LinkedListPrint(resultData))
	util.Cut("first", "")

	util.Start("second", "")
	for i := 0; i < loopCount; i++ {
		resultData = linkedTwoNumberAdd2(headNodes1, headNodes2)
	}
	fmt.Println(util.LinkedListPrint(resultData))
	util.Cut("second", "")

	//util.Start("third", "")
	//for i := 0; i < loopCount; i++ {
	//	resultData = linkedTwoNumberAdd3(headNodes)
	//}
	//fmt.Println(resultData)
	//util.Cut("third", "")

}

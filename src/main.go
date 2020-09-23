package main

import (
	"alg/dynamic"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())
	//------------------array---------------------//
	//array.TwoArrayInterSectionTest()
	//array.TestRotateArray()
	//array.TestRemoveArray()
	//array.TestAddOneArray()
	//array.TestGetSumArray()
	//array.TestGetThreeSumArray()

	//------------------string---------------------//
	//string.TestLongestSamePrefix()
	//string.TestGetStringZChange()

	//------------------list---------------------//
	//linked.TestDeleteLinkedReciprocalNode()
	//linked.TestMergeSortList()
	//linked.TestCheckLinkedRing()
	//linked.TestLinkedTwoNumberAdd()

	//------------------life---------------------//
	//life.TestTheBestTimeForStocks()

	//------------------dynamic programming---------------------//
	//dynamic.TestClimbStairs()
	//dynamic.TestArrayMaxSubArray()
	//dynamic.TestMaxAscendingSubArray()
	//dynamic.TestMinTriangleSum()
	dynamic.TestMinPathSum()

}

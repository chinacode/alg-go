package main

import (
	"alg/linked"
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
	linked.TestMergeSortList()

	//------------------life---------------------//
	//life.TestTheBestTimeForStocks()

}

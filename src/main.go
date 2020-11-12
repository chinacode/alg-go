package main

import (
	"demo"
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
	//string.TestStringReverse()
	//string.TestFirstNoRepeatCharacter()

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
	//dynamic.TestMinPathSum()
	//dynamic.TestTheftMaxMoney()
	//dynamic.TestMinCoinChange()

	//const (
	//	letterIdxMask = 1<<6 - 1
	//)
	//rand_num := rand.Int63()
	//abc := rand_num | letterIdxMask
	//println(rand_num, abc)

	//demo
	//demo.ExampleClient()
	//demo.Dump()

	//list := []string{}
	//for i := 0; i < 2000; i++ {
	//	list = append(list, util.GetRandomName(5))
	//}
	//log.Println(strings.Join(list, ","))

	demo.Dump()
	//demo.Main()

	//demo.EmailTest()

}

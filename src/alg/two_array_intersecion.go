package alg

import (
	"fmt"
	"sort"
)

func arrayInterSection(arr1 []int32, arr2 []int32) []int32 {
	kv := map[int32]int32{}
	for _, v := range arr1 {
		kv[v] += 1
	}

	num := 0
	for _, v2 := range arr2 {
		if kv[v2] > 0 {
			arr2[num] = v2
			num++
		}
	}

	return arr2[:num]
}

func arrayInterSectionBySort(arr1 []int, arr2 []int) []int {
	j, k := 0, 0
	//fmt.Println(arr1)
	//fmt.Println(arr2)
	sort.Ints(arr1)
	sort.Ints(arr2)

	//var intersectionArr []int
	i := 0
	for j < len(arr1) && k < len(arr2) {
		v1 := arr1[j]
		v2 := arr2[k]
		//println(v1, v2)
		if v2 == v1 {
			//intersectionArr = append(intersectionArr, v1)
			arr1[i] = v1
			i++
			j++
			k++
		} else if v2 > v1 {
			j++
		} else {
			k++
		}
	}
	//return intersectionArr
	return arr1[:i]
}

func TwoArrayInterSectionTest() {
	arr1 := []int32{
		11, 22, 32, 44, 55,
	}
	arr2 := []int32{
		11, 21, 33, 44, 55,
	}
	interSection := arrayInterSection(arr1, arr2)
	fmt.Println(interSection)

	arr3 := []int{2, 6, 3, 4, 4, 10, 13}
	arr4 := []int{1, 2, 3, 9, 10, 4, 6, 13}

	interSectionSort := arrayInterSectionBySort(arr3, arr4)
	fmt.Println(interSectionSort)
}

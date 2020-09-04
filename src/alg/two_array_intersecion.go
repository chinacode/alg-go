package alg

import "fmt"

func ArrayInterSection(arr1 []int32, array []int32) []int32 {

	return []int32{}
}

func main() {
	arr1 := []int32{
		11, 22, 33, 44, 55,
	}
	arr2 := []int32{
		11, 22, 33, 44, 55,
	}
	interSection := ArrayInterSection(arr1, arr2)
	fmt.Print(interSection)
}

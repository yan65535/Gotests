package main

import "fmt"

// 泛型列表
// 这是一个泛型切片，类型约束为int | int32 | int64
type GenericsSlice[T int64 | float64] []T

// 泛型哈希表，键必须为可比较类型(comparable)
type GenericMap[K comparable, V int | string | byte] map[K]V

func main() {
	var aList GenericsSlice[int64]
	bList := GenericsSlice[float64]{7, 9, 8, 10}
	aList = append(aList, 1, 2, 3, 4, 5)
	fmt.Println(aList, bList)

	//gMap := GenericMap[string, string]{}
	//gMap["name"] = "yhw"

	gMap := make(GenericMap[string, string], 5)
	gMap["age"] = "4"
	fmt.Println(gMap)
}

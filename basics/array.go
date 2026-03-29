package basics

import "fmt"

var nums [5]int

func DemoArr() {
	nums2 := [3]int{}
	nums3 := [3]int{1, 3, 5}
	fmt.Println(nums)  // 0 0 0 0 0
	fmt.Println(nums2) // 0 0 0
	fmt.Println(nums3) // 1 3 5
}

// 72. Как проинициализировать массив при создании?
// Что если не все элементы будут переданы?
// Что если их будет больше?
// Какой синтаксис предусмотрен если длянну массива нужн вычеслить при компиляции т.е. что бы её не задавать явно?

func DemoArr2() {
	arr1 := [3]int{1, 2, 3}
	fmt.Println(arr1) // 1 2 3
	arr2 := [5]int{1, 2, 3, 4}
	fmt.Println(arr2) // 1 2 3 4 0
	// arr3 := [3]int{1, 2, 3, 4} // err
	arr4 := [...]int{1, 2, 3, 4, 5}
	fmt.Println(arr4) // 1 2 3 4 5
}

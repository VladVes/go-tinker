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

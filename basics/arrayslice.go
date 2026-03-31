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

// 75. Как массивы передаются и возвращаются из фунции?
// Сделать функцию modify что бы изменяла полученный массив как параметр

// массивы передаются по значению (в отличии от срезов) по этому нужно явно передвать указатель
func modify(arr *[3]int) {
	arr[0] = 100
}

func DemoModify() {
	arr1 := [3]int{1, 2, 3}
	// получиаем указатель на массив и передаём в как артумент в функцию
	modify(&arr1)
	fmt.Println(arr1)

}

// 76. Написать пример переворота строки - функцию reverse с помощь массива и классического цикла
func reverse(s string) string {
	runes := []rune(s)
	l := len(runes)
	result := make([]rune, l)
	for i := 0; i < l; i++ {
		result[i] = runes[l-1-i]
	}
	return string(result)
}

func DemoReverse() {
	s1 := "hello"
	s2 := "привет"
	fmt.Println(reverse(s1))
	fmt.Println(reverse(s2))
}

package basics

import (
	"fmt"
	"slices"
	"strings"

	"github.com/samber/lo"
)

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

// 80. С помощью какой специальной функции ещё можно￼ создать срез?
// Какие преимущества это даёт? Создать срез если изветсно что он будет расширяться но точно до 1000 элементов?
// С помощью какой функции можно посмотреть объём среза? Проэксперементировать с проверкой объёма среза на пустом срезе.

var sl1 []int
var sl2 = []int{}
var sl3 = []int{1, 2, 3}

func DemoMake() {
	fmt.Printf("%#v, len: %d cap: %d\n", sl1, len(sl1), cap(sl1))
	fmt.Printf("%#v, len: %d cap: %d\n", sl2, len(sl1), cap(sl2))
	fmt.Printf("%#v, len: %d cap: %d\n", sl3, len(sl3), cap(sl3))

	sl4 := make([]int, 0, 1000)
	fmt.Printf("%#v, len: %d cap: %d\n", sl4, len(sl4), cap(sl4))
}

// 81. Go основы. Курс Hexlet.

// Как добавить элемент в срез?
// Что при этом происходит?
// Что если превышается вместимость?
// Как добавить несколько элементов?
// Создать срез  литерально []int{1,2,3}, добавить один, затем несколько элементов.
func DemoSlAppendMake() {
	slA := []int{1, 2, 3}
	fmt.Println(slA)
	slA = append(slA, 4)
	fmt.Println(slA)
	slA = append(slA, 6, 7, 8)
	fmt.Println(slA)
	// Создать срез с заданым значением cap с помощью спец функции и попробовать добавить элементов больше чем cap
	slB := make([]int, 0, 5)
	fmt.Println(slB)
	fmt.Println(cap(slB))
	slB = append(slB, 1, 2, 3, 4, 5, 6)
	fmt.Println(slB)
	fmt.Println(cap(slB))
	// Как объединить два среза?
	slAB := append(slA, slB...)
	fmt.Println(slAB)
	fmt.Println(cap(slAB))
}

// 82. Как получить доступ к элементу среза?
// Что будет если обратиться по несуществующиму индексу?
// Что нужн сделать что бы предотвартить ошибку?
func DemoSlIndex() {
	sl := []int{1, 2, 3}
	index := 3
	// fmt.Println(sl[index]) // panic: runtime error: index out of range [3] with length 3
	if index > len(sl)-1 {
		fmt.Println("index out of range")
	}
}

// 86. Написать срез состоящий из структур User (Name, Email)
// и функцию писка пользователя по Email. Фунция должна вернуть указатель на элемнт среза.

func findUser(users []User) *User {
	for i, u := range users {
		if strings.HasPrefix(u.Email, "ivan") {
			return &users[i]
		}
	}
	return nil
}

func FindUserDemo() {
	users := []User{
		{Name: "John", Email: "john@mail.com"},
		{Name: "Alice", Email: "alice@mail.com"},
		{Name: "Ivan", Email: "ivan@mail.com"},
	}

	fmt.Println(findUser(users))

}

// 88. Как и с помощью чего можно сравнить содержимое срезов?
// Создать два слайса и сравнить их содежимое.

func DemoSliceEqual() {
	fmt.Println("----slices.Equal----")
	sl1 := []int{1, 2, 3}
	sl2 := []int{1, 2, 3}
	sl3 := []int{3, 3, 3}

	fmt.Println(slices.Equal(sl1, sl2))
	fmt.Println(slices.Equal(sl1, sl3))
}

// 89. Как сравнить срезы на равенство ссылок?

func DemoSlicePtrEqual() {
	sl1 := []int{1, 2, 3}
	sl2 := sl1

	fmt.Println(&sl1[0] == &sl2[0]) // способ работает при условии что оба среза не пустые
}

// 92. Как можно удалить элемент из среза?
// Создать срез из элементов "a", "b", "c", "d", "e".
// Удалить первый элемент, удалить последний, удалить диапазона

func DemoSliceDeleteElem() {
	s := []string{"a", "b", "c", "d", "e"}
	fmt.Println(s[1:])
	fmt.Println(s[:len(s)-1])
	fmt.Println(append(s[:1], s[3:]...))
}

// 93. Как проверить срез на вхождение элемента?
// Создать срез с элементами 1, 2, 3 и проверить входит ли в него к прмеру 2 или 5
func DemoSliceContains() {
	s := []int{1, 2, 3}
	fmt.Println(slices.Contains(s, 5))
}

// 94. Как можно удалить дубликаты из среза?
// Использую два способа из среза {2, 1, 2, 3, 1}
// удалить дубликаты. В чем разница этих способов?
func DemoSliceUniq() {
	s1 := []int{1, 1, 2, 2, 3, 3}
	s3 := []int{0, 22, 1, 1, 22, 2, 3, 3}
	s2 := []int{22, 1, 2, 1, 3, 33, 3}
	fmt.Println(slices.Compact(s1)) // только для упорядоченных
	fmt.Println(slices.Compact(s3)) // только для упорядоченных [0 22 1 22 2 3]
	fmt.Println(lo.Uniq(s2))
}

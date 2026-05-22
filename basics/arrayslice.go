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

// 95. Как объяединить (union) два среза?
// Создать 2 среза {1, 2, 3} и {3, 4, 5} и объединить с помощью спец. функции
// из определённого пакета.
// Попробовать сделать тоже но уже с стандартной функции

func DemoSlicesUnion() {
	s1 := []int{1, 2, 3}
	s2 := []int{4, 5, 6}
	s3 := lo.Union(s1, s2)
	fmt.Println(s3)
	s4 := append(s1, s2...)
	fmt.Println(s1)
	fmt.Println(s4)
}

// Как найти пересечение (intersection) двух срезов? Создать 2 среза {1, 2, 3} и {2, 3}
func DemoSlicesIntersec() {
	s1 := []int{1, 2, 3}
	s2 := []int{2, 3}
	fmt.Println(lo.Intersect(s1, s2))

}

// Как найти разность (difference) двух срезов?
// Создать 2 среза: A {1, 2, 3} и B {2, 4} и найти элементы
// которые есть только в срезе А с помощью спец. функции
func DemoSlicesDiff() {
	s1 := []int{1, 2, 3}
	s2 := []int{2, 4}
	a, b := lo.Difference(s1, s2)
	fmt.Println(a, b)
	// Что такое и как найти симметрическую разность?
	// Создать 2 среза: A {1, 2, 3} и B {2, 4}
	// и найти элементы которые есть только в срезе А или только в срезе B?
	fmt.Println(append(a, b...))

}

// 99. Создать срез A с данными {1, 2, 3}, создать переменную B и
// проинициализировать сразуже присвоив ей значение A.
// Что произойдёт при попытке просто присваивания среза переменной B?
// Почему (как передаются срезы)?
// Попробовать измеить элемент по индексу 0 в []B. Что произойдёт?

func DemoSlicesAssign() {
	a := []int{1, 2, 3}
	b := a
	b[0] = 100
	fmt.Println(a) // {100, 2, 3}
}

// 100. С помощью какой функции и из какого пакета можно создать новый срез
// с тем же содержимым не изменяя оригина?
// Создать срез original c данными {1,2,3} и скопировать его с помощью этой функции
// в переменную clone. Убедиться что изменение элеметна в clone не влияет на original

func DemoSlicesClone() {
	original := []int{1, 2, 3}
	clone := slices.Clone(original)
	clone[0] = 100
	fmt.Println(clone)
	fmt.Println(original)
}

// 101. Как можно произвести частичное копирование элементов одного среза в другой?
// Какую специальную функцию можно для этого использовать?
// Создать срез src с данными {1, 2, 3}, создать dst с заданной длинной
// (с помощью спец. функции) в 2 элемента, скопировать src в dst,
// вывести содержимое dst

func DemoSlicesCopy() {
	src := []int{1, 2, 3}
	dst := make([]int, 2)
	copy(dst, src)
	fmt.Println(dst)
}

// 102. Каким способом можно скопировать срез не используюя специально
// предназначенные для этого функции типа slices.Clone и copy?
// Создать src с элементами {1, 2, 3} и скопировать его этим способом
// в переменную dst

func DemoSlicesCopyWithoutCopyAndClone() {
	src := []int{1, 2, 3}
	dst := append([]int(nil), src...)
	fmt.Println(dst)
}

// 103. Написать пример сортировки пузырьком для []int.
func BubbleSort(nums []int) {
	n := len(nums)
	for i := 0; i < n-1; i++ {
		swapped := false
		for j := 0; j < n-1-i; j++ {
			if nums[j] > nums[j+1] {
				nums[j], nums[j+1] = nums[j+1], nums[j]
				swapped = true
			}
		}
		if !swapped {
			break
		}
	}
}

func DemoBubbleSort() {
	values := []int{7, 3, 10, 1, 4, 2}
	fmt.Println("before:", values)
	BubbleSort(values)
	fmt.Println("after: ", values)
}

// 104. Создать функцию addElement которая получате срез и добавляет в него
// с помощью append новый элемент. Как срезы передаются в функцию?
// Что нужно учитывать при работе с append в функции изменяющей полученный срез?
// Сделать пример в котором демонстрируется особенность append - доработать,
// разобраться как сделать так что бы append добавляла элемент но так
// что бы не создавался новый массив т.е. как сделать что бы функция которая
// в данном примере исп. append изменяла массив исходного среза?

func modifySlice[T any](sl *[]T, elem T) {
	_ = append(*sl, elem)
	fmt.Println("from modify", sl)
}

func DemoSlicesModify() {
	a := make([]int, 0, 10)
	modifySlice(&a, 100)
	fmt.Println(a)
}

// 105. Какая проблема может возникнуть если из исходного среза с 1000 элементами
// взять подсрез с элементами от 0 до 10 индекса? Как решить такую проблему?
// Создать int срез А из 10000 элементов
// Cоздать срез B который ссылкается на окно элементов из среза А от
// 0 индекса до 10 способом воспроизводящим проблему.
// Создать срез С из тогоже набора элементов что и срез B (т.е. из окна с 0 по 10 индекс)
// но так что бы проблема решилась.

func DemoSlicePart() {
	A := make([]int, 1000)
	B := A[:10]
	fmt.Println(B)
	C := slices.Clone(A[:10])
	fmt.Println(C)
}

// 107. Go основы. Курс Hexlet. Как отсортировать срез упорядоченного типа (int, string)
// по возрастанию быстрее всего?
// Создать срез int с элементами {5, 2, 9, 1, 3} и отсортировать двумя способами
// один из них должен изменить исходный срез а другой создать новый.
// Какой тут есть ньюан с использованием функций стандартной библиотеки?
// Как проверить отсортирован ли срез? Проверить до сортировки и после.

func DemoSliceSort() {
	x := []int{5, 2, 9, 1, 3}
	slices.Sort(x)
	fmt.Println(x)
	fmt.Println(slices.IsSorted(x))
}

// 108. Как можно произвести сортировку среза с кастомной логикой?
// Создать срез с данными {"banana", "apple", "cherry"} и отсортировать по длинне слов

func DemoSliceCustomSort() {
	words := []string{"banana", "apple", "cherry"}
	fmt.Println(words)
	slices.SortFunc(words, func(a, b string) int {
		return len(a) - len(b)
	})
	fmt.Println(words)
}

// 109. Как отсортировать срез в обратном порядке?
// Создать срез {5, 2, 9, 1, 3} и отсортировать так,
// что бы в резульате получился срез [9 5 3 2 1]

func DemoSliceSortReverse() {
	s := []int{5, 2, 9, 1, 3}
	fmt.Println(s)
	slices.Sort(s)
	slices.Reverse(s)
	fmt.Println(s)
}

// 110. C помощью какой функции и из какого пакета можно определять
// минимальный и максимальный по значению элемет среза?

func DemoSliceMinMax() {
	s := []int{5, 2, 9, 1, 3}
	fmt.Println(s)
	max := slices.Max(s)
	min := slices.Min(s)
	fmt.Println(max)
	fmt.Println(min)
	fmt.Println(max)
	fmt.Println(min)
}

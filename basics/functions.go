package basics

import (
	"cmp"
	"errors"
	"fmt"
)

// 9. Какой особенностью обладают функции в Go в плане возвращения значений?
// Написать функцию divide принимающую два значения x y и возвращающая результат деления x на y.
// Как должна обрабатываться ситуация когда в кач-ве y передаётся 0? Как можно в Go создать объект ошибки?

func DivideByZero(x, y int) (float64, error) {
	if y == 0 {
		return 0, errors.New("divide by zero!")
	}
	return float64(x / y), nil
}

// 66. Как написать анонимную функци складывающую два числа?
// Для чего нужны? Как её вызвать? Как называется такая конструкция (когда идёт описание функции и сразу вызов)?

var Result = func(a, b int) int {
	return a + b
}(15, 5)

// 67. Как работают замыкания с анонимнными функциями?
// Пример когда внешняя переменная содержит часть строки используемой
// в формировании результата работы анонимной функции.

var strOuter = "Hi!"
var GreetStr = func(str string) string {
	return fmt.Sprintf("%s What is youre %s?", strOuter, str)
}("name")

// 163. Сделать пример дженерик-функции Max для решения задачи поиска максимума в срезе,
// при том, что срез может быть int или float64 или string типов
func Max[T cmp.Ordered](values []T) T {
	max := values[0]
	for _, v := range values {
		if v > max {
			max = v
		}
	}
	return max
}
func DemoMax() {
	iAr := []int{23, 44, 2, 55, 2}
	fAr := []float64{23.323, 414.55, 2.123, 55.333, 2.54}
	sAr := []string{"abc", "dfg", "zzzdddd"}
	fmt.Println(Max(iAr))
	fmt.Println(Max(fAr))
	fmt.Println(Max(sAr))
}

// 164. Сделать пример дженерик-функции которая можетс складывать
// два числа как int так и float64 типов. Сделать свой интерфейс - ограничение

type Number interface {
	int | float64
}

func Adder[T Number](a, b T) T {
	return a + b
}

func DemoAdder() {
	fmt.Println(Adder(44.22, 76.23))
	fmt.Println(Adder(22, 23))
}

// 68. Что такое функция высшего порядка? Написать функцию MakeMultiplier,
// которая принимает параметр factor int и возвращает другую функцию которая
// так же принимает x int и используя механизм замыканий обращается к переменной
// factor умножая на неё свой параметр x и возвращая значение?

// first class citizen
func MakeMultiplier(factor int) func(x int) int {
	return func(x int) int {
		return factor * x
	}
}

func TestMultiplier() {
	double := MakeMultiplier(2)
	triple := MakeMultiplier(3)
	fmt.Println(double(4)) // 8
	fmt.Println(triple(5)) // 15
}

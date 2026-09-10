package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// echo
// 177. Написать CLI выводящее все переданные аргументы
// на подобии утилиты Unix echo
// и цикл на стандартную функцию
// Измените программу так, чтобы она выводила также
// os.Args[0], имя выполняемой команды.
func echo1() {
	var s, sep string
	for i := 1; i < len(os.Args); i++ {
		s += sep + os.Args[i]
		sep = " "
	}
	fmt.Println(s)
}

// Измените программу так, чтобы она выводила индекс и
// значение каждого аргумента по одному аргументу в строке.
func echo2() {
	for i, ch := range os.Args[1:] {
		fmt.Printf("%d: %s\n", i, ch)
	}
}

// Переделать так что бы не создавалось много строк т.е. заменить += и цикл на стандартую функцию
func echo3() {
	fmt.Println(strings.Join(os.Args[1:], " "))
}

func echo4() {
	fmt.Println(os.Args)
}

// Поэкспериментируйте с измерением разницы времени выполне­
// ния потенциально неэффективных версий и версии с применением strings.Join
func echo3withTimeCheck() {
	start1 := time.Now()
	var s, sep string
	for i := 1; i < len(os.Args); i++ {
		s += sep + os.Args[i]
		sep = " "
	}
	fmt.Println(s)
	fmt.Println(time.Since(start1).Seconds())
	fmt.Println(time.Since(start1).Microseconds())

	start2 := time.Now()
	fmt.Println(strings.Join(os.Args[1:], " "))
	fmt.Println(time.Since(start2).Seconds())
	fmt.Println(time.Since(start2).Microseconds())
}

func main() {
	echo1()
	echo2()
	echo3()
	echo4()
	echo3withTimeCheck()
}

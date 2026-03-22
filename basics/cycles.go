package basics

import "fmt"

func CyclesDemo() {
	// 32.Как можно использовать for как while т.е. цикл только с условием?
	y := 10
	x := 0
	for x < y {
		fmt.Printf("x = %d\n", x)
		x++
	}
}

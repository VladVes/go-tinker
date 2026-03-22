package basics

import "fmt"

func CondDemo() {
	// 24. Как можно использовать переменнe в условном операторе if, где они видны?
	str := "some demo str"
	if x := len(str); x < 10 {
		fmt.Println("x < 10")
	} else if x > 10 {
		fmt.Println("x > 10")
	} else {
		fmt.Println("x == 10")
	}
}

package main

import "fmt"

var i, j, k int                 // все будут с дефолтным для int значением 0
var b, f, s = true, 2.3, "four" // можно и так проинициализировать

func main() {
	li, lj, lk := 0, 1, 3 // в краткой форме тоже работает такое определение в строку
	lb, lf, ls := true, 2.4, "four"
	fmt.Println(lb, lf, ls)
	fmt.Println(li, lj, lk)
	fmt.Println(i, j, k)
	fmt.Println(f, b, s)
}

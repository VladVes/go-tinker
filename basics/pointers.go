package basics

import "fmt"

// 45. Как в функцию передать переменную по адресу?
// Сделать простой пример где func change(number int) {number = 10}
// должна изменить значение исходной переменной которая ей передаётся.
// Как называется оператор ? Объснить значения операторов *  и  &
func incNum(n *int) {
	*n++
}

func DemoPtr() {
	n := 0
	fmt.Printf("DemoPtr - num == %d", n)
	incNum(&n)
	fmt.Printf("DemoPtr - change num: %d", n)
	incNum(&n)
	fmt.Printf("DemoPtr - change num: %d", n)
}

// 49. Что гарантирует Go при возврате указателя в плане безопасности?
func CreateEntity() *User {
	u := User{Name: "Charlie", Age: 40}
	return &u // безопасно: Go поместит структуру в heap
} // компилятор гарантирует, что указатель не станет висячим (dangling pointer).


// 51. Какое специальное значение используется для обозначения отсутсвтия данных?
// Для какиех типов оно может быть исползовано?
var x *int
var z = 10

func DemoNil() {
	x = &z
	fmt.Println(x)
	fmt.Println(*x)
	x = nil
	fmt.Println(x)
	// fmt.Println(*x) // паника
}

// 52. var ptr *int Чему будет равен указатель?
// Что будет если попытаться разыменовать указатель (не был проинициализирован) т.е. получить по нему значение?
// Как решить эту проблему?
// Что если поле структруры будет иметь такое знаение?
var ptr *int

func DemoNilPtr() {
	fmt.Println(ptr)
	// fmt.Println(*ptr) // паника
	if ptr != nil {
		fmt.Print(*ptr)
	}
}

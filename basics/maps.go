package basics

import "fmt"

// 111. Какой есть встроенный тип данных для хранения пар ключ - значение.
// Как создаётся (объявляется)?
// Чему будет равна переменная этого типа если не определить значения?
// Какого типа может быть ключ?
// Создать с ключом типа string и значением типа int.
// Что нужно сделать что бы можно было использвать?
// Внести несколько пар Alice - 25, Bob - 30, Joe - 40 Вывести на экран.
// Как обявить такую переменную и сразу проинициализировать значениями?

func DemoMap() {
	var m1 map[string]int
	fmt.Println(m1)        // map[]
	fmt.Println(m1 == nil) //true

	ages := make(map[string]int)
	fmt.Println(ages)        // map[]
	fmt.Println(ages == nil) // false

	ages["Alice"] = 25
	ages["Bob"] = 30
	fmt.Println(ages) // map[Alice:25 Bob:30]

	scores := map[string]int{
		"Tom":  90,
		"Jack": 80,
	}
	fmt.Println(scores) // map[Jack:80 Tom:90]

}

// 111.1 Как обратиться к элементу map?
// Создать map c данными Alice - 25, Bob - 30, John - 40,
// Получить значение по любому ключу
// Что будет если такого элемнта нет?
// Как проверить существует ли ключ в map?
// попробовать обратиться к ключу Zella, попробовать проверить существует ли ключ Joe

func DemoMapElement() {
	ages := map[string]int{
		"Alice": 25,
		"Bob":   30,
		"Vova":  43,
	}
	fmt.Println(ages["Alice"])
	fmt.Println(ages["Ivan"]) // 0 !
	age, exists := ages["Ivan"]
	fmt.Println(age, exists) // 0 false

}

// 112. Как удалить элемент из map?
// Создать map c данными
// "Alice": 25,
// "Bob": 30,
// "John": 40,
// Раелизовать логику при которой проиходит проверка сеществования ключа
// и если он существует, то удалять элемент,
// если нет то писать что такого элемента в мап нет.
// Удалить по ключу Alice. Что будет если попробовать удалить по ключу Joe без проверки?

func DemoMapRemoveElem() {
	ages := map[string]int{
		"Alice": 25,
		"Bob":   30,
		"John":  40,
	}
	fmt.Println(ages)
	if age, exists := ages["Alice"]; exists {
		fmt.Printf("Element \"Alice\" has been found. Value = %d.\nDeleting element \"Alice\"...\n", age)
		delete(ages, "Alice")
		fmt.Println(ages)
	}
	delete(ages, "Ivan") // no error
}

// 113. Как узнать число пар в map?
// Создать map c данными Alice - 25, Bob - 30, John - 40, узнать длинну

func DemoMapLen() {
	ages := map[string]int{
		"Alice": 25,
		"Bob":   30,
		"John":  40,
	}
	fmt.Println(ages)
	fmt.Println(len(ages))
}

// 114. Go основы. Курс Hexlet
// Создать map c парами Alice - 25, Bob - 30, John - 40
// сделать обход по всем элементам с выводом и ключей и значений.
// Что важно учитывать при обходе map (чем отличается от обхода среза)?

func DemoMapLoop() {
	ages := map[string]int{
		"Alice": 25,
		"Bob":   30,
		"John":  40,
	}
	for k, v := range ages {
		fmt.Println(k, v)
	}
	// 115. Создать map c парами Alice - 25, Bob - 30, John - 40
	// сделать обход при которм должен выбираться
	// первая попавшаяся пара у которй значение == 30

	for name, age := range ages {
		if age > 25 {
			fmt.Println(name, age)
			break
		}
	}

}

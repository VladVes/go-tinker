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

// 117. Как можно хранить структуры в map?
// Создать структуру User с полями Name, Email;
// создать map где ключи типа int (1, 2, 3...), а значения типа User:
// {Name: "Alice", Email: "alice@example.com"}, {Name: "Bob", Email: "bob@example.com"},.
// Создать переменную user и проинициализировать структурой лежащей в map по ключу 1.
// Что на самом деле попадёт в переменную? Что будет в map по ключю 1 если изменить поле Name в переменной user?
// Как можно обойти добиться обратного эффекта (два способа)?

type Usr struct {
	Name  string
	Emali string
}

func DemoMapStruct() {
	fmt.Println("---- map A ----")
	var usersA = map[int]User{
		1: {Name: "Alice", Email: "alice@example.com"},
		2: {Name: "Bob", Email: "bob@example.com"},
	}
	fmt.Println(usersA)
	user1 := usersA[1]
	user1.Name = "Alize"
	fmt.Println(usersA)
	fmt.Println(user1)

	usersA[1] = user1
	fmt.Println(usersA)

	fmt.Println("---- map B ----")
	u1 := User{Name: "Alice", Email: "alice@example.com"}
	u2 := User{Name: "Bob", Email: "bob@example.com"}
	var usersB = map[int]*User{
		1: &u1,
		2: &u2,
		3: {Name: "Jane", Email: "Jane@example.com"},
	}
	user2 := usersB[2]
	fmt.Println(usersB)
	fmt.Println(user2)

	user2.Name = "John"
	fmt.Println(usersB)
	fmt.Println(usersB[2])
	fmt.Println(user2)

}

// 119. Созадать map settings с ключами string (имена пользвателей)
// и значением map в котором ключи и значения string, к примеру,
// это настройки интерфейса типа "theme": "dark", "lang": "en".
// Что будет если попытаться получить настройки для несуществующего
// пользователя?
// Добавить новую пару (пользователь, настройки) в map settings.
// Изменить настройки существующего пользователя.
// Удалить из map одну из настроек какого-нибудь пользователя.
// Удалить пользователя со всем его настройками из map.

func DemoMapMap() {
	users := map[string]map[string]string{
		"Ivan":  {"theme": "darck", "lang": "ru"},
		"John":  {"theme": "light", "lang": "en"},
		"Jango": {"theme": "custom", "lang": "es"},
	}
	fmt.Println(users)
	fmt.Println(users["Ivan"])
	fmt.Println(users["Alice"])

	users["Alice"] = map[string]string{
		"theme": "purple",
		"lang":  "de",
	}
	fmt.Println(users)

	users["Jango"]["theme"] = "blue"
	users["Jango"]["lang"] = "fr"
	fmt.Println(users)

	delete(users["John"], "lang")
	delete(users, "Alice")
	fmt.Println(users)

}

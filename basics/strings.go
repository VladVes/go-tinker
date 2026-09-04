package basics

import (
	"encoding/json"
	"fmt"
	"strings"
)

// 15. Каким способом удобней записать многострочную строку или строку со спец символами в переменную?
var query = `
	SELECT *
    FROM users
    WHERE active = true
`

func LStr() {
	fmt.Println(query)
}

// 20. Как можно обратиться к символу в строке
func PrintCharByIndexStrDemo() {
	var str = "Привет"
	ch := str[0]               //
	fmt.Println(ch)            // выведет 10тичное предстваление первого байта строки
	runes := []rune(str)       // нужно создать срез рун
	fmt.Println(runes[0])      // выведет код первого символа в utf8
	fmt.Printf("%c", runes[0]) // выведет символ
}

// 40. В чем может быть проблема если использовать много оперций конкатенации
// (к примеру при сборке HTML, SQL или при генерации отчёта)?
// Какой есть механизм для эффективного работа со строками (в чем его эффективность)?
// Составить строку из 3х частей ("Hello", " ", "World" )
// не используя конкатенацию + и без шаблонизации с помощью fmt и записать её в переменную result?
var (
	sp1 = "Hello"
	sp2 = "Go!"
)

func StrBldr() {
	var sb strings.Builder
	sb.WriteString(sp1)
	sb.WriteRune(' ')
	sb.WriteString(sp2)
	sb.WriteByte('!')

	fmt.Println(sb.String())
}

// 43. Дана строка "abcdefg", нужно реализовать логику
// (можно в функции main) которая вернёт строку "a,b,c,d,e,f,g"
// но без конкатенации (+) и шаблонизации (fmt.Sprintf)

func DemoStringTransform() {
	s := "abcdefg"
	var sb strings.Builder
	for i, ch := range s {
		sb.WriteRune(ch)
		if i < (len(s) - 1) {
			sb.WriteRune(',')
		}
	}
	fmt.Println(sb.String())
}

// 76. Написать пример переворота строки -
// функцию reverse с помощью среза и цикла

func StrReverse() {
	s := "Reverse me! Ха-Ха"
	runeS := []rune(s)
	fmt.Println(len(s))
	fmt.Println(len(runeS))
	result := make([]rune, 0, len(s))
	for i := range runeS {
		result = append(result, runeS[len(runeS)-1-i])
	}

	fmt.Println(string(result))
}

// 149. C помощью каких функции можно сериализовать и десериализовать данные в Go?
// Создать структуру Person c типичными полями и доработать так что бы:
// - было поле с паролем но что бы оно не сериализовалось
// - поля ID, Age в json было в нижнем регистре,
// - поле Name в json записалoсь как Username
// - поле age не записывалось если оно пустое т.е. имеет нулевое значение
// Что возвращает сериализующая функция из стандарнтной либы и как это корретно вывести в консоль в читабельном виде?
// C помощью какой функции из стандартной бибилотеки можно сделать более красивый построчтный вывод полей json?

type Person struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"Username"`
	Age      int    `json:",omitempty"`
	password string
}

func SerializeDemo() {
	person := Person{
		ID:       234234,
		Email:    "person@mail.com",
		Name:     "Ivan",
		Age:      0,
		password: "superpass",
	}

	personSerialized, err := json.Marshal(person)
	if err != nil {
		fmt.Printf("Json serialize error: %v", err)
	}

	personSerialized2, err2 := json.MarshalIndent(person, " ", " ")
	if err2 != nil {
		fmt.Printf("Json serialize error: %v", err2)
	}

	fmt.Println(string(personSerialized))
	fmt.Println(string(personSerialized2))
}

package basics

import (
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

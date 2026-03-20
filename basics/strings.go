package basics

import "fmt"

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

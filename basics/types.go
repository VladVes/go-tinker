package basics

import (
	"fmt"
	"strconv"
)

// 17. Как можно преобразовать число к строке что нужно учитывать?
// Сделать пример где создать переменные типа int, int64, float64
// вывести все переменные как строки, int64 вывести в разных системах счисления
// вывести float с заданым числом знаков
// Создать переменные типа строки "234243", "33.4452" и преобразовть к int, int64, float64
// Что будет являтся результатом работы функций прербразовывающий строку к числу?

func DemoNumToStrToNum() {
	var i int = 111
	var i64 int64 = 111243234
	var f float64 = 444.3323

	fmt.Println(strconv.Itoa(i))
	fmt.Println(strconv.FormatInt(i64, 10))
	fmt.Println(strconv.FormatInt(i64, 2))
	fmt.Println(strconv.FormatFloat(f, 'f', 2, 64))
	fmt.Println(strconv.FormatFloat(f, 'f', 5, 64))

	var si = "111"
	var si64 = "234234234"
	var sf = "3324.7567"

	ri, _ := strconv.Atoi(si)
	fmt.Println(ri)
	ri64, _ := strconv.ParseInt(si64, 10, 64)
	fmt.Println(ri64)
	rf, _ := strconv.ParseFloat(sf, 64)
	fmt.Println(rf)
}

// 57. Как преобразовать bool к числу или к строке?
var b bool

func Btos() {
	bstr := fmt.Sprintf("%v", b)
	fmt.Printf("bool to string: %s", bstr)
}

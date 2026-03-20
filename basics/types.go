package basics

import "fmt"

// 57. Как преобразовать bool к числу или к строке?
var b bool

func Btos() {
	bstr := fmt.Sprintf("%v", b)
	fmt.Printf("bool to string: %s", bstr)
}

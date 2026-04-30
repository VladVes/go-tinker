package main

import (
	"fmt"
	"strings"

	"github.com/VladVes/go-tinker/v2/basics"
	grtg "github.com/VladVes/go-tinker/v2/greeting"
	clr "github.com/fatih/color"
	// "github.com/VladVes/go-tinker/v2/websrv"
	// "github.com/VladVes/go-tinker/v2/webfibersrv"
)

func PrintTopic(div, name string) {
	const maxDivLen = 100
	var b strings.Builder

	nL := len([]rune(name))
	divPartL := (maxDivLen - nL) / 2
	divPart := strings.Repeat(div, divPartL)
	b.WriteString(divPart)
	b.WriteString(" ")
	b.WriteString(name)
	b.WriteString(" ")
	b.WriteString(divPart)
	if divPartL*2+nL < maxDivLen {
		b.WriteString(div)
	}
	fmt.Println(b.String())
}

func main() {
	fmt.Println(grtg.Hello())
	clr.Magenta(grtg.Hello())
	clr.Green(grtg.Hello())
	// // nums := []int{1, 2, 3, 4, 5, 33, 7, 56, 99, 13, 12}
	// // for _, v := range ev.Even(nums) {
	// // 	clr.Red(strconv.Itoa(v))
	// }
	// ------------------------------- BASICS --------------------------
	// strSet := basics.NewSet[string]()
	// fmt.Println(strSet)
	// strSet.Add("first")
	// strSet.Add("second")
	// strSet.Add("third")
	// strSet.Add("fourth")
	// fmt.Println(strSet)
	// fmt.Println(strSet.Has("second"))
	// strSet.Remove("second")
	// fmt.Println(strSet.Has("second"))

	// basics.DemoMake()

	basics.DemoSlAppendMake()
	basics.DemoSlIndex()
	basics.DemoSliceEqual()
	basics.DemoSlicePtrEqual()
	PrintTopic("*", "Slices Union")
	basics.DemoSlicesUnion()
	PrintTopic("*", "Slices Intersection")
	basics.DemoSlicesIntersec()
	PrintTopic("*", "Slices difference")
	basics.DemoSlicesDiff()
	PrintTopic("*", "Modify Slice")
	basics.DemoSlicesModify()
}

package main

import (
	"fmt"
	"strings"

	"github.com/VladVes/go-tinker/v2/basics"
	grtg "github.com/VladVes/go-tinker/v2/internal/greeting"
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
	PrintTopic("*", "Bubble sort")
	basics.DemoBubbleSort()
	PrintTopic("*", "Modify Slice")
	basics.DemoSlicesModify()
	PrintTopic("=", "Map")
	basics.DemoMap()
	PrintTopic("=", "Map Element")
	basics.DemoMapElement()
	PrintTopic("=", "Map Delete Element")
	basics.DemoMapRemoveElem()
	PrintTopic("=", "Map len")
	basics.DemoMapLen()
	PrintTopic("=", "Map loop")
	basics.DemoMapLoop()

	PrintTopic("=", "Map with struct")
	basics.DemoMapStruct()

	PrintTopic("=", "Map of maps")
	basics.DemoMapMap()

	PrintTopic("=", "Sort map by keys")
	basics.DemoMapSort()

	PrintTopic("=", "Map of maps update")
	basics.DemoMapMapMutate()

	PrintTopic("=", "Map for counters")
	basics.DemoMapForCounters()

	PrintTopic("=", "Check Map for keys")
	basics.DemoCheckMapKey()

	PrintTopic("=", "Extract keys and values from map")
	basics.DemoExtractKeysAndValuesFromMap()
}

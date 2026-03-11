package main

import (
	"fmt"

	grtg "github.com/VladVes/go-tinker/greeting"
	clr "github.com/fatih/color"
)

func main() {
	fmt.Println(grtg.Hello())
	clr.Magenta(grtg.Hello())
}

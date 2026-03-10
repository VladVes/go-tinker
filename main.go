package main

import (
	"fmt"
	grtg "go-tinker/greeting"
	"github.com/fatih/color"
)

func main() {
	fmt.Println(grtg.Hello())
	color.Magenta(grtg.Hello())
}
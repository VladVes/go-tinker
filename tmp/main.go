package main

import (
	"fmt"

	p1 "github.com/VladVes/go-tinker/v2/tmp/pack1"
	p12 "github.com/VladVes/go-tinker/v2/tmp/pack1/v2"
)

func main() {
	greeting := p1.Greeting(p12.A)
	fmt.Println(greeting)
}

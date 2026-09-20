package main

import (
	// "bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func echo(args []string) {
	t := time.Now()
	fmt.Println(strings.Join(os.Args[1:], " "))
	fmt.Println(time.Since(t).Microseconds())
}

func main() {
	echo(os.Args)
}

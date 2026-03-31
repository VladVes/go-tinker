package main

import (
	"fmt"
	"strconv"

	"github.com/VladVes/go-tinker/v2/basics"
	ev "github.com/VladVes/go-tinker/v2/even"
	grtg "github.com/VladVes/go-tinker/v2/greeting"
	clr "github.com/fatih/color"
	// "github.com/VladVes/go-tinker/v2/websrv"
	// "github.com/VladVes/go-tinker/v2/webfibersrv"
)

func main() {
	fmt.Println(grtg.Hello())
	clr.Magenta(grtg.Hello())
	clr.Green(grtg.Hello())
	nums := []int{1, 2, 3, 4, 5, 33, 7, 56, 99, 13, 12}
	for _, v := range ev.Even(nums) {
		clr.Red(strconv.Itoa(v))
	}

	strSet := basics.NewSet[string]()
	fmt.Println(strSet)
	strSet.Add("first")
	strSet.Add("second")
	strSet.Add("third")
	strSet.Add("fourth")
	fmt.Println(strSet)
	fmt.Println(strSet.Has("second"))
	strSet.Remove("second")
	fmt.Println(strSet.Has("second"))

	// Запуск веб сервера с пом. стандартного пакета net/http
	// websrv.Run()

	// Запуск веб-сервера на микрофреймворке fiber
	// webfibersrv.Run()
}

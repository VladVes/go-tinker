package main

import (
	"fmt"
	"strconv"

	ev "github.com/VladVes/go-tinker/v2/even"
	grtg "github.com/VladVes/go-tinker/v2/greeting"
	"github.com/VladVes/go-tinker/v2/websrv"
	clr "github.com/fatih/color"
)

func main() {
	fmt.Println(grtg.Hello())
	clr.Magenta(grtg.Hello())
	clr.Green(grtg.Hello())
	nums := []int{1, 2, 3, 4, 5, 33, 7, 56, 99, 13, 12}
	for _, v := range ev.Even(nums) {
		clr.Red(strconv.Itoa(v))
	}
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	// тело ответа — это массив байт
	// 	w.Write([]byte("Hello world!"))
	// })

	// // запускаем веб-приложение для обработки запросов
	// http.ListenAndServe(":80", nil)

	websrv.Run()
}

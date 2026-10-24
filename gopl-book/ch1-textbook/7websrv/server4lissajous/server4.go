package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func main() {
	server4()
}

func server4() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/lissajous", lissajousHandler)
	log.Fatal(http.ListenAndServe("localhost:8181", nil))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	// в одном запросе нельзя отправить и текст и банарные данные
	// но ответ можно сформировать не просто в виде текста, а в виде html
	// что бы браузер сам сделал дополнительный запрос и ниже подгрузилась картинка
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "%s %s %s\n", r.Method, r.Host, r.Proto)
	for k, v := range r.Header {
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	}
	if err := r.ParseForm(); err != nil {
		log.Print(err)
	}
	for k, v := range r.Form {
		fmt.Fprintf(w, "Form[%q] = %q\n", k, v)
	}
	// получаем число циклов как параметр и передаём в следующий урл подгружающий картинку
	cycles := r.URL.Query().Get("cycles")
	cyclesParam := ""
	if cycles != "" {
		cyclesParam = fmt.Sprintf("?cycles=%s", cycles)
	}

	// при обработке html браузер сам сходит за картинокй по пути в тэге
	fmt.Fprintf(w, `<img src="/lissajous%s" alt="Lissajous" />`, cyclesParam)

}

func lissajousHandler(w http.ResponseWriter, r *http.Request) {
	cycles := 5 //default
	val := r.URL.Query().Get("cycles")
	if val != "" {
		var err error
		cycles, err = strconv.Atoi(val)
		fmt.Printf("URL: %s, cycles: %d\n", r.URL, cycles)
		if err != nil {
			http.Error(w, "invalid cycles parametr", http.StatusBadRequest)
			return
		}
	}
	// выстваляем верный заголовое для бинарника
	w.Header().Set("Content-Type", "image/gif")
	Lissajous(w, cycles)
}

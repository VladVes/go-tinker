package main

import (
	"fmt"
	"log"
	"net/http"
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

	// при обработке html браузер сам сходит за картинокй по пути в тэге
	fmt.Fprintln(w, `<img src="/lissajous" alt="Lissajous" />`)

}

func lissajousHandler(w http.ResponseWriter, r *http.Request) {
	// выстваляем верный заголовое для бинарника
	w.Header().Set("Content-Type", "image/gif")
	Lissajous(w)
}

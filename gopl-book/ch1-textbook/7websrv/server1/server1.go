// log
// http.HandleFunc("/", handler)
// log.Fatal(http.ListenAndServe("localhost:8181", nil))
// func handler(w http.ResponseWriter, r *http.Request)
// r.URL.Path

package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	Server1()
}

func Server1() {
	http.HandleFunc("/", handler)
	// log.Fatal ждет ошибку которую может вернуть ListenAndServe а пока её нет процесс продолжает крутится
	// для учебной задачи нет смысла отдельно обрабатывать ошибку и пытаться восстанавливаться
	log.Fatal(http.ListenAndServe("localhost:8181", nil))
	// nil — второй аргумент, это handler. Когда передаётся nil,
	// Go использует глобальный мультиплексор (http.DefaultServeMux),
	// к которому ранее был привязан обработчики через http.HandleFunc.
	// То есть вся логика маршрутов уже «записана» в него,
	// и здесь просто говоришь: «слушай и используй то, что уже настроено».
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
}

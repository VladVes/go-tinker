// минимальный "echo"-сервер со счетчиком запросов

package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

var mu sync.Mutex
var count int

func main() {
	Server2()
}

func Server2() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/count", countHandler)
	log.Fatal(http.ListenAndServe("localhost:8181", nil))
}

// инкрементируем счётчик и отдаём в ответ путь из урл
func rootHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	log.Print("root handler")
	count++
	log.Printf("Counter: %d\n", count)
	mu.Unlock()
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
}

// отдаём в ответ строку со счётчиком
func countHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	log.Print("counter handler")
	log.Printf("Counter: %d\n", count)
	fmt.Fprintf(w, "Count %d\n", count)
	mu.Unlock()
}

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
	log.Fatal(http.ListenAndServe("localhost:8181", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
}

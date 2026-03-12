package websrv

import (
	"fmt"
	"net/http"
)

func handleGet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from Go"))
}

func Run() {
	http.HandleFunc("/", handleGet)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

package websrv

import (
	"fmt"
	"net/http"
)

const (
	greeting = "Hello from Go"
	about    = "This web server is written in Go!"
)

func handleGetRoot(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(greeting))
}

func handeGetAbout(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(about))
}

func handleGetCourses(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	var pageCourses string
	switch page {
	case "":
		pageCourses = "Enter course number!"
	case "1":
		pageCourses = "How to write your first \"Hello world\" in Go..."
	case "2":
		pageCourses = "How does garbage colletor works?"
	default:
		pageCourses = "Under construction"
	}
	w.Write([]byte(pageCourses))
}

func Run() {
	http.HandleFunc("/", handleGetRoot)
	http.HandleFunc("/about", handeGetAbout)
	http.HandleFunc("/courses", handleGetCourses)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

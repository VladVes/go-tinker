package websrv

import (
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

const (
	greeting = "Hello from Go"
	about    = "This web server is written in Go!"
)

func Run() {
	var logger = logrus.New()
	logger.SetOutput(os.Stdout)
	port := "8080"

	handleGetRoot := func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(greeting))
		if err != nil {
			// Ошибка логируется функцией WithError
			logrus.WithError(err).Error("get roor handler error")
		}
	}

	handeGetAbout := func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(about))
		if err != nil {
			// Ошибка логируется функцией WithError
			logrus.WithError(err).Error("about handler error")
		}
	}

	handleGetCourses := func(w http.ResponseWriter, r *http.Request) {
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
		_, err := w.Write([]byte(pageCourses))
		if err != nil {
			// Ошибка логируется функцией WithError
			logrus.WithError(err).Error("courses handler error")
		}
	}

	http.HandleFunc("/", handleGetRoot)
	http.HandleFunc("/about", handeGetAbout)
	http.HandleFunc("/courses", handleGetCourses)
	// Дополнительная информация передается функцией WithFields
	logrus.WithFields(logrus.Fields{
		"port": port,
	}).Info("Starting a web-server on port")
	// err := http.ListenAndServe(":"+port, nil)
	// if err != nil {
	// 	fmt.Println(err)
	// 	panic(err)
	// }
	logrus.Fatal(http.ListenAndServe(":"+port, nil))
}

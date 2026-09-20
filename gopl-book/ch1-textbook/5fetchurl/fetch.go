// net/http

// http.Get(url)
// os.Exit(1)
// b, err := os.ReadAll(resp.Body)
// resp.Body.Close()

// b, err := io.Copy(os.Stdout, resp.Body)
// fmt.Printf("HTTP code: %s\n\n", resp.Status)

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func fetch1() {
	for _, url := range os.Args[1:] {
		// Функция http.Get выполняет HTTP-запрос и при отсутствии ошибок возвращает
		// результат в структуре resр.
		resp, err := http.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
			// В случае любой ошибки os.Exit(1) завершает
			// работу процесса с кодом состоя­ния 1.
			os.Exit(1)
		}
		// Поле Body этой структуры resp содержит ответ сервера
		// в виде потока, доступного для чтения.
		// io.ReadAll считывает весь ответ; резуль­тат
		// сохраняется в переменной Ь.
		b, err := io.ReadAll(resp.Body)
		// Поток Body закрывается для предотвращения утечки ресурсов
		resp.Body.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
			os.Exit(1)
		}
		// функция Printf записывает ответ в стандартный вывод
		fmt.Printf("%s", b)
	}
}

func fetch2() {
	for _, url := range os.Args[1:] {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
			os.Exit(1)
		}
		// Вызов функции io.Copy(dst, src) выполняет чтение src и
		// запись в dst. Воспользуйтесь ею вместо io.ReadAll для копирования
		// тела ответа в поток os.Stdout без необходимости выделения
		// достаточно большого для хранения всего ответа буфера
		b, err := io.Copy(os.Stdout, resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: copy %s: %v\n", url, err)
		}
		fmt.Printf("Bytes copied: %v\n", b)
	}
}

func fetch3() {
	for _, url := range os.Args[1:] {
		normUrl := url
		// добавить к url http:// если нет
		if !strings.HasPrefix(url, "http") {
			normUrl = "http://" + url
		}
		fmt.Printf("Fetching URL: %s...\n", normUrl)
		resp, err := http.Get(normUrl)
		// выводить статус запроса
		fmt.Printf("HTTP status code: %s\n\n", resp.Status)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
			os.Exit(1)
		}
		b, err := io.Copy(os.Stdout, resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: copy %s: %v\n", url, err)
		}
		fmt.Printf("Bytes copied: %v\n", b)
	}
}

func main() {
	fetch3()
}

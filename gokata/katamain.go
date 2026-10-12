package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	fetchAll(os.Args[1:])
}

func fetchAll(urls []string) {
	ch := make(chan string)
	for _, url := range urls {
		go fetch(url, ch)
	}
	fmt.Println("Results:")
	for range urls {
		fmt.Println(<-ch)
	}
}

func fetch(url string, ch chan<- string) {
	prot := "http"
	if !strings.HasPrefix(url, prot) {
		url = prot + "s://" + url
	}
	timer := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprintf("fetch: %v\n", err)
		return
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		ch <- fmt.Sprint("Reading body: %v", err)
		resp.Body.Close()
		return
	}
	filename := time.Now().String()
	wferr := os.WriteFile(filename, data, 0777)
	if wferr != nil {
		ch <- fmt.Sprint("writing body to file: %v", err)
		resp.Body.Close()
		return
	}
	nbytes, err := io.Copy(io.Discard, resp.Body)
	resp.Body.Close()
	if err != nil {
		ch <- fmt.Sprint("Coping body: %v", err)
		return
	}
	t := time.Since(timer).Seconds()
	// result := string(data)
	ch <- fmt.Sprintf("Fetched url:\t%s\nStatus:\t\t%s\nSize:\t\t%db\nFile:\t\t%s\nElapsed:\t%.3fs\n", url, resp.Status, nbytes, filename, t)
}

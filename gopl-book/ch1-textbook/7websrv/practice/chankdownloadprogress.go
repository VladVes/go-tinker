package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No url recieved")
		os.Exit(1)
	}
	fetchWithProgress(os.Args[1])
}

func fetchWithProgress(url string) {
	normUrl := url
	if !strings.HasPrefix("http", url) {
		normUrl = "https://" + url
	}
	fmt.Printf("Fetching url: %s\n", normUrl)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fetch url: %v", err)
		os.Exit(1)
	}
	buff := make([]byte, 2049)
	for {
		n, err := resp.Body.Read(buff)
		fmt.Printf("Bytes read\t%d\n", n)
		fmt.Println(string(buff))
		if err != nil {
			break
		}
	}

}

// ch := make(chan string)
// fmt.Println(<-ch) // получение из канала ch

// func fetch(url string, ch chan<- string)
// ch <- fmt.Sprint(err) // отправка в канал ch
// nbytes, err := io.Copy(io.Discard, resp.Body)

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	start := time.Now()
	// Канал является механизмом связи, который позволяет одной go-подпрограмме
	// передавать значения определенного типа другой go-подпрограмме.
	ch := make(chan string) // канал для передачи данных типа string
	for _, url := range os.Args[1:] {
		go fetch(url, ch) // создание и запуск go-подпрограммы (горутины)
	}
	fmt.Println("Fetching finished with results:")
	for range os.Args[1:] {
		fmt.Println(<-ch) // получение из канала ch
	}
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func fetch(url string, ch chan<- string) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err) // отправка в канал ch
		return
	}
	// Функция io.Copy считывает тело ответа и игнорирует его,
	// записывая в выходной поток io.Discard.
	// Возвращает количество байтов и информа­цию о происшедших ошибках
	nbytes, err := io.Copy(io.Discard, resp.Body)
	resp.Body.Close() // исключаем утечки ресурсов
	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v\n", url, err)
		return
	}
	secs := time.Since(start).Seconds()
	// При получении каждого результата fetch отправляет итоговую строку в канал ch
	// Когда одна go-подпрограмма пытается отправить или получить информацию по каналу,
	// она блокируется, пока другая go-подпрограмма пытается выполнить соответствующие
	// операции получения или отправки в этот же канал,
	// и после передачи информации обе go-подпрограммы продолжают работу.
	ch <- fmt.Sprintf("%.2fs %d %s\n", secs, nbytes, url)
	// В данном примере каждая функция fetch отправляет значение
	// (ch <- expression) в канал ch, и main получает их все (<- ch).
}

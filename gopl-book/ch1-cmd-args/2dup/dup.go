// os
// bufio

// input := bufio.NewSacnner(file) // file *os.File (os.Stdin)
// input.Scan()
//
// f, err := os.Open(fileName) 
// fmt.Fprintf(os.Stderr, "dup2 error with %v: %v\n", fileName, err)  
// f.Close()
//
package main

import (
	"bufio"
	"fmt"
	"os"
)

func dup1() {
	counts := make(map[string]int)
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		counts[input.Text()]++
	}
	// Примечание: игнорируем потенциальные
	// ошибки из input.Err()
	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%d\t%s\n", n, line)
		}
	}
}

// Следующая версия dup может как выполнять чтение стандартного ввода,
// так и работать со спис­ком файлов, используя os.Open для их открывания:
func dup2() {
	counts := make(map[string]int)
	files := os.Args[1:]
	if len(files) == 0 {
		countLines(os.Stdin, counts)// вызов countLines предшествует объявлению этой функ­ции.
		// Функции и другие объекты уровня пакета могут быть объявлены в любом по­рядке.
	} else {
		for _, arg := range files {
			f, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "dup2: %v\n", err)
				continue
			}
			countLines(f, counts)
			f.Close()
		}
	}
	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%d\t%s\n", n, line)
		}
	}

}

func countLines(f *os.File, counts map[string]int) {
	input := bufio.NewScanner(f)
	for input.Scan() {
		// на каждой итерации Scan читает строку из файла (которым может быть как открытый файл так и stdin)
		// если строка прочитана, то Scan возвращает true и помещает прочитанную строку в свою структуру, далее 
		// цикл получает результат true и выполняет своё тело и т.д. пока scan не вернёт fals из-за конца файла или ошибки
		counts[input.Text()]++
	}
}

// func main_() {
// dup1() // <Ctrl+D> - конец вводда в linux/machOs, либо ./dup < data1.txt
// }

func main() {
	dup2()
}

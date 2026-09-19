// os
// bufio
//

// input := bufio.NewSacnner(file) // file *os.File (os.Stdin)
// input.Scan()
//
// f, err := os.Open(fileName)
// fmt.Fprintf(os.Stderr, "dup2: %v\n", err)
// f.Close()
//
// data, err := os.ReadFile(filename)

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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
	counts := make(map[string]map[string]int)
	files := os.Args[1:]
	if len(files) == 0 {
		countLines(os.Stdin, counts, "stdin") // вызов countLines предшествует объявлению этой функ­ции.
		// Функции и другие объекты уровня пакета могут быть объявлены в любом по­рядке.
	} else {
		for _, filename := range files {
			f, err := os.Open(filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "dup2: %v\n", err)
				continue
			}
			countLines(f, counts, filename)
			f.Close()
		}
	}
	for filename, counter := range counts {
		doubles := getDobles(counter)
		if len(doubles) > 0 {
			fmt.Printf("File: %s\n", filename)
			for line, n := range doubles {
				fmt.Printf("%d\t%s\n", n, line)
			}
		}
	}

}

func countLines(f *os.File, counts map[string]map[string]int, filename string) {
	input := bufio.NewScanner(f)
	counts[filename] = make(map[string]int)
	for input.Scan() {
		// на каждой итерации Scan читает строку из файла (которым может быть как открытый файл так и stdin)
		// если строка прочитана, то Scan возвращает true и помещает прочитанную строку в свою структуру, далее
		// цикл получает результат true и выполняет своё тело и т.д. пока scan не вернёт fals из-за конца файла или ошибки
		counts[filename][input.Text()]++
	}
}

func getDobles(counter map[string]int) map[string]int {
	result := make(map[string]int)
	for k, v := range counter {
		if v > 1 {
			result[k] = v
		}
	}
	return result
}

func dup3() {
	counter := make(map[string]int)
	for _, filename := range os.Args[1:] {
		data, err := os.ReadFile(filename)
		// Функция ReadFile возвращает байтовый срез, который должен быть преобразо­ван в string так,
		// чтобы его можно было разбить с помощью функции strings.Split
		if err != nil {
			fmt.Fprintf(os.Stderr, "dup3: %v\n", err)
			continue
		}
		for _, line := range strings.Split(string(data), "\n") {
			counter[line]++
		}
	}

	for line, n := range counter {
		if n > 1 {
			fmt.Printf("%d\t\t%s\n", n, line)
		}
	}
}

// func main_() {
// dup1() // <Ctrl+D> - конец вводда в linux/machOs, либо ./dup < data1.txt
// }

func main() {
	dup2()
}

// func main() {
// 	dup3()
// }

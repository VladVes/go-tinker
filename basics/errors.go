package basics

import (
	"fmt"
	"os"
)

// 0.2 Написать небольшой пример где функция ReadFile("file.txt")
// из пакета os читает файл и возвращает ошибку, которую нужно проверить и вернуть сообщение
// об ошибке используя специальную функцию (из пакета предназначенного для работы с выводом).
func GetFileContent(fname string) ([]byte, error) {
	res, err := os.ReadFile(fname)
	if err != nil {
		return []byte{}, fmt.Errorf("file read: %w", err)
	}
	return res, nil
}

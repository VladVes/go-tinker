package basics

import "errors"

// 9. Какой особенностью обладают функции в Go в плане возвращения значений?
// Написать функцию divide принимающую два значения x y и возвращающая результат деления x на y.
// Как должна обрабатываться ситуация когда в кач-ве y передаётся 0? Как можно в Go создать объект ошибки?

func DivideByZero(x, y int) (float64, error) {
	if y == 0 {
		return 0, errors.New("divide by zero!")
	}
	return float64(x / y), nil
}

package pack1

import "fmt"

const (
	A = "Alice"
	B = "Ivan"
)

func Greeting(name string) string {
	return fmt.Sprintf("Hello, %s!", name)
}

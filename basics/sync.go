package basics

import (
	"fmt"
	"sync"
)

// 154. Какой есть механизм в Go позволяющий запланировать какие-то действия перед завершением функции?
func DeferDemo() {
	defer fmt.Println("First code")
	defer fmt.Println("Second code")
	fmt.Println("Main code")
}

// Main code
// Second code
// First code
// defer выполниться перед выходом из фукнции в любом случае (даже при ошибке)

// Почему такой подход особенно важен при работе с мутексами?
type Order struct {
	Id   string
	Name string
}
type OrderStorage struct {
	m      sync.Mutex
	Orders map[string]Order
}

func (o *OrderStorage) CreateOrder(order Order) (string, error) {
	o.m.Lock()
	defer o.m.Unlock()
	// можно просто забыть в конце вызвать Unlock
	// либо выполнение до вызова Unlock не дойдёт из-за возникшей ошибке где-то в середине кода функции
	// может быть несколько точек выхода тогда нужно будлировать Unlock что тоже можно забыть сделать
	// отложенный вызов с defer в самом начале когда - удобно и надёжно т.к. отработает в любом слвуяае перед
	// выходом из функции
	o.Orders[order.Id] = order
	return order.Id, nil
}

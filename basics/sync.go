package basics

import (
	"errors"
	"fmt"
	"sync"
)

var errOrderNotFound = errors.New("Order nor found")

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

// 155. Зачем инкапуслировать мутекс в структуру? (из примера OrderStorage) В чём плюсы такого подхода?
func (o *OrderStorage) GetOrder(id string) (Order, error) {
	o.m.Lock()
	defer o.m.Unlock()
	order, ok := o.Orders[id]
	if !ok {
		return Order{}, errOrderNotFound
	}
	return order, nil
}

// Kогда мутекс является частью структуры это позволяет иметь большую гибкость и удобство работы с ним:
// - инкапсуляция что бы пользователь структуры не думал о необходимости синхронизации
// - при вызове методов безопасность автоматически обеспечивается, не нужно верхнеуровнево оборачивать их в Lock() и Unlock()
// - в каждом методе при необходимости можно применить свою вариацию синхронизации
// - исключается риск забыть Unlock() т.к. всё происходит автоматом внутри методов

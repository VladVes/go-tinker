package main

import (
	"fmt"
	"sync"
)

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
	o.Orders[order.Id] = order
	return order.Id, nil
}

func main() {
	defer fmt.Println("first code")
	defer fmt.Println("second code")
	fmt.Println("Main code")
}

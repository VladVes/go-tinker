package entities

import (
	"errors"
	"fmt"

	"github.com/VladVes/go-tinker/v2/data/schemas"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Пример простейшей логики работы сознания и получения заказа реализованный по принципу слоёной архитекутуры
// Слой обработчика -> слой бизне-логики (тут опущен) -> слой хранилица

type (
	// Слой хранилища - хранение заказа (для примера просто в памяти)
	// сущность заказ
	Order struct {
		ID         string
		UserID     int64
		ProductIDs []int64
	}
	// Само хранилище описывается структурой (мап строка - id заказа: структ. заказа):
	OrderStorage struct {
		Orders map[string]Order
	}

	// Слой обработчика
	// Описывается стркутурой но не связывается на прямую с хранилищем
	// для сохранения гибкости т.е. опосредовано через контрак - интерфейс которому должно соответствовать хранилище.

	// Структура обработчика в дальнейшем должна получить методы GetOrder и CreateOrder которые и будут являтся хэндлерами и связываться
	// с маршрутами GET "/order" и POST "/order" и уже сами хэндлеры после парсинга тела запроса или параметров
	// будут обращаться к полю storage и вызывать на нём соответсвующие (даже одноименные как в этом примере)
	// методы для работой с хранилищем описанные интерфейсом.
	OrderHandler struct {
		Storage OrderCreatorGetter
	}

	OrderCreatorGetter interface {
		GetOrder(orderID string) (Order, error)
		CreateOrder(o Order) (string, error)
	}
)

// Реализация методов хранилища соответсвующиx контратку
func (o *OrderStorage) GetOrder(orderID string) (Order, error) {
	order, ok := o.Orders[orderID]
	if !ok {
		errMsg := fmt.Sprintf("Order with id %s not found", orderID)
		return Order{}, errors.New(errMsg)
	}
	return order, nil
}

func (o *OrderStorage) CreateOrder(order Order) (string, error) {
	o.Orders[order.ID] = order

	return order.ID, nil
}

// Реализация методов обработчика
func (h *OrderHandler) GetOrder(c *fiber.Ctx) error {
	id := c.Params("id")

	order, err := h.Storage.GetOrder(id)
	if err != nil {
		return fmt.Errorf("get order: %w", err)
	}

	return c.JSON(schemas.GetOrderResponse(order))
}

func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	var request schemas.CreateOrderRequest
	if err := c.BodyParser(&request); err != nil {
		return fmt.Errorf("body parser: %w", err)
	}

	order := Order{
		ID:         uuid.New().String(),
		UserID:     request.UserID,
		ProductIDs: request.ProductIDs,
	}

	id, err := h.Storage.CreateOrder(order)
	if err != nil {
		return fmt.Errorf("create order: %w", err)
	}

	return c.JSON(schemas.CreateOrderResponse{
		ID: id,
	})
}

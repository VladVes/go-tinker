package entities

import (
	"fmt"
	"sync"

	"github.com/VladVes/go-tinker/v2/data/schemas"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Пример CRUD по сущности Employee c хранением в памяти
type Employee struct {
	ID    string
	Email string
	Role  string
}

// Добавляем возможность блокировать изменение данных разными горутинами одновременно
type EmployeeStorageInMemory struct {
	mu       sync.Mutex
	Employee map[string]Employee
}

func (s *EmployeeStorageInMemory) Create(e Employee) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	e.ID = uuid.New().String()
	s.Employee[e.ID] = e

	return e.ID, nil
}

// Связь обработчика с хранилищем через интерфейс
type EmployeeStorage interface {
	Create(e Employee) (string, error)
}

type EmployeeHanlder struct {
	Storage EmployeeStorage
}

// Создание
func (h *EmployeeHanlder) CreateEmployee(c *fiber.Ctx) error {
	var req schemas.CreateEmployeeRequest
	if err := c.BodyParser(&req); err != nil {
		return fmt.Errorf("body parser: %w", err)
	}

	id, err := h.Storage.Create(Employee{
		Email: req.Email,
		Role:  req.Role,
	})
	if err != nil {
		return fmt.Errorf("create in storage: %w", err)
	}

	return c.JSON(schemas.CreateEmployeeResponse{ID: id})
}

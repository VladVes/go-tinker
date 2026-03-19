package entities

import (
	"errors"
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

// Методы хранилища в соответствии с интерфейсом EmployeeStorage
func (s *EmployeeStorageInMemory) Create(e Employee) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	e.ID = uuid.New().String()
	s.Employee[e.ID] = e

	return e.ID, nil
}

func (s *EmployeeStorageInMemory) List() []Employee {
	// Инициализируем массив с размером равным количеству
	// всех сотрудников в хранилище
	employees := make([]Employee, 0, len(s.Employee))

	// заполняем новый срез и возвращаем
	for _, e := range s.Employee {
		employees = append(employees, e)
	}

	return employees
}

func (s *EmployeeStorageInMemory) Get(id string) (Employee, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	e, ok := s.Employee[id]
	if !ok {
		return Employee{}, errors.New("employee not found")
	}

	return e, nil
}

// Связь обработчика с хранилищем через интерфейс
type EmployeeStorage interface {
	Create(e Employee) (string, error)
	List() []Employee
	Get(id string) (Employee, error)
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

// Получение всех и одного по id
// объекты ответов:
type (
	ListEmployeesResponse struct {
		Employees []EmployeePayload `json:"employees"`
	}

	GetEmployeeResponse struct {
		// анонимное поле (встраивание структуры) позволяет создавать композицию
		EmployeePayload
	}

	// можно было бы использовать уже имеющуюся структруру 
	// Employee для ответа, но для правильней и для ответа использовать отделньую
	// т.к. поля модели и ответа часто не совпадают
	EmployeePayload struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Role  string `json:"role"`
	}
)

// Получение списка сотрудников
func (h *EmployeeHanlder) GetEmployees(c *fiber.Ctx) error {
	// Получаем список всех сотрудников из хранилища
	employees := h.Storage.List()
	// Формируем ответ
	resp := ListEmployeesResponse{
		Employees: make([]EmployeePayload, len(employees)),
	}
	for i, e := range employees {
		// создание экземпляра EmployeePayload путём приведения
		resp.Employees[i] = EmployeePayload(e)
	}

	return c.JSON(resp)
}

// Получение сотрудника по id
func (h *EmployeeHanlder) GetEmployeeByID(c *fiber.Ctx) error {
	id := c.Params("id")

	e, err := h.Storage.Get(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString(err.Error())
	}

	// т.к. структруа ответа имеет встроенну структуру
	// ответ формируется так
	resp := GetEmployeeResponse{
		EmployeePayload{
			ID:    e.ID,
			Email: e.Email,
			Role:  e.Role,
		}}
	// но благодаря приведению типов
	// можно сделать короче:
	resp = GetEmployeeResponse{EmployeePayload(e)}

	return c.JSON(resp)
}

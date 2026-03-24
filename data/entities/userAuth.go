package entities

import (
	"errors"
	"fmt"

	"github.com/VladVes/go-tinker/v2/data/schemas"
	"github.com/gofiber/fiber/v2"
)

var (
	ErrBadCredentials = errors.New("email or password is incorrect")
	JwtSecretKey      = []byte("some-secret-key")
)

type (
	User struct {
		Email    string
		Name     string
		password string
	}
	AuthHandler struct {
		Storage *AuthStorage
	}
	AuthStorage struct {
		Users map[string]User
	}
)

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var regReq schemas.RegisterRequest
	if err := c.BodyParser(&regReq); err != nil {
		return fmt.Errorf("Body parser: %w", err)
	}

	// Проверяем, что пользователь с таким email еще не зарегистрирован
	if _, exists := h.Storage.Users[regReq.Email]; exists {
		return errors.New("The user already exists")
	}

	// Сохраняем в хранилище (в примере - память) нового зарегистрированного пользователя
	h.Storage.Users[regReq.Email] = User{
		Email: regReq.Email,
		Name:  regReq.Name,
		// для упрощения примера пароль не хэшируется
		password: regReq.Password,
	}
	return c.SendStatus(fiber.StatusCreated)
}

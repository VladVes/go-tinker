package entities

import (
	"errors"
	"fmt"
	"time"

	"github.com/VladVes/go-tinker/v2/data/schemas"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
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
	UserHandler struct {
		Storage *AuthStorage
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
		// fiber.StatusConflict
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

// Обработчик HTTP-запросов на вход в аккаунт
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	req := schemas.LoginRequest{}
	if err := c.BodyParser(&req); err != nil {
		return fmt.Errorf("body parser: %w", err)
	}

	// Ищем пользователя в памяти приложения по электронной почте
	user, exists := h.Storage.Users[req.Email]
	if !exists {
		return ErrBadCredentials
	}
	// Если пользователь найден, но у него другой пароль, возвращаем ошибку
	if user.password != req.Password {
		return ErrBadCredentials
	}

	// Генерируем JWT-токен для пользователя,
	// который он будет использовать в будущих HTTP-запросах
	// С пом. пакета jwt-go токен создаётся в три этапа:
	// 1. Генерируем полезные данные, которые будут храниться в токене
	payload := jwt.MapClaims{
		"sub": user.Email,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}
	// 2. Создаем новый JWT-токен и подписываем его по алгоритму HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	// 3. Подписываем и сериализуем
	t, err := token.SignedString(JwtSecretKey)

	if err != nil {
		logrus.WithError(err).Error("JWT token signin")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(schemas.LoginResponse{AccessToken: t})
}

// Получение данных о пользователе из токена
// передаваемого в заголовке Authorization
const (
	ContextKeyUser = "user"
)

func jwtPayloadFromRequest(c *fiber.Ctx) (jwt.MapClaims, bool) {
	// доразобраться как это работает:
	jwtToken, ok := c.Context().Value(ContextKeyUser).(*jwt.Token)
	if !ok {
		logrus.WithFields(logrus.Fields{
			"jwt_token_context_value": c.Context().Value(ContextKeyUser),
		}).Error("wrong type of JWT token in context")
		return nil, false
	}

	payload, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		logrus.WithFields(logrus.Fields{
			"jwt_token_claims": jwtToken.Claims,
		}).Error("wrong type of JWT token claims")
		return nil, false
	}

	return payload, true
}

// Обработчик HTTP-запросов на получение профиля пользователя
func (h *UserHandler) Profile(c *fiber.Ctx) error {
	jwtPayload, ok := jwtPayloadFromRequest(c)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	userInfo, ok := h.Storage.Users[jwtPayload["sub"].(string)]
	if !ok {
		return errors.New("user not found")
	}

	return c.JSON(schemas.ProfileResponse{
		Email: userInfo.Email,
		Name:  userInfo.Name,
	})
}

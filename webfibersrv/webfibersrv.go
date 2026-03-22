package webfibersrv

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/VladVes/go-tinker/v2/data"
	"github.com/VladVes/go-tinker/v2/data/entities"
	"github.com/VladVes/go-tinker/v2/data/schemas"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const profileUnknown = "unknown"

func Run() {
	app := fiber.New()
	app.Get("/address", func(c *fiber.Ctx) error {
		return c.SendString("Hello, go! Don't gitve up! Fight!\n")
	})

	// -------------------------------------Query---------------------------------------------------------------
	// работа с параметрами запроса
	app.Get("/profile", func(c *fiber.Ctx) error {
		profileId := c.Query("profile_id", profileUnknown)
		if profileId == "" {
			profileId = profileUnknown
		}
		if profileId == profileUnknown {
			return c.Status(fiber.StatusUnprocessableEntity).SendString("porofile_id is required")
		}

		return c.SendString(fmt.Sprintf("user Profile ID is: %s", profileId))
	})

	// ------------------------------------Dynamin route----------------------------------------------------------------
	// пример динамического роутинга
	app.Get("/likes/:postId", func(c *fiber.Ctx) error {
		postId := c.Params("postId")
		likes, ok := data.PostLikes[postId]
		if !ok {
			logrus.WithFields(logrus.Fields{
				"postId": postId,
			}).Info(fmt.Sprintf("A post with this id is not found: %s", postId))
			return c.Status(fiber.StatusNotFound).SendString("no post Id\n")
		}

		return c.SendString(fmt.Sprintf("Post id: %s, Likes: %s\n", postId, strconv.FormatInt(likes, 10)))
	})

	app.Post("/likes/:postId", func(c *fiber.Ctx) error {
		postId := c.Params("postId")

		data.PostLikes[postId]++
		likes := data.PostLikes[postId]

		status := fiber.StatusOK
		if likes == 1 {
			status = fiber.StatusCreated
		}
		return c.Status(status).SendString(fmt.Sprintf("Post id: %s, Likes: %s\n", postId, strconv.FormatInt(likes, 10)))

	})

	// -------------------------------JSON-----------------------------------------------------------------
	// Пример десериализации данных передаваемых в теле запсроса
	// и сериализации отправляемого ответа
	app.Post("/logs", func(c *fiber.Ctx) error {
		var request schemas.CreateLogEntryRequest
		if err := c.BodyParser(&request); err != nil {
			return fmt.Errorf("body parser: %w", err)
		}

		logEntry := schemas.LogEntry{
			ID:        uuid.New().String(),
			Message:   request.Message,
			Level:     request.Level,
			Timestamp: request.Timestamp,
		}

		// Упрощенное хранение в памяти приложения
		data.Logs = append(data.Logs, logEntry)

		return c.JSON(schemas.CreateLogEntryResponse{
			ID: logEntry.ID,
		})
	})

	// ----------------------------------POST JSON---------------------------------------------------------
	// Пример обработки POST запроса в теле которого передаётся отсортированый массив чисел
	// и искомое число, в ответе в теле возвращается JSON с полем target_index с индексом искомого числа в массиве
	// В примере используются две стркутруты - схемы для сериализации и десериализации. Схемы имеют тегированные поля
	// в которых содержится описание того из каких полей JSON что использовать для каких полей структуры и на оборот
	//  какой JSON должен получиться из каких полей
	app.Post("/searchIndex", func(c *fiber.Ctx) error {
		req := schemas.SearchIndexReq{}
		resp := schemas.SearchIndexResp{
			TargetIndex: -1,
		}

		if err := c.BodyParser(&req); err != nil {
			resp.Error = "Invalid JSON"
			logrus.WithError(err).Error(resp.Error)
			return c.Status(fiber.StatusBadRequest).JSON(resp)
		}

		targetIndex := slices.Index(req.Numbers, req.Target)
		if targetIndex == -1 {
			resp.Error = "Target not found"
			logrus.WithFields(logrus.Fields{
				"target":  req.Target,
				"numbers": req.Numbers,
			}).Info(resp.Error)
			return c.Status(fiber.StatusNotFound).JSON(resp)
		}

		resp.TargetIndex = targetIndex
		return c.Status(fiber.StatusOK).JSON(resp)
	})

	// -----------------------------------Layered arch-----------------------------------------------------------------
	// Пример простейшей логики создания и получения заказа реализованный по принципу слоёной архитекутуры
	// Слой обработчика -> слой бизнеc-логики (тут опущен) -> слой хранилица
	// определены в пакете entities
	orderHandler := &entities.OrderHandler{
		Storage: &entities.OrderStorage{
			Orders: data.Orders,
		},
	}
	app.Post("/orders", orderHandler.CreateOrder)
	app.Get("/orders/:id", orderHandler.GetOrder)

	// -------------------------------------CRUD------------------------------------------------------------
	// Пример CRUD по сущности Employee c хранением в памяти
	//
	employeeHandler := &entities.EmployeeHanlder{
		Storage: &entities.EmployeeStorageInMemory{
			Employee: data.Employee,
		},
	}

	app.Post("/employees", employeeHandler.CreateEmployee)
	app.Get("/employees", employeeHandler.GetEmployees)
	app.Get("/employees/:id", employeeHandler.GetEmployeeByID)
	app.Patch("/employees/:id", employeeHandler.UpdateEmployee)
	app.Delete("/employees/:id", employeeHandler.DeleteEmployee)

	// ------------------------------------Validation-------------------------------------------------------------
	// Пример простой ручной валидации запроса по созданию поста

	app.Post("/post", func(c *fiber.Ctx) error {
		// CreatePostReq определена в Validators.go
		var req CreatePostReq
		if err := c.BodyParser(&req); err != nil {
			return fmt.Errorf("body parser: %w", err)
		}

		// метод определен в Validators.go, проверяет поля структруры и возвращает ошибку c определенным текстом если поле не соответсвтует
		err := req.Validate()
		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).SendString(err.Error())
		}
		// @TODO сохранение данных в бд

		return c.SendStatus(fiber.StatusOK)
		// curl --location --request POST 'http://localhost/post' \
		// --header 'Content-Type: application/json' \
		// --data-raw '{"user_id": -1, "text": ""}'
	})

	// Пример использвоания библиотеки go-playground/validator
	// валидирующей поля используя их аннотации
	type CreatePostRequest struct {
		// Описываем правила валидации в аннотациях полей структуры.
		UserID int64  `json:"user_id" validate:"required,min=0"`
		Text   string `json:"text" validate:"required,max=140"`
	}

	validate := validator.New()

	app.Post("/post", func(c *fiber.Ctx) error {
		var req CreatePostRequest
		if err := c.BodyParser(&req); err != nil {
			return fmt.Errorf("body parser: %w", err)
		}

		// Проверка запроса
		err := validate.Struct(req)
		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).SendString(err.Error())
		}

		// @TODO создание поста и запись в бд

		return c.SendStatus(fiber.StatusOK)

	})

	// ----------------------------------------------------------------------------------------------------
	logrus.Fatal(app.Listen(":" + HttpPort))
}

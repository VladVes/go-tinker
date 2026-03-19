package webfibersrv

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/VladVes/go-tinker/v2/data"
	"github.com/VladVes/go-tinker/v2/data/entities"
	"github.com/VladVes/go-tinker/v2/data/schemas"
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

	// Пример десериализации данные передаваемых в теле запсроса
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

	logrus.Fatal(app.Listen(":" + HttpPort))
}

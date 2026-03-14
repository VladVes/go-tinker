package webfibersrv

import (
	"fmt"
	"strconv"

	"github.com/VladVes/go-tinker/v2/data"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const profileUnknown = "unknown"

func Run() {
	app := fiber.New()
	app.Get("/address", func(c *fiber.Ctx) error {
		return c.SendString("Hello, go! Don't gitve up! Fight!")
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
		var request data.CreateLogEntryRequest
		if err := c.BodyParser(&request); err != nil {
			return fmt.Errorf("body parser: %w", err)
		}

		logEntry := data.LogEntry{
			ID:        uuid.New().String(),
			Message:   request.Message,
			Level:     request.Level,
			Timestamp: request.Timestamp,
		}

		// Упрощенное хранение в памяти приложения
		data.Logs = append(data.Logs, logEntry)

		return c.JSON(data.CreateLogEntryResponse{
			ID: logEntry.ID,
		})
	})

	logrus.Fatal(app.Listen(":" + HttpPort))
}

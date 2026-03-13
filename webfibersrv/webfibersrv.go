package webfibersrv

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func Run() {
	app := fiber.New()
	app.Get("/address", func(c *fiber.Ctx) error {
		return c.SendString("Hello, go! Don't gitve up! Figth!")
	})

	logrus.Fatal(app.Listen(":80"))
}

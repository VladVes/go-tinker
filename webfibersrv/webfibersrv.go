package webfibersrv

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

const profileUnknown = "unknown"

func Run() {
	app := fiber.New()
	app.Get("/address", func(c *fiber.Ctx) error {
		return c.SendString("Hello, go! Don't gitve up! Figth!")
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

	logrus.Fatal(app.Listen(":" + HttpPort))
}

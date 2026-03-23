package webfibersrv

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// -------------------------------------Middlewares and Route Groups-----------------------------------------
// Пример простейшей middleware
func TestMiddleware(c *fiber.Ctx) error {
	sayMiddleware := c.Params("middleware")
	if sayMiddleware != "" {
		fmt.Printf("Middleware test: %s", sayMiddleware)
	}
	return c.Next()
}

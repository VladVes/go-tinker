package webfibersrv

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// -------------------------------------Middlewares and Route Groups-----------------------------------------
// Пример простейшей middleware
func TestMiddleware1(c *fiber.Ctx) error {
	fmt.Println("Run Common Middleware!")
	return c.Next()
}

func TestMiddleware2(c *fiber.Ctx) error {
	fmt.Println("Middleware for route group starts with /mw_tests")
	return c.Next()
}

package middleware

import "github.com/gofiber/fiber/v3"

func AllowAdmin(c fiber.Ctx) error {
	if isAdmin, ok := c.Locals("x-admin").(bool); ok {
		if isAdmin {
			return c.Next()
		}
		return c.SendStatus(fiber.StatusForbidden)
	}
	return c.SendStatus(fiber.StatusInternalServerError)

}

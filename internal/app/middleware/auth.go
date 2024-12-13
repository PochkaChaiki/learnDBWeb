package middleware

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

func NewAuthMiddleware(keyFunc jwt.Keyfunc) fiber.Handler {
	return func(c fiber.Ctx) error {
		authString := c.Get("Authorization")
		if len(authString) == 0 {
			return c.SendStatus(fiber.StatusForbidden)
		}

		parts := strings.Split(authString, " ")
		if len(parts) != 2 && parts[0] != "Bearer" {
			return c.SendStatus(fiber.StatusForbidden)
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, keyFunc)
		if err != nil {
			log.Printf("auth middleware error: %v", err)
			return c.SendStatus(fiber.StatusForbidden)
		}

		if !token.Valid {
			return c.SendStatus(fiber.StatusForbidden)
		}

		tokenClaims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		if userId, ok := tokenClaims["sub"].(float64); ok {
			c.Locals("x-user-id", int(userId))
		} else {
			log.Printf("token parse claims user id error")
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		if adminRole, ok := tokenClaims["admin"].(bool); ok {
			c.Locals("x-admin", adminRole)
		} else {
			log.Printf("token parse claims error")
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.Next()
	}
}

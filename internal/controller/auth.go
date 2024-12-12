package controller

import (
	"learnDB/internal/domain"
	"learnDB/internal/service"
	"log"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

type AuthController struct {
	service *service.AuthService
}

func NewAuthController(s *service.AuthService) *AuthController {
	return &AuthController{service: s}
}

func (a *AuthController) Login(c fiber.Ctx) error {
	// I understand that this way totally sucks but i don't have much time
	u := new(domain.User)
	if err := c.Bind().JSON(u); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	switch check, res := a.service.CheckUserCreds(u); res {
	case service.Ok:
		if !check {
			c.JSON(fiber.Map{
				"status":  "error",
				"message": "invalid user credentials",
			})
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		t, ok := a.service.CreateAccessToken(u)
		if !ok {
			return c.SendStatus((fiber.StatusInternalServerError))
		}
		c.JSON(fiber.Map{"access_token": t})
		return c.SendStatus(fiber.StatusOK)
	case service.InternalError:
		return c.SendStatus(fiber.StatusInternalServerError)
	case service.BadRequest:
		return c.SendStatus(fiber.StatusUnauthorized)
	default:
		return c.SendStatus(fiber.StatusNoContent)
	}

}

func (a *AuthController) AuthMiddleware(c fiber.Ctx) error {
	authString := c.Get("Authorization")
	if len(authString) == 0 {
		return c.SendStatus(fiber.StatusForbidden)
	}

	parts := strings.Split(authString, " ")
	if len(parts) != 2 && parts[0] != "Bearer" {
		return c.SendStatus(fiber.StatusForbidden)
	}

	tokenString := parts[1]

	token, err := jwt.Parse(tokenString, a.service.KeyFunc)
	if err != nil {
		log.Printf("auth middleware error: %v", err)
		return c.SendStatus(fiber.StatusForbidden)
	}

	if !token.Valid {
		return c.SendStatus(fiber.StatusForbidden)
	}

	sub, _ := token.Claims.GetSubject()
	c.Locals("x-username", sub)

	return c.Next()
}

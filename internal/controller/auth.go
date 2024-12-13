package controller

import (
	"learnDB/internal/domain"
	"learnDB/internal/service"

	"github.com/gofiber/fiber/v3"
)

type AuthController struct {
	Service *service.AuthService
}

func NewAuthController(s *service.AuthService) *AuthController {
	return &AuthController{Service: s}
}

func (a *AuthController) Login(c fiber.Ctx) error {
	// I understand that this way totally sucks but i don't have much time
	u := new(domain.User)
	if err := c.Bind().JSON(u); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	switch id, res := a.Service.CheckUserCreds(u); res {
	case service.Ok:
		if id == -1 {
			c.JSON(fiber.Map{
				"status":  "error",
				"message": "invalid user credentials",
			})
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		u.Id = id
		t, ok := a.Service.CreateAccessToken(u)
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

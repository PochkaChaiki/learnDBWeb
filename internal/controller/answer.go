package controller

import (
	"learnDB/internal/service"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

// GET /api/answer
// GET /api/answer/{id}
// POST /api/answer --> /api/query
// DELETE /api/answer/{id}

type AnswerController struct {
	service *service.ServiceAnswer
}

func NewAnswerController(s *service.ServiceAnswer) *AnswerController {
	return &AnswerController{service: s}
}

func (cnt *AnswerController) GetAllAnswers(c fiber.Ctx) error {
	switch anses, result := cnt.service.GetAll(); result {
	case service.Ok:
		if err := c.JSON(anses); err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		c.SendStatus(fiber.StatusOK)
	case service.InternalError:
		return c.SendStatus(fiber.StatusInternalServerError)
	default:
		return c.SendStatus(fiber.StatusNoContent)
	}

	return nil
}

func (cnt *AnswerController) GetAnswer(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		log.Fatalf("get answer error: %s", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	switch ans, res := cnt.service.Get(id); res {
	case service.Ok:
		if err := c.JSON(ans); err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		c.SendStatus(fiber.StatusOK)
	case service.InternalError:
		return c.SendStatus(fiber.StatusInternalServerError)
	default:
		return c.SendStatus(fiber.StatusNoContent)
	}

	return nil
}

func (cnt *AnswerController) DeleteAnswer(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	switch res := cnt.service.Delete(id); res {
	case service.Ok:
		c.JSON(fiber.Map{
			"status":  "ok",
			"message": "deleted",
		})
		return c.SendStatus(fiber.StatusOK)
	case service.InternalError:
		return c.SendStatus(fiber.StatusInternalServerError)
	case service.BadRequest:
		c.JSON(fiber.Map{
			"status":  "error",
			"message": "there is no answer with this id",
		})
		return c.SendStatus(fiber.StatusBadRequest)
	default:
		return c.SendStatus(fiber.StatusNoContent)
	}

	return nil
}

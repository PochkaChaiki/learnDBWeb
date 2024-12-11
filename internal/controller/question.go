package controller

import (
	"learnDB/internal/domain"
	"learnDB/internal/service"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

// GET /api/question
// GET /api/question/{id}
// POST /api/question
// PUT /api/question
// DELETE /api/question/{id}

type QuestionController struct {
	service *service.ServiceQuestion
}

func NewQuestionController(s *service.ServiceQuestion) *QuestionController {
	return &QuestionController{service: s}
}

func (cnt *QuestionController) CreateQuestion(c fiber.Ctx) error {
	q := new(domain.Question)
	if err := c.Bind().JSON(q); err != nil {
		log.Fatalf("create question error: %s", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	switch res := cnt.service.Create(q); res {
	case service.Ok:
		c.JSON(fiber.Map{
			"status":  "ok",
			"message": "created",
		})
		return c.SendStatus(fiber.StatusCreated)
	case service.InternalError:
		c.SendStatus(fiber.StatusInternalServerError)
	default:
		return c.SendStatus(fiber.StatusNoContent)
	}

	return nil
}

func (cnt *QuestionController) GetQuestion(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		log.Fatalf("get query error: %s", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	switch q, res := cnt.service.Get(id); res {
	case service.Ok:
		if err := c.JSON(q); err != nil {
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

func (cnt *QuestionController) GetAllQuestions(c fiber.Ctx) error {
	switch qs, result := cnt.service.GetAll(); result {
	case service.Ok:
		if err := c.JSON(qs); err != nil {
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

func (cnt *QuestionController) UpdateQuestion(c fiber.Ctx) error {
	q := new(domain.Question)
	if err := c.Bind().JSON(q); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	switch res := cnt.service.Update(q); res {
	case service.Ok:
		c.JSON(fiber.Map{
			"status":  "ok",
			"message": "updated",
		})
		return c.SendStatus(fiber.StatusOK)
	case service.InternalError:
		return c.SendStatus(fiber.StatusInternalServerError)
	case service.BadRequest:
		c.JSON(fiber.Map{
			"status":  "error",
			"message": "there is no question with this id",
		})
		c.SendStatus(fiber.StatusBadRequest)
	default:
		return c.SendStatus(fiber.StatusNoContent)
	}

	return nil
}

func (cnt *QuestionController) DeleteQuestion(c fiber.Ctx) error {
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
			"message": "there is no question with this id",
		})
		return c.SendStatus(fiber.StatusBadRequest)
	default:
		return c.SendStatus(fiber.StatusNoContent)
	}

	return nil
}

package controller

import (
	"learnDB/internal/domain"
	"learnDB/internal/service"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

// GET /api/query
// GET /api/query/{id}
// POST /api/query

type QueryController struct {
	service *service.ServiceQuery
}

func NewQueryController(s *service.ServiceQuery) *QueryController {
	return &QueryController{service: s}
}

func (cnt *QueryController) GetAllQueries(c fiber.Ctx) error {
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

func (cnt *QueryController) GetQuery(c fiber.Ctx) error {
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

func (cnt *QueryController) CreateQuery(c fiber.Ctx) error {

	q := new(domain.Query)
	if err := c.Bind().JSON(q); err != nil {
		log.Fatalf("create query error: %s", err)
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

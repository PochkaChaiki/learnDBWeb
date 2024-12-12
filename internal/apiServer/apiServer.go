package apiserver

import (
	"learnDB/internal/controller"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

type APIServer struct {
	address    string
	controller *controller.APIController
	auth       *controller.AuthController
}

func New(a string, s []byte, c *controller.APIController, auth *controller.AuthController) *APIServer {
	return &APIServer{address: a, controller: c, auth: auth}
}

// GET /api/answer
// GET /api/answer/{id}
// POST /api/answer --> /api/query
// DELETE /api/answer/{id}

// GET /api/query
// GET /api/query/{id}
// POST /api/query

// GET /api/question
// GET /api/question/{id}
// POST /api/question
// PUT /api/question
// DELETE /api/question/{id}

func (s *APIServer) Run() {
	app := fiber.New()

	Logger := logger.New()
	app.Use(Logger)

	app.Post("/login", s.auth.Login)

	api := app.Group("/api")

	question := api.Group("/question")
	question.Get("/", s.controller.GetAllQuestions)
	question.Get("/:id", s.controller.GetQuestion)

	question.Use(s.auth.AuthMiddleware)

	question.Post("/", s.controller.CreateQuestion)
	question.Put("/", s.controller.UpdateQuestion)
	question.Delete("/:id", s.controller.DeleteQuestion)

	query := api.Group("/query")
	query.Use(s.auth.AuthMiddleware)
	query.Get("/", s.controller.GetAllQueries)
	query.Get("/:id", s.controller.GetQuery)
	query.Post("/", s.controller.CreateQuery)

	answer := api.Group("/answer")
	answer.Use(s.auth.AuthMiddleware)
	answer.Get("/", s.controller.GetAllAnswers)
	answer.Get("/:id", s.controller.GetAnswer)
	answer.Post("/", s.controller.CreateAnswer)
	answer.Delete("/:id", s.controller.DeleteAnswer)

	log.Fatal(app.Listen(s.address))
}

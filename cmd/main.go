package main

import (
	"learnDB/internal/config"
	"learnDB/internal/storage"

	"github.com/gofiber/fiber/v3"
)

func main() {

	config := config.MustLoad()
	storage := storage.New(config.StoragePath)

	app := fiber.New()

	api := app.Group("/api")

	question := api.Group("/question")
	question.Use(authMiddleware)
	question.Get("/", getQuestions)
	question.Get("/:id", getQuestion)
	question.Post("/", postQuestion)
	question.Put("/", putQuestion)
	question.Delete("/:id", deleteQuestion)

	answer := api.Group("/answer")
	answer.Use(authMiddleware)
	answer.Get("/", getAnswer)
	answer.Get("/:id", getAnswer)
	answer.Post("/", postAnswer)
	answer.Put("/", putAnswer)
	answer.Delete("/:id", deleteAnswer)
}

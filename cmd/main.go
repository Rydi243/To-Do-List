package main

import (
	"encoding/json"
	"log"

	"To-Do-List/internal/app"
	"To-Do-List/internal/contract"
	"To-Do-List/internal/pkg/db"

	"github.com/gofiber/fiber/v3"
	_ "github.com/lib/pq"
)

func main() {
	conn := db.NewDB()

	r := fiber.New()

	r.Post("/tasks", func(c fiber.Ctx) error {
		var s contract.Task
		if err := json.Unmarshal(c.Body(), &s); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Ошибка при парсинге post")
		}

		if err := app.PostHandler(conn, s); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Ошибка при добавлении задачи")
		}

		return c.SendString("Запись добавлена")
	})

	r.Get("/tasks", func(c fiber.Ctx) error {
		tasks, err := app.GetHandler(conn)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Ошибка при получении задач")
		}

		return c.JSON(tasks)
	})

	r.Put("/tasks/:id", func(c fiber.Ctx) error {
		var s contract.PutDelTask
		if err := json.Unmarshal(c.Body(), &s); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Ошибка при парсинге put")
		}

		if err := app.PutHandler(conn, s); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Ошибка при обновлении задачи")
		}

		return c.SendString("Запись обновлена")
	})

	r.Delete("/tasks/:id", func(c fiber.Ctx) error {
		var s contract.PutDelTask
		if err := json.Unmarshal(c.Body(), &s); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Ошибка при парсинге delete")
		}
		if err := app.DelHandler(conn, s); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Ошибка при удалении задачи")
		}
		return c.SendString("Запись удалена")
	})

	log.Fatal(r.Listen(":8080"))
}

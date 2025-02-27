package main

import (
	"context"
	"encoding/json"
	"log"

	"To-Do-List/internal/app"
	"To-Do-List/internal/contract"
	"To-Do-List/internal/pkg/db"

	"github.com/gofiber/fiber/v3"
	_ "github.com/lib/pq"
)

// @title To-Do List API
// @version 1.0
// @description API для управления задачами в To-Do List
// @host localhost:8080
// @BasePath /
func main() {
	conn := db.NewDB()
	defer conn.Close(context.Background())

	r := fiber.New()
	log.Println("Роутер fiber создан")

	// @Summary Добавить задачу
	// @Description Добавляет новую задачу в список
	// @Tags tasks
	// @Accept json
	// @Produce json
	// @Param task body contract.Task true "Данные задачи"
	// @Success 200 {string} string "Запись добавлена"
	// @Failure 400 {string} string "Ошибка при парсинге JSON"
	// @Failure 500 {string} string "Ошибка при добавлении задачи"
	// @Router /tasks [post]
	r.Post("/tasks", func(c fiber.Ctx) error {
		var s contract.Task
		if err := json.Unmarshal(c.Body(), &s); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Ошибка при парсинге JSON")
		}

		if err := app.PostHandler(conn, s); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Ошибка при добавлении задачи")
		}

		return c.SendString("Запись добавлена")
	})

	// @Summary Получить список задач
	// @Description Возвращает список всех задач
	// @Tags tasks
	// @Produce json
	// @Success 200 {array} contract.GetTask
	// @Failure 500 {string} string "Ошибка при получении задач"
	// @Router /tasks [get]
	r.Get("/tasks", func(c fiber.Ctx) error {
		tasks, err := app.GetHandler(conn)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Ошибка при получении задач")
		}

		return c.JSON(tasks)
	})

	// @Summary Обновить задачу
	// @Description Обновляет существующую задачу
	// @Tags tasks
	// @Accept json
	// @Produce json
	// @Param id path int true "ID задачи"
	// @Param task body contract.PutDelTask true "Данные задачи"
	// @Success 200 {string} string "Запись обновлена"
	// @Failure 400 {string} string "Ошибка при парсинге JSON"
	// @Failure 500 {string} string "Ошибка при обновлении задачи"
	// @Router /tasks/{id} [put]
	r.Put("/tasks/:id", func(c fiber.Ctx) error {
		var s contract.PutDelTask
		if err := json.Unmarshal(c.Body(), &s); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Ошибка при парсинге JSON")
		}

		if err := app.PutHandler(conn, s); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Ошибка при обновлении задачи")
		}

		return c.SendString("Запись обновлена")
	})

	// @Summary Удалить задачу
	// @Description Удаляет задачу по ID
	// @Tags tasks
	// @Accept json
	// @Produce json
	// @Param id path int true "ID задачи"
	// @Param task body contract.PutDelTask true "Данные задачи"
	// @Success 200 {string} string "Запись удалена"
	// @Failure 400 {string} string "Ошибка при парсинге JSON"
	// @Failure 500 {string} string "Ошибка при удалении задачи"
	// @Router /tasks/{id} [delete]
	r.Delete("/tasks/:id", func(c fiber.Ctx) error {
		var s contract.PutDelTask
		if err := json.Unmarshal(c.Body(), &s); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Ошибка при парсинге JSON")
		}

		if err := app.DelHandler(conn, s); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Ошибка при удалении задачи")
		}

		return c.SendString("Запись удалена")
	})

	log.Fatal(r.Listen(":8080"))
}

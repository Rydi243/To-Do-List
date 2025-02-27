package main

import (
	"database/sql"
	"encoding/json"
	"time"

	"log"

	"github.com/gofiber/fiber/v3"
	_ "github.com/lib/pq"
)

type Task struct {
	Title       string
	Description int
	Status      string
}

type GetTask struct {
	Id          int
	Title       string
	Description int
	Status      string
	Created_at  time.Time
	Updated_at  time.Time
}

type PutDelTask struct {
	Id          int
	Title       string
	Description int
	Status      string
}

func postHandler(db *sql.DB, t Task) error {
	query := "INSERT INTO tasks (title, description, status) VALUES ($1, $2, $3)"
	_, err := db.Exec(query, t.Title, t.Description, t.Status)
	return err
}
func getHandler(db *sql.DB) ([]GetTask, error) {
	query := "SELECT id, title, description, status, created_at, updated_at FROM tasks"
	row, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	var tasks []GetTask

	for row.Next() {
		var t GetTask
		err := row.Scan(&t.Id, &t.Title, &t.Description, &t.Status, &t.Created_at, &t.Updated_at)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}
func putHandler(db *sql.DB, i PutDelTask) error {
	query := "UPDATE tasks SET title = $1, description = $2, status = $3, updated_at = NOW() WHERE id = $4"
	_, err := db.Exec(query, i.Title, i.Description, i.Status, i.Id)
	return err
}
func delHandler(db *sql.DB, i PutDelTask) error {
	query := "DELETE FROM tasks WHERE id = $1"
	_, err := db.Exec(query, i.Id)
	return err
}

func main() {
	dsn := "user=postgres dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := fiber.New()

	r.Post("/tasks", func(c fiber.Ctx) error {
		var s Task
		if err := json.Unmarshal(c.Body(), &s); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Ошибка при парсинге post")
		}

		if err := postHandler(db, s); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Ошибка при добавлении задачи")
		}

		return c.SendString("Запись добавлена")
	})

	r.Get("/tasks", func(c fiber.Ctx) error {
		tasks, err := getHandler(db)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Ошибка при получении задач")
		}

		return c.JSON(tasks)
	})

	r.Put("/tasks/:id", func(c fiber.Ctx) error {
		var s PutDelTask
		if err := json.Unmarshal(c.Body(), &s); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Ошибка при парсинге put")
		}

		if err := putHandler(db, s); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Ошибка при обновлении задачи")
		}

		return c.SendString("Запись обновлена")
	})

	r.Delete("/tasks/:id", func(c fiber.Ctx) error {
		var s PutDelTask
		if err := json.Unmarshal(c.Body(), &s); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Ошибка при парсинге delete")
		}
		if err := delHandler(db, s); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Ошибка при удалении задачи")
		}
		return c.SendString("Запись удалена")
	})

	log.Fatal(r.Listen(":8080"))
}

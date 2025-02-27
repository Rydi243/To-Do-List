package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5"

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

func postHandler(conn *pgx.Conn, t Task) error {
	query := "INSERT INTO tasks (title, description, status) VALUES ($1, $2, $3)"
	_, err := conn.Exec(context.Background(), query, t.Title, t.Description, t.Status)
	return err
}
func getHandler(conn *pgx.Conn) ([]GetTask, error) {
	query := "SELECT id, title, description, status, created_at, updated_at FROM tasks"

	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []GetTask

	for rows.Next() {
		var t GetTask
		err := rows.Scan(&t.Id, &t.Title, &t.Description, &t.Status, &t.Created_at, &t.Updated_at)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}
func putHandler(conn *pgx.Conn, i PutDelTask) error {
	query := "UPDATE tasks SET title = $1, description = $2, status = $3, updated_at = NOW() WHERE id = $4"

	_, err := conn.Exec(context.Background(), query, i.Title, i.Description, i.Status, i.Id)

	return err
}
func delHandler(conn *pgx.Conn, i PutDelTask) error {

	_, err := conn.Exec(context.Background(), "DELETE FROM tasks WHERE id = $1", i.Id)

	return err
}

func main() {
	dsn := "postgres://postgres:postgres@localhost:5432/postgres"

	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Unable to connection to database: %v\n", err)
	}
	defer conn.Close(context.Background())
	log.Print("Connected db!")

	r := fiber.New()

	r.Post("/tasks", func(c fiber.Ctx) error {
		var s Task
		if err := json.Unmarshal(c.Body(), &s); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Ошибка при парсинге post")
		}

		if err := postHandler(conn, s); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Ошибка при добавлении задачи")
		}

		return c.SendString("Запись добавлена")
	})

	r.Get("/tasks", func(c fiber.Ctx) error {
		tasks, err := getHandler(conn)
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

		if err := putHandler(conn, s); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Ошибка при обновлении задачи")
		}

		return c.SendString("Запись обновлена")
	})

	r.Delete("/tasks/:id", func(c fiber.Ctx) error {
		var s PutDelTask
		if err := json.Unmarshal(c.Body(), &s); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Ошибка при парсинге delete")
		}
		if err := delHandler(conn, s); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Ошибка при удалении задачи")
		}
		return c.SendString("Запись удалена")
	})

	log.Fatal(r.Listen(":8080"))
}

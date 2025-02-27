package app

import (
	"context"

	"github.com/jackc/pgx/v5"

	_ "github.com/lib/pq"

	"To-Do-List/internal/contract"
)

func PostHandler(conn *pgx.Conn, t contract.Task) error {
	query := "INSERT INTO tasks (title, description, status) VALUES ($1, $2, $3)"
	_, err := conn.Exec(context.Background(), query, t.Title, t.Description, t.Status)
	return err
}

func GetHandler(conn *pgx.Conn) ([]contract.GetTask, error) {
	query := "SELECT id, title, description, status, created_at, updated_at FROM tasks"

	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []contract.GetTask

	for rows.Next() {
		var t contract.GetTask
		err := rows.Scan(&t.Id, &t.Title, &t.Description, &t.Status, &t.Created_at, &t.Updated_at)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}

func PutHandler(conn *pgx.Conn, i contract.PutDelTask) error {
	query := "UPDATE tasks SET title = $1, description = $2, status = $3, updated_at = NOW() WHERE id = $4"

	_, err := conn.Exec(context.Background(), query, i.Title, i.Description, i.Status, i.Id)

	return err
}

func DelHandler(conn *pgx.Conn, i contract.PutDelTask) error {

	_, err := conn.Exec(context.Background(), "DELETE FROM tasks WHERE id = $1", i.Id)

	return err
}

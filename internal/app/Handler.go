package app

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"

	_ "github.com/lib/pq"

	"To-Do-List/internal/contract"
)

// PostHandler обрабатывает запрос на добавление задачи.
// @Summary Добавить задачу
// @Description Добавляет новую задачу в список
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body contract.Task true "Данные задачи"
// @Success 200 {string} string "Запись добавлена"
// @Failure 500 {string} string "Ошибка при добавлении задачи"
// @Router /tasks [post]
func PostHandler(conn *pgx.Conn, t contract.Task) error {
	query := "INSERT INTO tasks (title, description, status) VALUES ($1, $2, $3)"
	_, err := conn.Exec(context.Background(), query, t.Title, t.Description, t.Status)
	log.Printf("Добавлена запись в таблицу tasks")
	return err
}

// GetHandler обрабатывает запрос на получение списка задач.
// @Summary Получить список задач
// @Description Возвращает список всех задач
// @Tags tasks
// @Produce json
// @Success 200 {array} contract.GetTask
// @Failure 500 {string} string "Ошибка при получении задач"
// @Router /tasks [get]
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

	log.Println("Выполнен запрос таблицы tasks")
	return tasks, nil
}

// PutHandler обрабатывает запрос на обновление задачи.
// @Summary Обновить задачу
// @Description Обновляет существующую задачу
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body contract.PutDelTask true "Данные задачи"
// @Success 200 {string} string "Запись обновлена"
// @Failure 500 {string} string "Ошибка при обновлении задачи"
// @Router /tasks/{id} [put]
func PutHandler(conn *pgx.Conn, i contract.PutDelTask) error {
	query := "UPDATE tasks SET title = $1, description = $2, status = $3, updated_at = NOW() WHERE id = $4"

	_, err := conn.Exec(context.Background(), query, i.Title, i.Description, i.Status, i.Id)

	log.Println("Обновлена информация в таблице tasks")
	return err
}

// DelHandler обрабатывает запрос на удаление задачи.
// @Summary Удалить задачу
// @Description Удаляет задачу по ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body contract.PutDelTask true "Данные задачи"
// @Success 200 {string} string "Запись удалена"
// @Failure 500 {string} string "Ошибка при удалении задачи"
// @Router /tasks/{id} [delete]
func DelHandler(conn *pgx.Conn, i contract.PutDelTask) error {

	_, err := conn.Exec(context.Background(), "DELETE FROM tasks WHERE id = $1", i.Id)

	log.Println("Удалена информация из таблицы tasks")
	return err
}
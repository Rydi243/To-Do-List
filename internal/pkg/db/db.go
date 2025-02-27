package db

import (
	"context"
	"log"

	"To-Do-List/internal/config"

	"github.com/jackc/pgx/v5"
)

func NewDB() *pgx.Conn {

	conn, err := pgx.Connect(context.Background(), config.PgDSN)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v\n", err)
	}
	// defer conn.Close(context.Background())
	log.Print("Подключение к PostgreSQL успешно на порту 5432")

	return conn
}

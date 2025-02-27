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
		log.Fatalf("Unable to connection to database: %v\n", err)
	}
	defer conn.Close(context.Background())
	log.Print("Connected db!")

	return conn
}

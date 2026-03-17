package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

func ConnectDB() *pgx.Conn {
	connection, err := pgx.Connect(context.Background(), "postgres://anshal:strongpassword@192.168.1.12:5432/passwordStorage")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database connected")
	return connection
}

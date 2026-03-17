package main

import (
	"context"

	"github.com/anshal1/passwordStorage/src/db"
)

func main() {
	connection := db.ConnectDB()
	defer connection.Close(context.Background())
}

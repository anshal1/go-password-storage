package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func ConnectDB(DBUri string) *sql.DB {
	connection, err := sql.Open("postgres", DBUri)
	if err != nil {
		log.Fatal(err)
	}
	err = connection.Ping()
	if err != nil {
		fmt.Println(err)
		return connection
	}
	return connection
}

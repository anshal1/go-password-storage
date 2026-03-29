package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	mp "github.com/anshal1/migrations-package/utils"
	"github.com/anshal1/passwordStorage/src/db"
	userRepo "github.com/anshal1/passwordStorage/src/repo/user"
	userService "github.com/anshal1/passwordStorage/src/services/user"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load() // loads ".env" by default
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	config := mp.GetConfig()
	dburi := os.Getenv("DB")
	err = mp.CreateMigrations(config, dburi)
	if err != nil {
		log.Fatal(err)
	}
	db := db.ConnectDB(dburi)
	if db == nil {
		fmt.Println("Unable to connect to db")
		return
	}
	newUserRepo := userRepo.NewUserRepo(db)
	newUserService := userService.NewUserService(newUserRepo)

	http.HandleFunc("/user", newUserService.UserHandler)
	http.HandleFunc("/user/login", newUserService.HandleLogin)
	err = http.ListenAndServe(":9999", nil)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
}

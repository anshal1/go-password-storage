package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	mp "github.com/anshal1/migrations-package/utils"
	"github.com/anshal1/passwordStorage/src/db"
	passwordRepo "github.com/anshal1/passwordStorage/src/repo/password"
	userRepo "github.com/anshal1/passwordStorage/src/repo/user"
	passwordService "github.com/anshal1/passwordStorage/src/services/password"
	userService "github.com/anshal1/passwordStorage/src/services/user"
	"github.com/anshal1/passwordStorage/utils"
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

	newPasswordRepo := passwordRepo.NewPasswordRepo(db)
	newPasswordService := passwordService.NewPasswordService(newPasswordRepo)

	http.HandleFunc("/user", utils.Log(newUserService.UserHandler))
	http.HandleFunc("/user/login", utils.Log(newUserService.HandleLogin))
	http.HandleFunc("/save-password", utils.Log(newPasswordService.SavePasswordHandler))
	http.HandleFunc("/get-password", utils.Log(newPasswordService.GetPasswordHandler))
	err = http.ListenAndServe(":9999", nil)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
}

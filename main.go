package main

import (
	"fmt"
	"log"
	"net/http"

	mp "github.com/anshal1/migrations-package/utils"
	"github.com/anshal1/passwordStorage/src/db"
	userModel "github.com/anshal1/passwordStorage/src/models/user"
	userRepo "github.com/anshal1/passwordStorage/src/repo/user"
	userService "github.com/anshal1/passwordStorage/src/services/user"
)

type Temp struct {
	Name string
}

func (t *Temp) AddUser(user userModel.User) error {
	return nil
}

func (t *Temp) UpdateUser(user userModel.User) error {
	return nil
}

func main() {
	config := mp.GetConfig()
	dburi := "postgres://anshal:strongpassword@192.168.1.12:5432/passwordStorage"
	err := mp.CreateMigrations(config, dburi)
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
	err = http.ListenAndServe(":9999", nil)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
}

package userRepo

import (
	"database/sql"
	"errors"
	"log"

	userModel "github.com/anshal1/passwordStorage/src/models/user"
	"golang.org/x/crypto/bcrypt"
)

type UserRepo struct {
	DB *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{
		DB: db,
	}
}

// use this function in user service
func (u *UserRepo) AddUser(user userModel.User) error {
	// add user logic goes here
	var exists bool
	err := u.DB.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`,
		user.Username,
	).Scan(&exists)
	if err != nil {
		log.Printf("AddUser: existence check failed: %v", err)
		return err
	}
	if exists {
		return errors.New("Username already taken")
	}

	// 5. Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("AddUser: bcrypt failed: %v", err)
		return err
	}

	// 6. Insert user
	_, err = u.DB.Exec(
		`INSERT INTO users (username, password) VALUES ($1, $2)`,
		user.Username, string(hash),
	)
	if err != nil {
		log.Printf("AddUser: insert failed: %v", err)
		return errors.New("Unable to create user")
	}
	return nil
}

func (u *UserRepo) UpdateUser(user userModel.User) error {
	return nil
}

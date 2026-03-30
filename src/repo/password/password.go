package passwordRepo

import (
	"database/sql"
	"errors"
	"fmt"

	passwordsModel "github.com/anshal1/passwordStorage/src/models/passwords"
	userRepo "github.com/anshal1/passwordStorage/src/repo/user"
	"github.com/anshal1/passwordStorage/utils"
)

type PasswordRepo struct {
	DB *sql.DB
}

func NewPasswordRepo(db *sql.DB) *PasswordRepo {
	return &PasswordRepo{
		DB: db,
	}
}

func (p *PasswordRepo) SavePassword(password passwordsModel.Password, jwtToken string) error {
	user, err := userRepo.GetCurrentUser(p.DB, jwtToken)
	if err != nil {
		return err
	}
	var domain string

	err = p.DB.QueryRow("select domain from passwords where domain = $1 and userId = $2", password.Domain, user.Id).Scan(&domain)
	if !errors.Is(err, sql.ErrNoRows) && err != nil {

		return err
	}
	if password.Secret == "" {
		return errors.New("password secret not found")
	}
	if domain != "" {
		return errors.New("domain already exists")
	}
	passwordHash, err := utils.GeneratePasswordHash(password.Password, password.Secret)
	if err != nil {
		fmt.Println(err)
		return err
	}
	_, err = p.DB.Exec("insert into passwords (password, domain, userId) values ($1, $2, $3)", passwordHash, password.Domain, user.Id)
	if err != nil {
		return errors.New("unable to save password")
	}
	return nil
}

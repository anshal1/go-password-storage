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
func (p *PasswordRepo) GetPassword(domain string, jwtToken string, secret string) (string, error) {
	user, err := userRepo.GetCurrentUser(p.DB, jwtToken)
	if err != nil {
		return "", err
	}
	var passwordHash string
	err = p.DB.QueryRow("select password from passwords where domain = $1 and userId = $2", domain, user.Id).Scan(&passwordHash)
	if errors.Is(err, sql.ErrNoRows) {
		return "", errors.New("password not found")
	}
	if err != nil {
		return "", err
	}
	password, apiErr := utils.VerifyPasswordAndGetPassword(passwordHash, []byte(secret))
	if apiErr != nil {
		return "", errors.New(apiErr.Message)
	}
	return password, nil
}

func (p *PasswordRepo) GetAllPasswords(page int, limit int, jwtToken string) ([]passwordsModel.AllPasswordsResponse, error) {
	var passwords []passwordsModel.AllPasswordsResponse
	user, err := userRepo.GetCurrentUser(p.DB, jwtToken)
	if err != nil {
		return nil, err
	}
	rows, err := p.DB.Query("select domain, id, password from passwords where userId = $1 limit $2 offset $3", user.Id, limit, (page-1)*limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var domain, password string
		var id int64
		if err := rows.Scan(&domain, &id, &password); err != nil {
			return nil, err
		}
		passwords = append(passwords, passwordsModel.AllPasswordsResponse{
			Domain:   domain,
			Id:       id,
			Password: password,
		})
	}
	return passwords, nil
}

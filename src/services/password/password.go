package passwordService

import (
	"encoding/json"
	"fmt"
	"net/http"

	passwordsModel "github.com/anshal1/passwordStorage/src/models/passwords"
	"github.com/anshal1/passwordStorage/utils"
)

type PasswordServiceContract interface {
	SavePassword(password passwordsModel.Password, jwtToken string) error
}

type PasswordService struct {
	passwordRepo PasswordServiceContract
}

func NewPasswordService(repo PasswordServiceContract) *PasswordService {
	return &PasswordService{
		passwordRepo: repo,
	}
}

func (p *PasswordService) SavePasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, passwordsModel.ErrMethodNotAllowed)
		return
	}
	var password passwordsModel.Password
	if err := json.NewDecoder(r.Body).Decode(&password); err != nil {
		utils.WriteError(w, passwordsModel.ErrInvalidJSON)
		return
	}
	defer r.Body.Close()

	if password.Domain == "" || password.Password == "" || password.Secret == "" {
		utils.WriteError(w, passwordsModel.ErrMissingPasswordFields)
		return
	}
	cookie, err := r.Cookie("access_token")
	if err != nil {
		utils.WriteError(w, &utils.APIError{Message: utils.UserNotFound, Code: 404})
		return
	}
	fmt.Println("Able to run")
	err = p.passwordRepo.SavePassword(password, cookie.Value)
	if err != nil {

		utils.WriteError(w, &utils.APIError{Message: err.Error(), Code: 500})
		return
	}
	utils.WriteJSON(w, 201, map[string]any{"message": "password saved successfully"})
}

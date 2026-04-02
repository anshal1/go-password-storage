package passwordService

import (
	"encoding/json"
	"net/http"
	"strconv"

	passwordsModel "github.com/anshal1/passwordStorage/src/models/passwords"
	"github.com/anshal1/passwordStorage/utils"
)

type PasswordServiceContract interface {
	SavePassword(password passwordsModel.Password, jwtToken string) error
	GetPassword(domain string, jwtToken string, secret string) (string, error)
	GetAllPasswords(page int, limit int, jwtToken string) ([]passwordsModel.AllPasswordsResponse, error)
	DeletePassword(id int64, jwtToken string) error
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
	token, tokenErr := utils.GetTokenFromHeader(r)
	if tokenErr != nil {
		utils.WriteError(w, tokenErr)
		return
	}
	err := p.passwordRepo.SavePassword(password, token)
	if err != nil {

		utils.WriteError(w, &utils.APIError{Message: err.Error(), Code: 500})
		return
	}
	utils.WriteJSON(w, 201, map[string]any{"message": "password saved successfully"})
}

func (p *PasswordService) GetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, &utils.APIError{Message: "method not allowed", Code: 405})
		return
	}
	var password struct {
		Domain string `json:"domain"`
		Secret string `json:"secret"`
	}
	if err := json.NewDecoder(r.Body).Decode(&password); err != nil {
		utils.WriteError(w, passwordsModel.ErrInvalidJSON)
		return
	}
	defer r.Body.Close()

	if password.Domain == "" || password.Secret == "" {
		utils.WriteError(w, passwordsModel.ErrMissingPasswordFields)
		return
	}
	token, tokenErr := utils.GetTokenFromHeader(r)
	if tokenErr != nil {
		utils.WriteError(w, tokenErr)
		return
	}
	plainPassword, err := p.passwordRepo.GetPassword(password.Domain, token, password.Secret)
	if err != nil {
		utils.WriteError(w, &utils.APIError{Message: err.Error(), Code: 500})
		return
	}
	utils.WriteJSON(w, 200, map[string]any{"message": "password retrieved successfully", "password": plainPassword})
}

func (p *PasswordService) GetAllPasswordsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, &utils.APIError{Message: "method not allowed", Code: 405})
		return
	}
	token, tokenErr := utils.GetTokenFromHeader(r)
	if tokenErr != nil {
		utils.WriteError(w, tokenErr)
		return
	}
	page := r.URL.Query().Get("page")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		utils.WriteError(w, &utils.APIError{Message: "page number not provided or invalid", Code: 400})
		return
	}
	if pageInt == 0 {
		pageInt = 1
	}
	limit := 10
	passwords, err := p.passwordRepo.GetAllPasswords(pageInt, limit, token)
	if err != nil {
		utils.WriteError(w, &utils.APIError{Message: err.Error(), Code: 500})
		return
	}
	utils.WriteJSON(w, 200, passwords)
}

func (p *PasswordService) DeletePasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.WriteError(w, &utils.APIError{Message: "method not allowed", Code: 405})
		return
	}
	cookie, err := r.Cookie("access_token")
	if err != nil {
		utils.WriteError(w, &utils.APIError{Message: utils.UserNotFound, Code: 404})
		return
	}
	id := r.URL.Query().Get("id")
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		utils.WriteError(w, &utils.APIError{Message: "invalid password id", Code: 400})
		return
	}
	if err := p.passwordRepo.DeletePassword(idInt, cookie.Value); err != nil {
		utils.WriteError(w, &utils.APIError{Message: err.Error(), Code: 500})
		return
	}
	utils.WriteJSON(w, 200, map[string]any{"message": "password deleted successfully"})
}

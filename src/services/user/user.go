package userService

import (
	"encoding/json"
	"net/http"

	userModel "github.com/anshal1/passwordStorage/src/models/user"
	"github.com/anshal1/passwordStorage/utils"
)

type UserRepoInterface interface {
	AddUser(user userModel.User) error
	Login(user userModel.User) error
}

type UserService struct {
	userRepo UserRepoInterface
}

func NewUserService(u UserRepoInterface) *UserService {
	return &UserService{
		userRepo: u,
	}
}

func (u *UserService) UserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, userModel.ErrMethodNotAllowed)
		return
	}
	var user userModel.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.WriteError(w, userModel.ErrInvalidJSON)
		return
	}
	defer r.Body.Close()

	if user.Username == "" || user.Password == "" {
		utils.WriteError(w, userModel.ErrMissingFields)
		return
	}

	if len(user.Password) < 8 {
		utils.WriteError(w, userModel.ErrWeakPassword)
		return
	}
	if err := u.userRepo.AddUser(user); err != nil {
		utils.WriteError(w, &utils.APIError{Message: err.Error(), Code: 400})
		return
	}
	jwtToken, err := utils.GenerateJWT(user.Username)
	if err != nil {
		utils.WriteError(w, &utils.APIError{Message: err.Error(), Code: 500})
		return
	}
	utils.SetAuthCookie(w, jwtToken)
	utils.WriteJSON(w, http.StatusCreated, map[string]string{
		"message":  "user created successfully",
		"username": user.Username,
		"access_token": jwtToken,
	})
}

func (u *UserService) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.WriteError(w, userModel.ErrMethodNotAllowed)
		return
	}

	var creds userModel.User
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		utils.WriteError(w, userModel.ErrInvalidJSON)
		return
	}
	defer r.Body.Close()

	if creds.Username == "" || creds.Password == "" {
		utils.WriteError(w, userModel.ErrMissingFields)
		return
	}
	err := u.userRepo.Login(creds)
	if err != nil {
		utils.WriteError(w, &utils.APIError{Message: err.Error(), Code: 400})
		return
	}
	jwt, err := utils.GenerateJWT(creds.Username)
	if err != nil {
		utils.WriteError(w, &utils.APIError{Message: err.Error(), Code: 500})
		return
	}
	utils.SetAuthCookie(w, jwt)
	utils.WriteJSON(w, 200, map[string]any{
		"message": "Login successfull",
		"access_token": jwt,
	})
}

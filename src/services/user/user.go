package userService

import (
	"encoding/json"
	"net/http"

	userModel "github.com/anshal1/passwordStorage/src/models/user"
	"github.com/anshal1/passwordStorage/utils"
)

type UserRepoInterface interface {
	AddUser(user userModel.User) error
	UpdateUser(user userModel.User) error
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
		utils.WriteError(w, userModel.ErrInternalServer)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, map[string]string{
		"message":  "user created successfully",
		"username": user.Username,
	})
}

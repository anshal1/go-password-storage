package userModel

import (
	"net/http"

	"github.com/anshal1/passwordStorage/utils"
)

var (
	ErrMissingFields    = &utils.APIError{Message: "username and password are required", Code: http.StatusBadRequest}
	ErrUsernameTaken    = &utils.APIError{Message: "username already taken", Code: http.StatusConflict}
	ErrInternalServer   = &utils.APIError{Message: "internal server error", Code: http.StatusInternalServerError}
	ErrMethodNotAllowed = &utils.APIError{Message: "method not allowed", Code: http.StatusMethodNotAllowed}
	ErrInvalidJSON      = &utils.APIError{Message: "invalid JSON body", Code: http.StatusBadRequest}
	ErrWeakPassword     = &utils.APIError{Message: "password must be at least 8 characters", Code: http.StatusBadRequest}
	ErrUserNotFound     = &utils.APIError{Message: "user not found", Code: http.StatusNotFound}
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

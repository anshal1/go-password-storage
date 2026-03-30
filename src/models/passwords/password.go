package passwordsModel

import (
	"net/http"

	"github.com/anshal1/passwordStorage/utils"
)

var (
	ErrMissingPasswordFields = &utils.APIError{Message: "password, domain and userId are required", Code: http.StatusBadRequest}
	ErrDomainTaken           = &utils.APIError{Message: "domain already exists", Code: http.StatusConflict}
	ErrPasswordNotFound      = &utils.APIError{Message: "password not found", Code: http.StatusNotFound}
	ErrInvalidUserID         = &utils.APIError{Message: "invalid userId", Code: http.StatusBadRequest}
	ErrMethodNotAllowed      = &utils.APIError{Message: "method not allowed", Code: http.StatusMethodNotAllowed}
	ErrInvalidJSON           = &utils.APIError{Message: "invalid JSON body", Code: http.StatusBadRequest}
	ErrInternalServer        = &utils.APIError{Message: "internal server error", Code: http.StatusInternalServerError}
)

type Password struct {
	Password string `json:"password"`
	Domain   string `json:"domain"`
	UserId   string `json:"userId"`
	Secret   string `json:"secret"`
}

package utils

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type APIError struct {
	Message string `json:"error"`
	Code    int    `json:"code"`
}

func (e *APIError) Error() string { return e.Message }

func WriteJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("writeJSON encode error: %v", err)
	}
}

func WriteError(w http.ResponseWriter, apiErr *APIError) {
	WriteJSON(w, apiErr.Code, apiErr)
}

func GenerateJWT(username string) (string, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))
	if len(secret) == 0 {
		return "", errors.New("JWT_SECRET env var is not set")
	}

	claims := Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   username,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func SetAuthCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,                    // not accessible via JS
		Secure:   true,                    // HTTPS only — set false for local dev
		SameSite: http.SameSiteStrictMode, // CSRF protection
		MaxAge:   int((24 * time.Hour).Seconds()),
	})
}

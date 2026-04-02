package utils

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}
type PasswordClaims struct {
	Password string `json:"password"`
	jwt.RegisteredClaims
}

type APIError struct {
	Message string `json:"error"`
	Code    int    `json:"code"`
}

var (
	LoginError   = "please login to continue"
	UserNotFound = "user not found"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	White  = "\033[37m"

	BoldRed    = "\033[1;31m"
	BoldGreen  = "\033[1;32m"
	BoldYellow = "\033[1;33m"
)

func (e *APIError) Error() string { return e.Message }

func WriteJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("status", strconv.Itoa(status))
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

func GeneratePasswordHash(password string, secret string) (string, error) {
	if len(secret) == 0 {
		return "", errors.New("JWT_SECRET env var is not set")
	}

	claims := PasswordClaims{
		Password: password,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:  password,
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
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

var (
	ErrInvalidToken  = &APIError{Message: "invalid or malformed token", Code: 401}
	ErrExpiredToken  = &APIError{Message: "token has expired", Code: 401}
	ErrMissingSecret = errors.New("JWT_SECRET env var is not set")
)

func VerifyJWT(tokenString string) (string, *APIError) {
	secret := []byte(os.Getenv("JWT_SECRET"))
	if len(secret) == 0 {
		return "", &APIError{Message: ErrMissingSecret.Error(), Code: 500}
	}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(t *jwt.Token) (any, error) {
			// Reject any token not signed with HMAC — prevents the "alg:none" attack
			if t.Method != jwt.SigningMethodHS256 {
				return "", ErrInvalidToken
			}
			return secret, nil
		},
	)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", ErrExpiredToken
		}
		return "", ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid || claims.Username == "" {
		return "", ErrInvalidToken
	}

	return claims.Username, nil
}

func VerifyPasswordAndGetPassword(passwordHash string, secret []byte) (string, *APIError) {

	if len(secret) == 0 {
		return "", &APIError{Message: ErrMissingSecret.Error(), Code: 500}
	}

	token, err := jwt.ParseWithClaims(
		passwordHash,
		&PasswordClaims{},
		func(t *jwt.Token) (any, error) {
			// Reject any token not signed with HMAC — prevents the "alg:none" attack
			if t.Method != jwt.SigningMethodHS256 {
				return "", ErrInvalidToken
			}
			return secret, nil
		},
	)
	if err != nil {
		return "", ErrInvalidToken
	}

	claims, ok := token.Claims.(*PasswordClaims)
	if !ok || !token.Valid || claims.Password == "" {
		return "", ErrInvalidToken
	}

	return claims.Password, nil
}

func PrintColoredLog(logText string, color string) {
	log.Print(color + logText + Reset)
}

func Log(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		fn(w, r)
		status := w.Header().Get("status")
		statusInt, err := strconv.Atoi(status)
		logText := path + " - Status: " + status
		if err != nil {
			PrintColoredLog("[ERROR] "+logText, Yellow)
			return
		}
		if statusInt < 400 {
			PrintColoredLog("[INFO] "+logText, Green)
			return
		}
		PrintColoredLog("[ERROR] "+logText, Red)

	}
}

func GetTokenFromHeader(r *http.Request) (string, *APIError) {
	headerValue := r.Header.Get("Authorization")
	if headerValue == "" {
		return "", &APIError{Message: "Authorization header missing", Code: 401}
	}
	if !strings.HasPrefix(headerValue, "Bearer") {
		return "", &APIError{Message: "Invalid authorization header", Code: 401}
	}
	return strings.Replace(headerValue, "Bearer ", "", 1), nil
}

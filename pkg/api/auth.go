// pkg/api/auth.go
package api

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type signInRequest struct {
	Password string `json:"password"`
}

type signInResponse struct {
	Token string `json:"token"`
}

type Claims struct {
	PasswordHash string `json:"ph"`
	jwt.RegisteredClaims
}

var (
	jwtKey = []byte(os.Getenv("TODOTODO_JWT_SECRET"))
)

func init() {
	if len(jwtKey) == 0 {
		jwtKey = []byte("my_secret_key")
		fmt.Println("WARNING: JWT secret not set – using default")
	}
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req signInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, "invalid json", http.StatusBadRequest)
		return
	}

	expected := os.Getenv("TODO_PASSWORD")
	if expected == "" {
		writeJSONError(w, "password not set", http.StatusInternalServerError)
		return
	}

	if req.Password != expected {
		writeJSONError(w, "Неверный пароль", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(8 * time.Hour)
	claims := &Claims{
		PasswordHash: hashPassword(expected),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		writeJSONError(w, "token error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expirationTime,
		Path:     "/",
		HttpOnly: true,
	})

	writeJSON(w, signInResponse{Token: tokenString})
}

func hashPassword(p string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(p)))
}
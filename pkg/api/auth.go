// pkg/api/auth.go
package api

import (
    "encoding/json"
    "net/http"
    "os"
    "time"
	"fmt"
	"crypto/md5"

    "github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("my_secret_key") // в реальности — из env

type Claims struct {
    PasswordHash string `json:"ph"`
    jwt.RegisteredClaims
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        writeJSONError(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var creds struct {
        Password string `json:"password"`
    }
    if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
        writeJSONError(w, "invalid json", http.StatusBadRequest)
        return
    }

    expected := os.Getenv("TODO_PASSWORD")
    if expected == "" {
        writeJSONError(w, "password not set", http.StatusInternalServerError)
        return
    }

    if creds.Password != expected {
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

    writeJSON(w, map[string]string{"token": tokenString})
}

func hashPassword(p string) string {
    return fmt.Sprintf("%x", md5.Sum([]byte(p)))
}
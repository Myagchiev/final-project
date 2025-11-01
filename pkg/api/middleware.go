// pkg/api/middleware.go
package api

import (
    "net/http"
    "os"

    "github.com/golang-jwt/jwt/v5"
)

func Auth(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        pass := os.Getenv("TODO_PASSWORD")
        if pass == "" {
            next(w, r)
            return
        }

        cookie, err := r.Cookie("token")
        if err != nil {
            http.Error(w, "Authentication required", http.StatusUnauthorized)
            return
        }

        tokenStr := cookie.Value
        claims := &Claims{}
        token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
            return jwtKey, nil
        })

        if err != nil || !token.Valid {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        if hashPassword(pass) != claims.PasswordHash {
            http.Error(w, "Password changed", http.StatusUnauthorized)
            return
        }

        next(w, r)
    }
}
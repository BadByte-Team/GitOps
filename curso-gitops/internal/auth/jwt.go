package auth

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Typed context key to avoid collisions
type contextKey string

const RoleKey contextKey = "role"
const UsernameKey contextKey = "username"

var jwtKey []byte

func init() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		fmt.Println("ADVERTENCIA: JWT_SECRET no definida, usando clave por defecto (NO usar en producción)")
		secret = "dev_default_key_change_me"
	}
	jwtKey = []byte(secret)
}

func GenerateJWT(username, role string) (string, error) {
	claims := &jwt.MapClaims{
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(8 * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtKey)
}

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, `{"error":"Token requerido"}`, http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method to prevent algorithm confusion attacks
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
			}
			return jwtKey, nil
		})

		if err != nil || token == nil || !token.Valid {
			http.Error(w, `{"error":"Token inválido"}`, http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, `{"error":"Token inválido"}`, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), RoleKey, claims["role"])
		ctx = context.WithValue(ctx, UsernameKey, claims["username"])
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Context().Value(RoleKey) != "admin" {
			http.Error(w, `{"error":"Acceso prohibido"}`, http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

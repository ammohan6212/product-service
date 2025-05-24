package main

import (
	"fmt"
	"net/http"
	"strings"
	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte("your-secret-key") // Must match user-service
var jwtAlgorithm = "HS256"

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	tokenString := ""
	authHeader := r.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		tokenString = strings.TrimPrefix(authHeader, "Bearer ")
	}
	if tokenString == "" {
		http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwtAlgorithm {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Welcome to the product service")
}

func main() {
	http.HandleFunc("/welcome", welcomeHandler)
	fmt.Println("Product service running on :8080")
	http.ListenAndServe(":8080", nil)
}
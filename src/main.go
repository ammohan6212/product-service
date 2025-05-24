package main

import (
	"fmt"
	"net/http"
)

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Welcome to the product service")
}

func main() {
	http.HandleFunc("/welcome", welcomeHandler)
	fmt.Println("Product service running on :8080")
	http.ListenAndServe(":8080", nil)
}

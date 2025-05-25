package main

import (
	"fmt"
	"log"
	"net/http"

	"yourapp/db"
)

func main() {
	db.Connect()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "✅ Go + Postgres app is running!")
	})

	log.Println("🚀 Server listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

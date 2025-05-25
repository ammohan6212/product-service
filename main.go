package main

import (
	"log"
	"yourapp/db"
)

func main() {
	db.Connect()
	log.Println("🚀 App is running...")
}

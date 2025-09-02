package main

import (
	"log"
	"net/http"
	"os"

	"github.com/BhadauriaAditya/LogLens/loglens"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ No .env file found, using system environment")
	}

	user := os.Getenv("ADMIN_USER")
	pass := os.Getenv("ADMIN_PASS")

	if user == "" || pass == "" {
		log.Fatal("ADMIN_USER and ADMIN_PASS must be set in .env or environment")
	}

	// Protected /logs route
	http.Handle("/logs", loglens.AuthMiddleware(http.HandlerFunc(loglens.ViewLogs)))

	// Start server
	log.Println("🚀 LogLens running at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"log"
	"net/http"
	"observe/database"
	"observe/handlers"
	"time"
)

func init() {
	db := database.GetDBConnection()
	database.CreateUsersTable(db)
	database.CreateProjectsTable(db)
	database.CreateLogsTable(db)
	database.CreateIndexes(db)
}

func main() {
	db := database.GetDBConnection()

	multiplexer := http.NewServeMux()
	multiplexer.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		handlers.UserRegistrationHandler(w, r, db)
	})
	multiplexer.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		handlers.UserAssertionHandler(w, r, db)
	})

	server := http.Server{
		Addr:         ":8080",
		Handler:      multiplexer,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("Server is listening on port 8080")
	log.Fatal(server.ListenAndServe())
}

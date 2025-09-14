package main

import (
	"log"
	"net/http"

	"eu.michalvalko.chess_arbiter_delegation_generator/internal/handlers"
)

func main() {
	// Serve static assets
	fs := http.FileServer(http.Dir("web/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	// Serve frontend
	http.HandleFunc("/", handlers.FrontendHandler)

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

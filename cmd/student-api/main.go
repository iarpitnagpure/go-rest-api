package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/iarpitnagpure/go-rest-api/internal/config"
)

func main() {
	// Load config
	cfg := config.MustLoad()

	// Database setup

	// Setup Router
	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to student api"))
	})

	// Setup Server
	server := http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	fmt.Println("Server Started")

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Fail to log server")
	}
}

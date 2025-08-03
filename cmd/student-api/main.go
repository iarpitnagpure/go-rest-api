package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/iarpitnagpure/go-rest-api/internal/config"
	"github.com/iarpitnagpure/go-rest-api/internal/http/handlers/students"
	"github.com/iarpitnagpure/go-rest-api/internal/storage/sqlite"
)

func main() {
	// Load config
	cfg := config.MustLoad()

	// Database setup
	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("storage initialized", storage)

	// Setup Router
	// Create new router to set REST APIs using http package NewServeMux method
	router := http.NewServeMux()

	// POST API to add new student
	router.HandleFunc("POST /api/students", students.NewStudent(storage))

	// Setup Server
	// Use Server method from http package and pass Address and router
	server := http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	fmt.Println("Server Started on", cfg.Address)

	// Make Gracefull exit from ongoing program,  Create chan with signal and perform use notify package to listen system notification
	done := make(chan os.Signal, 1)

	// Notify causes package signal to relay incoming signals to c.
	// If no signals are provided, all incoming signals will be relayed to c. Otherwise, just the provided signals will.
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Best practice to add in separate goroutine or add in main program
	go func() {
		// Listen to server using ListenAndServe method
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Fail to log server")
		}
	}()

	// Make Gracefull exit from ongoing program
	<-done

	slog.Info("Shutting down the request")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		slog.Error("Failed to shut down the server", slog.String("error", err.Error()))
	}

	slog.Info("Server shutdown susccessfully")
}

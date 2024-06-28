package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Todo struct {
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

func main() {

	todos := []Todo{}

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	r.Post("/todos", func(w http.ResponseWriter, r *http.Request) {
		todo := Todo{}
		err := json.NewDecoder(r.Body).Decode(&todo)

		r.Body = http.MaxBytesReader(w, r.Body, 1024)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		todos = append(todos, todo)

		w.WriteHeader(http.StatusCreated)
	})

	r.Get("/todos", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(todos)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	})

	srv := &http.Server{
		Addr:    ":4000",
		Handler: r,
	}

	// Create a channel to listen for OS signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		fmt.Println("Server is running on :4000")

		// This is just an unredable version of assigning the value to the "err" variable and checking if it's not nil
		if err := srv.ListenAndServe(); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}()

	fmt.Println("Press Ctrl+C to stop the server")

	// Wait for signals to gracefully shut down the server
	<-sigCh

	fmt.Println("Shutting down the server...")

	// Create a context with a timeout for the graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Println("Server gracefully stopped")

}

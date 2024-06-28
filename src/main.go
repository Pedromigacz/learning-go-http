package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Pedromigacz/learning-go-http/src/internal/handlers"
	"github.com/Pedromigacz/learning-go-http/src/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	todos := []store.Todo{}

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/healthcheck", handlers.NewHealthcheckHandler().ServerHTTP)

	r.Post("/todos", handlers.NewCreateTodoHandler(handlers.CreateTodoHandlerParams{
		Todos: &todos,
	}).ServerHTTP)

	r.Get("/todos", handlers.NewGetTodosHandler(handlers.GetTodosHandlerParams{
		Todos: &todos,
	}).ServerHTTP)

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

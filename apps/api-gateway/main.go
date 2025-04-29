package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

var redisClient *redis.Client

func main() {
	// Initialize Redis connection pool
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Create a new router
	r := mux.NewRouter()

	// Define routes
	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/health", healthHandler)

	// Create a new HTTP server
	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Implement graceful shutdown
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on :8080: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to the home page!")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		http.Error(w, "Redis is down", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, "OK")
}

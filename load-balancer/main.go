package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/patrickmn/go-cache"
	"golang.org/x/net/context"
)

var (
	redisClient *redis.Client
	lruCache    *cache.Cache
)

func main() {
	// Initialize Redis client
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Initialize LRU cache
	lruCache = cache.New(5*time.Minute, 10*time.Minute)

	// Create a new router
	r := mux.NewRouter()

	// Define routes
	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/health", healthHandler)

	// Set up TLS configuration
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	// Create a new HTTP server with TLS
	server := &http.Server{
		Addr:      ":443",
		Handler:   r,
		TLSConfig: tlsConfig,
	}

	// Start the server
	log.Println("Starting server on :443")
	log.Fatal(server.ListenAndServeTLS("server.crt", "server.key"))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	cacheKey := "home"

	// Check LRU cache
	if cachedResponse, found := lruCache.Get(cacheKey); found {
		fmt.Fprint(w, cachedResponse)
		return
	}

	// Check Redis cache
	cachedResponse, err := redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		lruCache.Set(cacheKey, cachedResponse, cache.DefaultExpiration)
		fmt.Fprint(w, cachedResponse)
		return
	}

	// Generate response
	response := "Welcome to the home page!"

	// Store response in Redis and LRU cache
	redisClient.Set(ctx, cacheKey, response, 0)
	lruCache.Set(cacheKey, response, cache.DefaultExpiration)

	fmt.Fprint(w, response)
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

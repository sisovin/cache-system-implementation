package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/patrickmn/go-cache"
)

type CacheMiddleware struct {
	redisClient *redis.Client
	lruCache    *cache.Cache
}

func NewCacheMiddleware(redisClient *redis.Client) *CacheMiddleware {
	return &CacheMiddleware{
		redisClient: redisClient,
		lruCache:    cache.New(5*time.Minute, 10*time.Minute),
	}
}

func (cm *CacheMiddleware) CacheKey(r *http.Request) string {
	return fmt.Sprintf("%s:%s", r.Method, r.URL.String())
}

func (cm *CacheMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	ctx := context.Background()
	cacheKey := cm.CacheKey(r)

	// Check LRU cache
	if cachedResponse, found := cm.lruCache.Get(cacheKey); found {
		fmt.Fprint(w, cachedResponse)
		return
	}

	// Check Redis cache
	cachedResponse, err := cm.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		cm.lruCache.Set(cacheKey, cachedResponse, cache.DefaultExpiration)
		fmt.Fprint(w, cachedResponse)
		return
	}

	// Capture the response
	recorder := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
	next(recorder, r)

	// Store response in Redis and LRU cache
	cm.redisClient.Set(ctx, cacheKey, recorder.body.String(), 0)
	cm.lruCache.Set(cacheKey, recorder.body.String(), cache.DefaultExpiration)
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (rr *responseRecorder) WriteHeader(statusCode int) {
	rr.statusCode = statusCode
	rr.ResponseWriter.WriteHeader(statusCode)
}

func (rr *responseRecorder) Write(b []byte) (int, error) {
	rr.body.Write(b)
	return rr.ResponseWriter.Write(b)
}

package ramcache

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

func BenchmarkCache(b *testing.B) {
	cache := NewCache(5 * time.Minute)
	ctx := context.Background()
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	b.Run("RAMCache_ConcurrentReadWrite", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				cache.Set("key", "value")
				cache.Get("key")
			}
		})
	})

	b.Run("RedisCache_ConcurrentReadWrite", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				redisClient.Set(ctx, "key", "value", 0)
				redisClient.Get(ctx, "key")
			}
		})
	})
}

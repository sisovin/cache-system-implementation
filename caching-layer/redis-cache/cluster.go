package rediscache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sony/gobreaker"
)

type Cluster struct {
	client         *redis.ClusterClient
	circuitBreaker *gobreaker.CircuitBreaker
}

func NewCluster(addrs []string) *Cluster {
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: addrs,
	})

	cbSettings := gobreaker.Settings{
		Name:        "RedisClusterCircuitBreaker",
		MaxRequests: 5,
		Interval:    60 * time.Second,
		Timeout:     30 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > 3
		},
	}

	return &Cluster{
		client:         client,
		circuitBreaker: gobreaker.NewCircuitBreaker(cbSettings),
	}
}

func (c *Cluster) Get(ctx context.Context, key string) (string, error) {
	result, err := c.circuitBreaker.Execute(func() (interface{}, error) {
		return c.client.Get(ctx, key).Result()
	})

	if err != nil {
		return "", err
	}

	return result.(string), nil
}

func (c *Cluster) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	_, err := c.circuitBreaker.Execute(func() (interface{}, error) {
		return nil, c.client.Set(ctx, key, value, expiration).Err()
	})

	return err
}

func (c *Cluster) MonitorHealth(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			_, err := c.client.Ping(ctx).Result()
			if err != nil {
				fmt.Println("Redis cluster health check failed:", err)
			} else {
				fmt.Println("Redis cluster is healthy")
			}
		case <-ctx.Done():
			return
		}
	}
}

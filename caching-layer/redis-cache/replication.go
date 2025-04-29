package rediscache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type Replication struct {
	client *redis.Client
	pubsub *redis.PubSub
}

func NewReplication(addr string) *Replication {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	pubsub := client.Subscribe(context.Background(), "cache-sync")

	return &Replication{
		client: client,
		pubsub: pubsub,
	}
}

func (r *Replication) PublishUpdate(ctx context.Context, key string, value string) error {
	err := r.client.Publish(ctx, "cache-sync", fmt.Sprintf("%s:%s", key, value)).Err()
	if err != nil {
		return fmt.Errorf("failed to publish update: %v", err)
	}
	return nil
}

func (r *Replication) SubscribeUpdates(ctx context.Context) {
	for {
		msg, err := r.pubsub.ReceiveMessage(ctx)
		if err != nil {
			fmt.Println("failed to receive message:", err)
			continue
		}

		// Process the message
		fmt.Println("received message:", msg.Payload)
	}
}

func (r *Replication) ResolveConflict(key string, value1 string, value2 string) string {
	// Implement CRDT conflict resolution logic here
	// For simplicity, we'll just return the lexicographically larger value
	if value1 > value2 {
		return value1
	}
	return value2
}

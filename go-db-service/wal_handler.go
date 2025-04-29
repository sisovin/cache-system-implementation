package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/segmentio/kafka-go"
)

type WALHandler struct {
	conn       *pgx.Conn
	kafkaWriter *kafka.Writer
}

func NewWALHandler(connString string, kafkaBrokers []string) (*WALHandler, error) {
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %v", err)
	}

	kafkaWriter := kafka.NewWriter(kafka.WriterConfig{
		Brokers: kafkaBrokers,
		Topic:   "cache-invalidation",
	})

	return &WALHandler{
		conn:       conn,
		kafkaWriter: kafkaWriter,
	}, nil
}

func (h *WALHandler) Close() {
	h.conn.Close(context.Background())
	h.kafkaWriter.Close()
}

func (h *WALHandler) ReadWALLogs() error {
	// Implement logic to read WAL logs and extract relevant changes
	// This is a placeholder implementation
	for {
		time.Sleep(1 * time.Second)
		changes := []string{"change1", "change2"} // Placeholder changes
		for _, change := range changes {
			if err := h.PublishCacheInvalidationEvent(change); err != nil {
				return fmt.Errorf("failed to publish cache invalidation event: %v", err)
			}
		}
	}
}

func (h *WALHandler) PublishCacheInvalidationEvent(change string) error {
	msg := kafka.Message{
		Value: []byte(change),
	}
	return h.kafkaWriter.WriteMessages(context.Background(), msg)
}

func main() {
	connString := "postgres://user:password@localhost:5432/mydb"
	kafkaBrokers := []string{"localhost:9092"}

	handler, err := NewWALHandler(connString, kafkaBrokers)
	if err != nil {
		log.Fatalf("Failed to create WAL handler: %v", err)
	}
	defer handler.Close()

	go func() {
		if err := handler.ReadWALLogs(); err != nil {
			log.Fatalf("Failed to read WAL logs: %v", err)
		}
	}()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Shutting down...")
}

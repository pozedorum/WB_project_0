package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/segmentio/kafka-go"
)

func main() {
	broker := os.Getenv("KAFKA_BROKER")
	topic := os.Getenv("KAFKA_TOPIC")
	groupID := os.Getenv("KAFKA_GROUP_ID")

	if broker == "" {
		broker = "localhost:9092"
	}
	if topic == "" {
		topic = "test-topic"
	}
	if groupID == "" {
		groupID = "test-group"
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   topic,
		GroupID: groupID,
	})
	defer reader.Close()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-sigchan:
			fmt.Println("Shutting down consumer")
			return
		default:
			msg, err := reader.ReadMessage(context.Background())
			if err != nil {
				fmt.Printf("Consumer error: %v\n", err)
				continue
			}
			fmt.Printf("Received: %s\n", string(msg.Value))
		}
	}
}

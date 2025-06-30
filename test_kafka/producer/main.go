package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	broker := os.Getenv("KAFKA_BROKER")
	topic := os.Getenv("KAFKA_TOPIC")

	if broker == "" {
		broker = "localhost:9092"
	}
	if topic == "" {
		topic = "test-topic"
	}

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{broker},
		Topic:   topic,
	})
	defer writer.Close()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-sigchan:
			fmt.Println("Shutting down producer")
			return
		default:
			msg := fmt.Sprintf("Hello World @ %v", time.Now())
			err := writer.WriteMessages(context.Background(),
				kafka.Message{
					Value: []byte(msg),
				},
			)
			if err != nil {
				fmt.Printf("Produce failed: %v\n", err)
			} else {
				fmt.Printf("Produced: %s\n", msg)
			}
			time.Sleep(3 * time.Second)
		}
	}
}

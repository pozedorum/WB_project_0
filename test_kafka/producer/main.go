package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IBM/sarama"
)

func main() {
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "kafka:9092"
	}

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	var producer sarama.SyncProducer
	var err error

	// Retry connection
	for i := 0; i < 10; i++ {
		producer, err = sarama.NewSyncProducer([]string{broker}, config)
		if err == nil {
			break
		}
		log.Printf("Connection attempt %d failed: %v\n", i+1, err)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		log.Fatalf("Failed to start producer: %v", err)
	}
	defer producer.Close()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	log.Println("Producer started. Press Ctrl+C to exit.")

	for {
		select {
		case <-ticker.C:
			msg := &sarama.ProducerMessage{
				Topic: "test-topic",
				Value: sarama.StringEncoder(fmt.Sprintf("Message at %s", time.Now().Format(time.RFC3339))),
			}

			partition, offset, err := producer.SendMessage(msg)
			if err != nil {
				log.Printf("Failed to send message: %v\n", err)
				continue
			}

			log.Printf("Sent message to partition %d at offset %d\n", partition, offset)

		case <-signals:
			log.Println("Shutting down producer...")
			return
		}
	}
}

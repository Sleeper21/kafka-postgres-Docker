package kafka

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func StartConsumer(topic string) {

	// Create a new consumer
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "my-group",
		"auto.offset.reset": "latest", // "earliest" to read from the beginning, "latest" to read only new messages
	})

	// Setting to "latest" will ignore all messages that were produced before the consumer started
	// it will only read messages that are produced after the consumer started

	if err != nil {
		panic(err)
	}

	// Subscribe to the topic
	err = consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("Kafka Consumer created and subscribed to topic:", topic)

	// Create channel to receive os signals ( enables us to gracefully shutdown the consumer)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start consuming messages
	// In this case the consumer will be running in a loop until we receive a signal
	run := true

	for run {
		select {
		case sig := <-sigChan:
			fmt.Println("Received signal:", sig)
			run = false
		default:
			ev := consumer.Poll(100) // Asks kafka for new messages and waits for 100ms
			switch e := ev.(type) {  // Switch over the event type that kafka has

			case *kafka.Message:
				// Handle the message - Print in this case
				fmt.Printf("Received message from Kafka: %s\n", string(e.Value))
				// Here you can add code to process the message, e.g., save it to the database

			case kafka.Error:
				// Handle errors - Print in this case
				fmt.Fprintf(os.Stderr, "Error: %v\n", e)
			}
		}
	}

	fmt.Println("Closing consumer")
	consumer.Close()
}

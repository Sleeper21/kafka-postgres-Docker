package kafka

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var Producer *kafka.Producer

func CreateProducer() {
	var err error
	Producer, err = kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092"})
	if err != nil {
		panic(err)
	}
	fmt.Println("Kafka Producer created")
}

func PublishMessage(topic string, message []byte) error {
	deliveryChan := make(chan kafka.Event)

	err := Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          message,
	}, deliveryChan)

	if err != nil {
		return fmt.Errorf("failed to produce message: %v", err)
	}

	e := <-deliveryChan
	msg := e.(*kafka.Message)

	if msg.TopicPartition.Error != nil {
		return msg.TopicPartition.Error
	}

	fmt.Printf("Message produced to topic %s [%d] at offset %v\n", *msg.TopicPartition.Topic, msg.TopicPartition.Partition, msg.TopicPartition.Offset)
	close(deliveryChan)
	return nil
}

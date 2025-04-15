package main

import "docker-kafka-postgres-kafkaUI/kafka"

func main() {
	kafka.StartConsumer("workers")
}

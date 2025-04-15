package main

import (
	"docker-kafka-postgres-kafkaUI/db"
	"docker-kafka-postgres-kafkaUI/kafka"
	"encoding/json"
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

type Worker struct {
	WorkerID uint   `gorm:"unique" json:"worker_id"`
	Name     string `json:"name"`
	Email    string `gorm:"unique" json:"email"`
	Roles    string `json:"roles"`
}

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to the database
	db.ConnectDB()

	// Start the Kafka producer
	kafka.CreateProducer()
	defer kafka.Producer.Close()

	// Perform database migrations
	err = db.DB.AutoMigrate(&Worker{})
	if err != nil {
		log.Fatal("Error migrating database: ", err)
	}

	// Create a new worker
	worker := Worker{
		WorkerID: 001,
		Name:     "Pedro Cruz",
		Email:    "email@example.com",
		Roles:    "frontend, backend",
	}

	result := db.DB.Create(&worker)

	if result.Error != nil {
		fmt.Println("Error creating worker: ", result.Error)
	} else {
		fmt.Println("Worker created successfully: ", worker)
	}

	// Add a new worker
	newWorker := Worker{
		WorkerID: 002,
		Name:     "John Doe",
		Email:    "new_one@example.com",
		Roles:    "Manager",
	}
	result = db.DB.Create(&newWorker)
	if result.Error != nil {
		fmt.Println("Error creating new worker: ", result.Error)
	} else {
		fmt.Println("New worker created successfully: ", newWorker)
	}

	// Query a worker
	var foundWorker Worker
	result = db.DB.First(&foundWorker, "worker_id = ?", 001)
	if result.Error != nil {
		fmt.Println("Error finding worker: ", result.Error)
	} else {
		fmt.Println("Worker found: ", foundWorker)
	}

	// Query all workers
	var workers []Worker
	result = db.DB.Find(&workers)
	if result.Error != nil {
		fmt.Println("Error finding workers: ", result.Error)
	} else {
		fmt.Println("Workers found: ", workers)
	}

	// Update a worker
	db.DB.Model(&Worker{}).Where("worker_id = ?", worker.WorkerID).Updates(Worker{
		Email: "newEmail@newemail.com",
	})

	// Update worker roles from worker_id
	workerID := 001
	roles := "fullstack"
	result = db.DB.Model(&Worker{}).Where("worker_id = ?", workerID).Update("roles", roles)
	if result.Error != nil {
		fmt.Println("Error updating worker roles: ", result.Error)
	} else {
		fmt.Println("Worker roles updated successfully")
	}

	// Delete all workers
	// allWorkers := workers
	// for _, worker := range allWorkers {
	// 	db.DB.Model(&Worker{}).Where("worker_id = ?", worker.WorkerID).Delete(&Worker{})
	// 	fmt.Printf("Worker %s deleted successfully\n", worker.Name)
	// }

	// Delete a worker by ID
	// workerIDToDelete := 001
	// db.DB.Model(&Worker{}).Where("worker_id = ?", workerIDToDelete).Delete(&Worker{})
	// fmt.Printf("Worker with ID %d deleted successfully\n", workerIDToDelete)

	// Publish workers updates to Kafka
	jsonData, _ := json.Marshal(workers)
	err = kafka.PublishMessage("workers", jsonData)
	if err != nil {
		fmt.Println("Error publishing message to Kafka: ", err)
	} else {
		fmt.Println("Message published to Kafka successfully")
	}
}

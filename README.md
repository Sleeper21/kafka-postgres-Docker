# Kafka - Postgres GORM - Emulator

This is a mini project made to interact with kafka as producer and as consumer using Docker. Also the Docker file contains a Postgres db image that can be interacted with. 

The intention of this is to interact with th DB using GORM , compare data coming from an external API with the db and sending new info received to a kakfa, and made a consumer subscribed to the kafka topic and read the event and send to another API 

## To run the Docker containers

-  Install Docker desktop and open it
- The docker-compose.yml script is using env variables to config the db connection. You can hardcode the port, user, password etc in the script or create a .env file.
- If using .env varibales run:
  ```bash
  docker-compose --env-file .env up -d
  ```
- If not using .env variables run:
  ```bash
  docker-compose up -d
  ```

### After running the docker-compose you will have

- Kafka service at localhost:9092
- Zookeeper at localhost:2181 (Kafka depends on this)
- PostgresSQL at your localhost env port
- Kafka UI at localhost:8080 ( Visual web interface )

### To run the Kafka consumer

If the NewConsumer(&kafka.ConfigMap{}) has the "auto.offset.reset": set to "latest", it means that the consumer will only read messages that are published after the consumer has started. 
So for that you need to:

- Open a second terminal
- On one terminal first run the consumer: ( from the project root )
```bash
  go run cmd/consumer_main.go
```
- After, on the second terminal run the main.go :
```bash
  go run main.go
```

If the "auto.offset.reset": set to "earliest" , you can run the main.go first because the consumer will read all messages available in kafka.






# Go Kafka Testing

This project is a Go application that demonstrates integration with Kafka and PostgreSQL. It includes an HTTP server that handles message creation, sending messages to Kafka, and storing them in PostgreSQL.

## Features

- Create and retrieve messages via HTTP endpoints.
- Send messages to a Kafka topic.
- Store messages in a PostgreSQL database.
- Retrieve statistics of messages.

## Starting the project

### Run with docker

```sh
docker-compose up -d --build
```

### Run manualy

**Clone the Repository**

    Open your terminal and clone the repository:
    ```sh
    git clone https://github.com/lavander40/go_kafka_testing
    cd yourproject
    ```

**Install Dependencies**

    Ensure all dependencies are installed. Go modules should manage this automatically if they are defined in `go.mod`:
    ```sh
    go mod tidy
    ```

**Run the Application**

    Use the `go run` command to start the application:
    ```sh
    go run cmd/main.go
    ```

# API Documentation

## API is avalable on the address

```url
http://localhost:8080
```

## Requesting with Curl 
```sh
curl -X POST http://localhost:8080/messages \
     -H 'Content-Type: application/json' \
     -d '{"content": "Hello, world!"}'
```

## Endpoints

### Create a Message

**Endpoint:** `POST /messages`

**Description:** Creates a new message.

**Request Body:**
```json
{
    "content": "string"
}
```
### Get all messages

**Endpoint:** `GET /messages`

**Description:** Gets all messages.

### Get messages stats

**Endpoint:** `GET /messages/stats`

**Description:** Gets statiscit about stored messages.
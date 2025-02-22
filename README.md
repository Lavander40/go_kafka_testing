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

    git clone http://github.com/Lavander40/go_kafka_testing
    cd go_kafka_testing

**Install Dependencies**

    go mod tidy

**Run the Application**

    go run cmd/main.go

# API Documentation

## API is avalable on the address

```url
https://lab.aksuma.ru/
```

## Requesting with Curl
### Linux:
```sh
curl -X POST http://lab.aksuma.ru/messages \
     -H 'Content-Type: application/json' \
     -d '{"content": "Hello, world!"}'
```
### Windows CMD:
```sh
    curl -X POST "http://lab.aksuma.ru/messages" -H "Content-Type: application/json" -d "{\"content\": \"Hello, world!\"}"
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

**Description:** Gets statistic about stored messages.
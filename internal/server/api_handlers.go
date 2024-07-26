package server

import (
	"encoding/json"
	"go_kafka_testing/internal/domain"
	"log/slog"
	"net/http"
	"time"
)

// getMessagesHandler handles GET requests to retrieve all messages from storage
func (s *Server) getMessagesHandler(w http.ResponseWriter, r *http.Request) {
	// getting messages from DB
	msgs, err := s.storage.GetMessages()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		s.log.Error("Error during messages reading:", slog.Any("error", err))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(msgs)
}

// createMessageHandler handles POST requests to create a new message
func (s *Server) createMessageHandler(w http.ResponseWriter, r *http.Request) {
	var msg domain.Message

	// Decode the request body into a Message struct
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if msg.Content == "" {
		http.Error(w, "content of the message is empty of field does not exist", http.StatusBadRequest)
		return
	}

	msg.CreatedAt = time.Now()

	// Save the message to the database
	id, err := s.storage.SaveMessage(&msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		s.log.Error("Error during message saving in db:", slog.Any("error", err))
		return
	}
	msg.ID = id

	// Send the message to Kafka
	if err = s.kafka.SendMessage(&msg); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		s.log.Error("Error during message sending to kafka:", slog.Any("error", err))
		return
	}
	s.log.Debug("message successfully send", slog.Any("message", msg))

	s.log.Debug("message successfully saved", slog.Any("message", msg))
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(msg)
}

// getStatsHandler handles GET requests to retrieve message statistics from storage
func (s *Server) getStatsHandler(w http.ResponseWriter, r *http.Request) {
	// getting stats from DB
	stats, err := s.storage.GetStats()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		s.log.Error("Error during stats reading:", slog.Any("error", err))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stats)
}
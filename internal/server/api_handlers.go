package server

import (
	"encoding/json"
	"go_kafka_testing/internal/domain"
	"log/slog"
	"net/http"
	"time"
)

func (s *Server) getMessagesHandler(w http.ResponseWriter, r *http.Request) {
	msgs, err := s.storage.GetMessages()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		s.log.Error("Error during messages reading:", slog.Any("error", err))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(msgs)
}

func (s *Server) createMessageHandler(w http.ResponseWriter, r *http.Request) {
	var msg domain.Message

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

	id, err := s.storage.SaveMessage(&msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		s.log.Error("Error during message saving in db:", slog.Any("error", err))
		return
	}
	msg.ID = id

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

func (s *Server) getStatsHandler(w http.ResponseWriter, r *http.Request) {
	stats, err := s.storage.GetStats()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		s.log.Error("Error during stats reading:", slog.Any("error", err))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stats)
}
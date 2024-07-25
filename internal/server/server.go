package server

import (
	"go_kafka_testing/internal/kafka"
	"go_kafka_testing/internal/storage"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	log     *slog.Logger
	port    string
	storage *storage.Storage
	kafka   kafka.KafkaWriter
	router  *mux.Router
}

func New(log *slog.Logger, port string, storage *storage.Storage, kafka kafka.KafkaWriter) *Server {
	return &Server{
		log:     log,
		port:    ":" + port,
		storage: storage,
		kafka:   kafka,
		router:  mux.NewRouter(),
	}
}

func (s *Server) Start() error {
	s.configureRouter()
	return http.ListenAndServe(s.port, s.router)
}

func (s *Server) configureRouter() {
	s.router.Use(s.middlewareFunc)	
	s.router.HandleFunc("/messages", s.getMessagesHandler).Methods("GET")
	s.router.HandleFunc("/messages", s.createMessageHandler).Methods("POST")
	s.router.HandleFunc("/messages/stats", s.getStatsHandler).Methods("GET")
}

func (s *Server) middlewareFunc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		s.log.Info("Got request from", slog.Any("method", r.Method), slog.Any("endpoint", r.RequestURI), slog.Any("ip", r.RemoteAddr))

		next.ServeHTTP(w, r)

		s.log.Debug("Completed request from", slog.Any("endpoint", r.RequestURI), slog.Any("ip", r.RemoteAddr), slog.Any("duration", time.Since(start)))
	})
}

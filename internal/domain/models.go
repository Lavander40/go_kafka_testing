// Domain package describes data structures used in different modules of the application
package domain

import (
	"database/sql"
	"time"
)

type Message struct {
	ID          int       `json:"id"`
	Content     string    `json:"content"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	ProcessedAt sql.NullTime `json:"processed_at,omitempty"`
}

type Stats struct {
	TotalMessages     int `json:"total_messages"`
	PendingMessages   int `json:"pending_messages"`
	ProcessedMessages int `json:"processed_messages"`
}
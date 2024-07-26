package storage

import "go_kafka_testing/internal/domain"

// SaveMessafe stores provided message to database
func (s *Storage) SaveMessage(msg *domain.Message) (id int, err error) {
	err = s.conn.QueryRow("INSERT INTO messages (content, status, created_at) VALUES ($1, 'pending', $2) RETURNING id;", msg.Content, msg.CreatedAt).Scan(&id)
	return
}

// GetMessages returns full list of stored messages
func (s *Storage) GetMessages() ([]*domain.Message, error) {
	var messages []*domain.Message

	rows, err := s.conn.Query("SELECT * FROM messages;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var message domain.Message
		err := rows.Scan(&message.ID, &message.Content, &message.Status, &message.CreatedAt, &message.ProcessedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, &message)
	}

	return messages, nil
}

// UpdateStatus applies changes to status of specific message in database 
func (s *Storage) UpdateStatus(msg *domain.Message) error {
	_, err := s.conn.Exec("UPDATE messages SET status=$1, processed_at=$2 WHERE id=$3;", msg.Status, msg.ProcessedAt, msg.ID)
	return err
}

// GetStats recive information about amount of pending, processed and overall amount for stored messages 
func (s *Storage) GetStats() (domain.Stats, error) {
	var stats domain.Stats
	
	err := s.conn.QueryRow("SELECT COUNT(*) FROM messages;").Scan(&stats.TotalMessages)
	if err != nil {
		return stats, err
	}
	
	err = s.conn.QueryRow("SELECT COUNT(*) FROM messages WHERE status='processed';").Scan(&stats.ProcessedMessages)
	if err != nil {
		return stats, err
	}

	stats.PendingMessages = stats.TotalMessages - stats.ProcessedMessages

	return stats, err
}

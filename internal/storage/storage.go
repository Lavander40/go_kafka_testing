package storage

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
)

type Storage struct {
	ctx  context.Context
	conn *sql.DB
	cfg  string
}

func New(ctx context.Context, connectString string) *Storage {
	return &Storage{
		ctx: ctx,
		cfg: connectString,
	}
}

func (s *Storage) Connect() error {
	conn, err := sql.Open("postgres", s.cfg)
	if err != nil {
		return err
	}

	if err = conn.Ping(); err != nil {
		return err
	}

	s.conn = conn
	return nil
}

func (s *Storage) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}

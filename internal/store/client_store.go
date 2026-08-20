package store

import "database/sql"

type ClientStore struct {
	db *sql.DB
}

func NewClientStore(db *sql.DB) *ClientStore {
	return &ClientStore{db: db}
}

// Ensure makes sure a client row exists for id, creating it if necessary.
func (s *ClientStore) Ensure(id string) error {
	_, err := s.db.Exec(`INSERT OR IGNORE INTO clients (id) VALUES (?)`, id)
	return err
}

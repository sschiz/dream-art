package shop

import "database/sql"

type Syncer struct {
	db *sql.DB
}

// Sync synchronizes the shop with the database
func (s *Syncer) Sync(*Shop) error {
	return nil
}

package shop

import "database/sql"

type Syncer struct {
	db *sql.DB
}

// Sync synchronizes the shop with the database
func (s *Syncer) Sync(shop *Shop) error {
	return nil
}

package db

import (
	"github.com/consumer-match-delete/internal/db/model"
)

// DeleteMatch marks a match as deleted in DB.
func(db *DB) DeleteMatch (m model.Match) error {
	_, err := db.stmtDeleteMatch.Exec(m.ID, m.DeletedAt)
	if err != nil {
		return err
	}

	return nil
}
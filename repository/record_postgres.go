
package repository

import (
	"fmt"

	recordsrestapi "github.com/Pinkman-77/records-restapi"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type RecordPostgres struct {
	db *sqlx.DB
}

func NewRecordPostgres(db *sqlx.DB) *RecordPostgres {
	return &RecordPostgres{db: db}
}

func (r *RecordPostgres) CreateRecord(record recordsrestapi.Record) (int, error) {
	// Retrieve the artist ID based on the artist's name
	var artistID int
	getArtistID := fmt.Sprintf("SELECT id FROM %s WHERE name = $1", artistTable)
	err := r.db.QueryRow(getArtistID, record.Artist).Scan(&artistID)
	if err != nil {
		return 0, fmt.Errorf("failed to get artist ID: %w", err)
	}

	// Insert the record using the artist ID
	tx, err := r.db.Beginx()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}

	var id int
	createRecord := fmt.Sprintf(`
		INSERT INTO %s (title, artist_id, year, tracklist, credits, duration)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id
	`, albumTable)

	err = tx.QueryRow(
		createRecord,
		record.Title,
		artistID,
		record.Year,
		pq.Array(record.Tracklist), // Convert slice to PostgreSQL array
		pq.Array(record.Credits),   // Convert slice to PostgreSQL array
		record.Duration,
	).Scan(&id)

	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to create record: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return id, nil
}

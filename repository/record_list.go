package repository

import (
    "database/sql"
    "fmt"

    recordsrestapi "github.com/Pinkman-77/records-restapi"
    "github.com/jmoiron/sqlx"
    "github.com/lib/pq" // Ensure you have the pq package for PostgreSQL
)

type RecordPostgres struct {
    db *sqlx.DB
}

func NewRecordPostgres(db *sqlx.DB) *RecordPostgres {
    return &RecordPostgres{db: db}
}

func (r *RecordPostgres) CreateRecord(record recordsrestapi.Record) (recordsrestapi.Record, error) {
	tx, err := r.db.Begin()
	if err != nil {
			return recordsrestapi.Record{}, fmt.Errorf("failed to start transaction: %w", err)
	}

	// 1. Get Artist ID 
	var artistID int
	err = r.db.QueryRow("SELECT id FROM artists WHERE name = $1", record.Artist).Scan(&artistID)
	if err != nil {
			tx.Rollback()
			if err == sql.ErrNoRows {
					return recordsrestapi.Record{}, fmt.Errorf("artist '%s' not found", record.Artist)
			}
			return recordsrestapi.Record{}, fmt.Errorf("failed to get artist ID: %w", err)
	}

	// 2. Insert Record 
	_, err = tx.Exec("INSERT INTO records (title, artist_id, year, tracklist, credits, duration) VALUES ($1, $2, $3, $4, $5, $6)", 
			record.Title, artistID, record.Year, pq.Array(record.Tracklist), pq.Array(record.Credits), record.Duration)
	if err != nil {
			tx.Rollback()
			return recordsrestapi.Record{}, fmt.Errorf("failed to create record: %w", err)
	}

	// 3. Commit the transaction
	if err := tx.Commit(); err != nil {
			return recordsrestapi.Record{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Return an empty record as we don't need to retrieve it immediately
	return recordsrestapi.Record{}, nil 
}
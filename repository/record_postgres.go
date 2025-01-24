
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

func (r *RecordPostgres) GetRecord(id int) (recordsrestapi.RecordWithArtist, error) {
 
        var record recordsrestapi.RecordWithArtist
    
        query := fmt.Sprintf(`
            SELECT r.id, r.title, a.name AS artist, r.year, r.tracklist, r.credits, r.duration
            FROM %s r
            INNER JOIN %s a ON r.artist_id = a.id
            WHERE r.id = $1
        `, albumTable, artistTable)
    
        err := r.db.QueryRow(
            query, id,
        ).Scan(
            &record.ID,
            &record.Title,
            &record.Artist, // Fetches artist name directly
            &record.Year,
            pq.Array(&record.Tracklist),
            pq.Array(&record.Credits),
            &record.Duration,
        )
    
        if err != nil {
            return record, err
        }
    
        return record, nil
    }
    func (r *RecordPostgres) GetAllRecords() ([]recordsrestapi.Record, error) {
        var records []recordsrestapi.Record
    
        query := fmt.Sprintf(`
            SELECT r.id, r.title, a.name AS artist, r.year, r.tracklist, r.credits, r.duration 
            FROM %s r
            INNER JOIN %s a ON r.artist_id = a.id
        `, albumTable, artistTable)
    
        rows, err := r.db.Query(query)
        if err != nil {
            return nil, err
        }
        defer rows.Close()
    
        for rows.Next() {
            var record recordsrestapi.Record
            err := rows.Scan(
                &record.ID,
                &record.Title,
                &record.Artist,
                &record.Year,
                pq.Array(&record.Tracklist),
                pq.Array(&record.Credits),
                &record.Duration,
            )
            if err != nil {
                return nil, err
            }
            records = append(records, record)
        }
    
        return records, nil
    }
    
   


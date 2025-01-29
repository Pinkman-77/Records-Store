package repository

import (
	"database/sql"

	"fmt"

	recordsrestapi "github.com/Pinkman-77/records-restapi"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type ArtistPostgres struct {
	db *sqlx.DB
}

func NewArtistPostgres(db *sqlx.DB) *ArtistPostgres {
	return &ArtistPostgres{db: db}
}

func (r *ArtistPostgres) CreateArtist(artist recordsrestapi.Artist) (int, error) {
    
	tx, err := r.db.Beginx()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}

	var id int

	err = tx.QueryRow("INSERT INTO artists (name) VALUES ($1) RETURNING id", artist.Name).Scan(&id)
	if err != nil {
		// Rollback transaction on error
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return 0, fmt.Errorf("failed to rollback transaction after error: %w", rollbackErr)
		}
		return 0, fmt.Errorf("failed to insert artist: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return id, nil
}


func (r *ArtistPostgres) GetAllArtists() ([]recordsrestapi.ArtistWithRecords, error) {
    var result []recordsrestapi.ArtistWithRecords // Explicitly initialize as an empty slice
    query := `
        SELECT a.id, a.name, r.id AS record_id, r.title, r.year, r.tracklist, r.credits, r.duration
        FROM artists a
        LEFT JOIN records r ON a.id = r.artist_id
    `

    rows, err := r.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    artistMap := make(map[uint]*recordsrestapi.ArtistWithRecords)

    for rows.Next() {
        var artistID uint
        var artistName string
        var recordID sql.NullString
        var recordTitle sql.NullString
        var recordYear sql.NullInt64
        var tracklist pq.StringArray
        var credits pq.StringArray
        var duration sql.NullString

        err := rows.Scan(&artistID, &artistName, &recordID, &recordTitle, &recordYear, &tracklist, &credits, &duration)
        if err != nil {
            return nil, err
        }

        if _, exists := artistMap[artistID]; !exists {
            artistMap[artistID] = &recordsrestapi.ArtistWithRecords{
                ID:      artistID,
                Name:    artistName,
                Records: []recordsrestapi.Record{},
            }
        }

        if recordID.Valid {
            record := recordsrestapi.Record{
                ID:        recordID.String,
                Title:     recordTitle.String,
                Artist:    artistName,
                Year:      recordYear.Int64,
                Tracklist: tracklist,
                Credits:   credits,
                Duration:  duration.String,
            }

            artistMap[artistID].Records = append(artistMap[artistID].Records, record)
        }
    }

    for _, artistWithRecords := range artistMap {
        result = append(result, *artistWithRecords)
    }

    return result, nil
}

func (r *ArtistPostgres) GetArtist(id int) (recordsrestapi.ArtistWithRecords, error) {
    var artist recordsrestapi.ArtistWithRecords
    artist.Records = []recordsrestapi.Record{} // Initialize empty slice to avoid null values

    query := `
        SELECT a.id, a.name, r.id AS record_id, r.title, r.year, r.tracklist, r.credits, r.duration 
        FROM artists a
        LEFT JOIN records r ON a.id = r.artist_id
        WHERE a.id = $1
    `

    rows, err := r.db.Queryx(query, id)
    if err != nil {
        return artist, err
    }
    defer rows.Close()

    for rows.Next() {
        var recordID sql.NullString
        var recordTitle sql.NullString
        var recordYear sql.NullInt64
        var tracklist pq.StringArray
        var credits pq.StringArray
        var duration sql.NullString

        err := rows.Scan(&artist.ID, &artist.Name, &recordID, &recordTitle, &recordYear, &tracklist, &credits, &duration)
        if err != nil {
            return artist, err
        }

        // Check if record exists, otherwise skip adding a record
        if recordID.Valid {
            record := recordsrestapi.Record{
                ID:        recordID.String,
                Title:     recordTitle.String,
                Artist:    artist.Name,
                Year:      recordYear.Int64,
                Tracklist: tracklist, // Properly scan tracklist array
                Credits:   credits,   // Properly scan credits array
                Duration:  duration.String,
            }
            artist.Records = append(artist.Records, record)
        }
    }

    return artist, nil
}


func (r *ArtistPostgres) UpdateArtist(id int, updatedArtist recordsrestapi.Artist) error {
    updateQuery := fmt.Sprintf("UPDATE %s SET name = $1 WHERE id = $2", artistTable)
    _, err := r.db.Exec(updateQuery, updatedArtist.Name, id)
    if err != nil {
            return fmt.Errorf("failed to update artist: %w", err)
    }
    return nil
}

func (r *ArtistPostgres) DeleteArtist(id int) error {
    deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE id = $1", artistTable)
    result, err := r.db.Exec(deleteQuery, id)
    if err != nil {
            return fmt.Errorf("failed to delete artist: %w", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
            return fmt.Errorf("failed to check rows affected: %w", err)
    }

    if rowsAffected == 0 {
            return fmt.Errorf("artist not found with ID: %d", id)
    }

    return nil
}

package repository

import (
	"database/sql"

	recordsrestapi "github.com/Pinkman-77/records-restapi"
	"github.com/jmoiron/sqlx"
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
		return 0, err
	}

	var id int

	err = tx.QueryRow("INSERT INTO artists (name) VALUES ($1) RETURNING id", artist.Name).Scan(&id)

	if err != nil {
		return 0, err
	}

	err = tx.Commit()

	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *ArtistPostgres) GetAllArtists() ([]recordsrestapi.ArtistWithRecords, error) {
    var result []recordsrestapi.ArtistWithRecords
    query := `
        SELECT a.id, a.name, r.id AS record_id, r.title, r.year 
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
        var recordYear sql.NullInt64 // Keep as sql.NullInt64

        err := rows.Scan(&artistID, &artistName, &recordID, &recordTitle, &recordYear)
        if err != nil {
            return nil, err
        }

        // Initialize artist entry if it doesn't exist
        if _, exists := artistMap[artistID]; !exists {
            artistMap[artistID] = &recordsrestapi.ArtistWithRecords{
                ID:      artistID,
                Name:    artistName,
                Records: []recordsrestapi.Record{},
            }
        }

        // Only add records if recordID is valid
        if recordID.Valid {
            record := recordsrestapi.Record{
                ID:     recordID.String,
                Title:  recordTitle.String,
                Artist: artistName,
                Year:   0, // Default value for Year
            }

            // Assign valid year if available
            if recordYear.Valid {
                record.Year = recordYear.Int64 // Use the valid year value
            }

            // Append the record to the artist's records
            artistMap[artistID].Records = append(artistMap[artistID].Records, record)
        }
    }

    // Convert map to slice for final result
    for _, artistWithRecords := range artistMap {
        result = append(result, *artistWithRecords)
    }

    return result, nil
}

func (r *ArtistPostgres) GetArtist(id int) (recordsrestapi.Artist, error) {
	var artist recordsrestapi.Artist
	err := r.db.Get(&artist, "SELECT id, name FROM artists WHERE id = $1", id)
	return artist, err
}
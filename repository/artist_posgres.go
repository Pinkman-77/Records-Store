package repository

import (
	recordsrestapi "github.com/Pinkman-77/records-restapi"
	"github.com/jmoiron/sqlx"
)

type ArtistPostgres struct {
	db *sqlx.DB
}

func NewArtistPosgres(db *sqlx.DB) *ArtistPostgres {
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
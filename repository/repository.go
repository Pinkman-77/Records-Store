package repository

import (
	recordsrestapi "github.com/Pinkman-77/records-restapi"
	"github.com/jmoiron/sqlx"
)

type Creator interface {
	CreateArtist(artist recordsrestapi.Artist) (int, error)
	GetAllArtists() ([]recordsrestapi.ArtistWithRecords, error)
}

type Record interface {

}



type Repository struct {
	Creator
	Record
}

func NewRepository(db sqlx.DB) *Repository {
	return &Repository{
		Creator: NewArtistPostgres(&db),
	}
}
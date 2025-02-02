package repository

import (
        recordsrestapi "github.com/Pinkman-77/records-restapi"
        "github.com/jmoiron/sqlx"
)

type Creator interface {
        CreateArtist(artist recordsrestapi.Artist) (int, error)
        GetAllArtists() ([]recordsrestapi.ArtistWithRecords, error)
        GetArtist(id int) (recordsrestapi.ArtistWithRecords, error)
        UpdateArtist(id int, updatedArtist recordsrestapi.Artist) error
        DeleteArtist(id int) error
        GetUserIDByEmail(email string) (int, error) // New method for getting user id by email

}


type Record interface {
        CreateRecord(record recordsrestapi.Record) (int, error)
        GetAllRecords() ([]recordsrestapi.Record, error)
        GetRecord(id int) (recordsrestapi.RecordWithArtist, error) 
        UpdateRecord(id int, updatedRecord recordsrestapi.Record) error
        PatchRecord(id int, updates map[string]interface{}) error
        DeleteRecord(id int) error

}

type Repository struct {
        Creator
        Record
}

func NewRepository(db sqlx.DB) *Repository {
        return &Repository{
                Creator: NewArtistPostgres(&db),
                Record:  NewRecordPostgres(&db), // Initialize RecordPostgres
        }
}
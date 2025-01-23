package service

import (
        recordsrestapi "github.com/Pinkman-77/records-restapi"
        "github.com/Pinkman-77/records-restapi/repository"
)

type Creator interface {
        CreateArtist(artist recordsrestapi.Artist) (int, error)
        GetAllArtists() ([]recordsrestapi.ArtistWithRecords, error)
        GetArtist(id int) (recordsrestapi.Artist, error)
        UpdateArtist(id int, updatedArtist recordsrestapi.Artist) error
        DeleteArtist(id int) error
 
}
type Record interface {

}

type Service struct {
        Creator
        Record
}

func NewService(repo *repository.Repository) *Service {
        return &Service{
                Creator:  NewArtistList(repo),
                Record:   NewRecordList(repo), // Initialize RecordList
        }
}
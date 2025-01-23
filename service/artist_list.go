package service

import (
	recordsrestapi "github.com/Pinkman-77/records-restapi"
	"github.com/Pinkman-77/records-restapi/repository"
)


type ArtistList struct {
	repo repository.Creator
}

func NewArtistList(repo repository.Creator) *ArtistList {
	return &ArtistList{repo: repo}
}

func (r *ArtistList) CreateArtist(artist recordsrestapi.Artist) (int, error) {
	return r.repo.CreateArtist(artist)
}

func (r *ArtistList) GetAllArtists() ([]recordsrestapi.ArtistWithRecords, error) {
	return r.repo.GetAllArtists()
}

func (r *ArtistList) GetArtist(id int) (recordsrestapi.Artist, error) {
	return r.repo.GetArtist(id)
}
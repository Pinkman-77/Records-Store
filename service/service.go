package service
import (
	recordsrestapi "github.com/Pinkman-77/records-restapi"
	"github.com/Pinkman-77/records-restapi/repository"
)

type Creator interface {
	CreateArtist(artist recordsrestapi.Artist) (int, error)
	GetAllArtists() ([]recordsrestapi.ArtistWithRecords, error)
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
	}
}

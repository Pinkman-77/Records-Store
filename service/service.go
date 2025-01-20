package service

import "github.com/Pinkman-77/records-restapi/repository"

type Records interface {

}

type Service struct {
	Records
}

func NewService(repo repository.Repository) *Service {
	return &Service{
		// Records: NewRecordsService(repo),
	}
}
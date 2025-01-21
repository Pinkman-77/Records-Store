package service

import (
	"github.com/Pinkman-77/records-restapi/repository"
)

type Creator interface {
}

type Record interface {
}

type RecordItems interface {

}

type Service struct {
	Creator
	Record
	RecordItems
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
	}
}
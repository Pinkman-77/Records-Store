package repository

import (
	"github.com/jmoiron/sqlx"
)

type Creator interface {

}

type Record interface {
}

type RecordItems interface {
}

type Repository struct {
	Creator
	Record
	RecordItems
}

func NewRepository(db sqlx.DB) *Repository {
	return &Repository{
	}
}

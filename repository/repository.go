package repository

import (
	"github.com/jmoiron/sqlx"
)

type Creator interface {

}

type Record interface {
}

type Repository struct {
	Creator
	Record
}

func NewRepository(db sqlx.DB) *Repository {
	return &Repository{
	}
}

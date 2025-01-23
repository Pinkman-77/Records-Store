package repository

import (

    "github.com/jmoiron/sqlx"
)

type RecordPostgres struct {
    db *sqlx.DB
}

func NewRecordPostgres(db *sqlx.DB) *RecordPostgres {
    return &RecordPostgres{db: db}
}


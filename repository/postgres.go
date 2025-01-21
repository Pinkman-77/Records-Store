package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/jackc/pgx/v5" 
	_ "github.com/lib/pq"
)

type Postgres struct {
	db *sqlx.DB
}

var p *Postgres

const (
	// Config
	dbName   = "rap_records_shop" 
	dbUser   = "postgres"
	dbPass   = "postgres" 
	dbHost   = "localhost"
	dbPort   = "5432"
)

const (
	// Tables
	artistTable = "artists"
	albumTable  = "records"
	itemTable   = "items"
)

func Connect() (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass, dbName))

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return db, nil	
}

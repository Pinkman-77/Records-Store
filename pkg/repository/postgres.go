package repository

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	_ "github.com/jackc/pgx/v5" // PostgreSQL driver
	_ "github.com/lib/pq"
)

type Postgres struct {
	db *sqlx.DB
}

var p *Postgres

const (
	// Tables
	artistTable = "artists"
	albumTable  = "records"
)
func LoadConfig() {
	viper.SetConfigName("config")        
	viper.SetConfigType("yaml")           
	viper.AddConfigPath("./configs")   

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
}

func Connect() (*sqlx.DB, error) {
	LoadConfig()

	dbHost := viper.GetString("database.host")
	dbPort := viper.GetString("database.port")
	dbUser := viper.GetString("database.user")
	dbPass := viper.GetString("database.password")
	dbName := viper.GetString("database.name")
	sslMode := viper.GetString("database.sslmode")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPass, dbName, sslMode,
	)

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Verify the connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

package recordsrestapi

import ()

type Artist struct {
	ID     uint     `gorm:"primaryKey"`
	Name   string   `gorm:"not null" binding:"required"` // Name of the artist
	Records []Record `gorm:"foreignKey:ArtistID"`
}

type Record struct {
	ID     string  `json:"id"`      // Unique identifier for the record
	Title  string  `json:"title" binding:"required"`    // Title of the record
	Artist string  `json:"artist"`   // ID or name of the artist (you may want to use an ID reference)
	Year   int64     `json:"year"`    // Release year of the record
	Tracklist []string `json:"tracklist"` // List of song titles in the record
	Credits   []string `json:"credits"`   // List of artists featured on the record
	Duration  string   `json:"duration"`  // Total duration of the record (e.g., "45:30")
}

type ArtistWithRecords struct {
	ID      uint     `json:"id"`
	Name    string   `json:"name"`
	Records []Record  `json:"records"`
}

type RecordWithArtist struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Artist string `json:"artist"` // This will hold the artist's name
	Year   int64    `json:"year"`
}

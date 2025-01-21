package recordsrestapi

type Artist struct {
	ID     uint     `gorm:"primaryKey"`
	Name   string   `gorm:"not null"`
	Records []Record `gorm:"foreignKey:ArtistID"`
}

type Record struct {
	ID     string  `json:"id"`      // Unique identifier for the record
	Title  string  `json:"title" binding:"required"`    // Title of the record
	Artist string  `json:"artist"`   // ID or name of the artist (you may want to use an ID reference)
	Year   int     `json:"year"`    // Release year of the record
}

type Item struct {
	Tracklist []string `json:"tracklist"` // List of song titles in the record
	Credits   []string `json:"credits"`   // List of artists featured on the record
	Duration  string   `json:"duration"`  // Total duration of the record (e.g., "45:30")
}

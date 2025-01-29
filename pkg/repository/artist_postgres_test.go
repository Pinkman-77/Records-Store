package repository

import (
	"testing"
	"fmt"

	"github.com/lib/pq"
	"github.com/zhashkevych/go-sqlxmock"

	recordsrestapi "github.com/Pinkman-77/records-restapi"

	"github.com/stretchr/testify/assert"
)


func TestArtistPostgres_CreateArtist(t *testing.T) {
	// Initialize sqlmock
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Initialize repository
	repo := ArtistPostgres{db: db}

	// Define test cases
	tests := []struct {
		name    string
		mock    func()
		input   recordsrestapi.Artist
		want    int
		wantErr bool
	}{
		{
			name: "Success",
			mock: func() {
				mock.ExpectBegin() // Expect transaction start

				// Mock successful artist insertion
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery("INSERT INTO artists").WithArgs("Kanye West").WillReturnRows(rows)

				mock.ExpectCommit() // Expect transaction commit
			},
			input: recordsrestapi.Artist{
				Name: "Kanye West",
			},
			want: 1,
		},
		{
			name: "Insert Error",
			mock: func() {
				mock.ExpectBegin() // Expect transaction start

				// Mock insertion failure
				mock.ExpectQuery("INSERT INTO artists").WithArgs("Kanye West").WillReturnError(assert.AnError)

				mock.ExpectRollback() // Expect transaction rollback
			},
			input: recordsrestapi.Artist{
				Name: "Kanye West",
			},
			wantErr: true,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up mocks
			tt.mock()

			// Call CreateArtist method
			got, err := repo.CreateArtist(tt.input)

			// Validate results
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			// Ensure all expectations are met
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
func TestArtistPostgres_GetArtist(t *testing.T) {
    // Initialize sqlmock
    db, mock, err := sqlmock.Newx()
    assert.NoError(t, err)
    defer db.Close()

    // Initialize repository
    repo := ArtistPostgres{db: db}

    // Define the mock behavior
    mock.ExpectQuery("SELECT a.id, a.name, r.id AS record_id, r.title, r.year, r.tracklist, r.credits, r.duration FROM artists a LEFT JOIN records r ON a.id = r.artist_id WHERE a.id = \\$1").
        WithArgs(1). // Artist ID: 1 (Pop Smoke)
        WillReturnRows(sqlmock.NewRows([]string{"id", "name", "record_id", "title", "year", "tracklist", "credits", "duration"}).
            AddRow(1, "Pop Smoke", "101", "Shoot for the Stars, Aim for the Moon", 2020, pq.Array([]string{"Dior", "For the Night"}), pq.Array([]string{"Quavo", "Lil Baby"}), "45:12").
            AddRow(1, "Pop Smoke", "102", "Meet the Woo 2", 2020, pq.Array([]string{"Christopher Walking", "Element"}), pq.Array([]string{"50 Cent", "Fivio Foreign"}), "38:45"))

    // Define the expected result
    expected := recordsrestapi.ArtistWithRecords{
        ID:   1,
        Name: "Pop Smoke",
        Records: []recordsrestapi.Record{
            {
                ID:        "101",
                Title:     "Shoot for the Stars, Aim for the Moon",
                Artist:    "Pop Smoke", // Include the artist name
                Year:      2020,
                Tracklist: []string{"Dior", "For the Night"},
                Credits:   []string{"Quavo", "Lil Baby"},
                Duration:  "45:12",
            },
            {
                ID:        "102",
                Title:     "Meet the Woo 2",
                Artist:    "Pop Smoke", // Include the artist name
                Year:      2020,
                Tracklist: []string{"Christopher Walking", "Element"},
                Credits:   []string{"50 Cent", "Fivio Foreign"},
                Duration:  "38:45",
            },
        },
    }

    // Call the GetArtist method
    got, err := repo.GetArtist(1) // ID: 1 (Pop Smoke)

    // Assertions
    assert.NoError(t, err)
    assert.Equal(t, expected, got)

    // Ensure all expectations were met
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestArtistPostgres_GetAllArtists(t *testing.T) {
	// Initialize sqlmock
	db, mock, err := sqlmock.Newx()
	assert.NoError(t, err)
	defer db.Close()

	// Initialize repository
	repo := ArtistPostgres{db: db}

	// Define test cases
	tests := []struct {
		name    string
		mock    func()
		want    []recordsrestapi.ArtistWithRecords
		wantErr bool
	}{
		{
			name: "Success",
			mock: func() {
				// Mock rows for artists and their records
				rows := sqlmock.NewRows([]string{"id", "name", "record_id", "title", "year", "tracklist", "credits", "duration"}).
					AddRow(1, "Pop Smoke", "101", "Shoot for the Stars, Aim for the Moon", 2020, pq.Array([]string{"Dior", "For the Night"}), pq.Array([]string{"Quavo", "Lil Baby"}), "45:12").
					AddRow(1, "Pop Smoke", "102", "Meet the Woo 2", 2020, pq.Array([]string{"Christopher Walking", "Element"}), pq.Array([]string{"50 Cent", "Fivio Foreign"}), "38:45").
					AddRow(2, "Jay-Z", nil, nil, nil, nil, nil, nil) // Jay-Z with no records

				mock.ExpectQuery("SELECT a.id, a.name, r.id AS record_id, r.title, r.year, r.tracklist, r.credits, r.duration FROM artists a LEFT JOIN records r ON a.id = r.artist_id").
					WillReturnRows(rows)
			},
			want: []recordsrestapi.ArtistWithRecords{
				{
					ID:   1,
					Name: "Pop Smoke",
					Records: []recordsrestapi.Record{
						{
							ID:        "101",
							Title:     "Shoot for the Stars, Aim for the Moon",
							Artist:    "Pop Smoke",
							Year:      2020,
							Tracklist: []string{"Dior", "For the Night"},
							Credits:   []string{"Quavo", "Lil Baby"},
							Duration:  "45:12",
						},
						{
							ID:        "102",
							Title:     "Meet the Woo 2",
							Artist:    "Pop Smoke",
							Year:      2020,
							Tracklist: []string{"Christopher Walking", "Element"},
							Credits:   []string{"50 Cent", "Fivio Foreign"},
							Duration:  "38:45",
						},
					},
				},
				{
					ID:      2,
					Name:    "Jay-Z",
					Records: []recordsrestapi.Record{},
				},
			},
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up mocks
			tt.mock()

			// Call GetAllArtists method
			got, err := repo.GetAllArtists()

			// Validate results
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			// Ensure all expectations are met
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
func TestArtistPostgres_DeleteArtist(t *testing.T) {
	// Initialize sqlmock
	db, mock, err := sqlmock.Newx()
	assert.NoError(t, err)
	defer db.Close()

	// Initialize repository
	repo := ArtistPostgres{db: db}

	// Define test cases
	tests := []struct {
		name    string
		mock    func()
		input   int
		wantErr bool
	}{
		{
			name: "Success",
			mock: func() {
				mock.ExpectExec("DELETE FROM artists WHERE id = \\$1").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(0, 1)) // Simulate 1 row affected
			},
			input: 1,
		},
		{
			name: "Artist Not Found",
			mock: func() {
				mock.ExpectExec("DELETE FROM artists WHERE id = \\$1").
					WithArgs(99).
					WillReturnResult(sqlmock.NewResult(0, 0)) // Simulate 0 rows affected
			},
			input:   99,
			wantErr: true,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up mocks
			tt.mock()

			// Call DeleteArtist method
			err := repo.DeleteArtist(tt.input)

			// Validate results
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// Ensure all expectations are met
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestArtistPostgres_UpdateArtist(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	assert.NoError(t, err)
	defer db.Close()

	repo := ArtistPostgres{db: db}


	tests := []struct {
		name    string
		mock    func()
		inputID int
		input   recordsrestapi.Artist
		wantErr bool
	}{
		{
			name: "Success",
			mock: func() {
				mock.ExpectExec("UPDATE artists SET name = \\$1 WHERE id = \\$2").
					WithArgs("Pop Smoke", 1).
					WillReturnResult(sqlmock.NewResult(0, 1)) 
			},
			inputID: 1,
			input:   recordsrestapi.Artist{Name: "Pop Smoke"},
			wantErr: false,
		},
		{
			name: "Artist Not Found",
			mock: func() {
				mock.ExpectExec("UPDATE artists SET name = \\$1 WHERE id = \\$2").
					WithArgs("Unknown", 99).
					WillReturnResult(sqlmock.NewResult(0, 0)) 
			},
			inputID: 99,
			input:   recordsrestapi.Artist{Name: "Unknown"},
			wantErr: false, // The method doesn't check RowsAffected(), so it won't return an error
		},
		{
			name: "Database Error",
			mock: func() {
				mock.ExpectExec("UPDATE artists SET name = \\$1 WHERE id = \\$2").
					WithArgs("Kanye West", 2).
					WillReturnError(fmt.Errorf("db error")) 
			},
			inputID: 2,
			input:   recordsrestapi.Artist{Name: "Kanye West"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := repo.UpdateArtist(tt.inputID, tt.input)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}





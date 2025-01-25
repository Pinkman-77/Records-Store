package repository

import (
	"testing"
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


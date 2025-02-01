package repository

import (
	"database/sql"
	"testing"

	recordsrestapi "github.com/Pinkman-77/records-restapi"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
)

func TestRecordsPostgres_CreateRecord(t *testing.T) {

	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	repo := RecordPostgres{db:db}

	tests := []struct {
		name string
		mock func()
		input recordsrestapi.Record
		want int
		wantErr bool

	}{
		{
			name: "Green light",
			mock: func() {
				mock.ExpectQuery("SELECT id FROM artists WHERE name = \\$1").
				WithArgs("Kanye West"). 
				WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectBegin()

				mock.ExpectQuery("INSERT INTO records").WithArgs(
					"Donda", 1, 2011,
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					"46:12",
				).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(10))

				mock.ExpectCommit()

			},
			input: recordsrestapi.Record{
				Title:     "Donda",
				Artist:    "Kanye West",
				Year:      2011,
				Tracklist: []string{"No Church In The Wild", "Otis"},
				Credits:   []string{"Frank Ocean", "Beyoncé"},
				Duration:  "46:12",
			},
			want: 10,	

		},

		{
			name: "Red light",
			mock: func() {
				mock.ExpectQuery("SELECT id FROM artists WHERE name = \\$1").
				WithArgs("Kanye West"). 
				WillReturnError(sql.ErrNoRows)

			},
			input: recordsrestapi.Record{
				Title:     "Donda",
				Artist:    "Kanye West",
				Year:      2011,
				Tracklist: []string{"No Church In The Wild", "Otis"},
				Credits:   []string{"Frank Ocean", "Beyoncé"},
				Duration:  "46:12",
			},
			wantErr: true,
		},

		{
			name: "2nd Red light",
			mock: func() {
				mock.ExpectQuery("SELECT id FROM artists WHERE name = \\$1").
				WithArgs("Kanye West"). 
				WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectBegin()

				mock.ExpectQuery("INSERT INTO records").WithArgs(
					"Donda", 1, 2011,
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					"46:12",
				).WillReturnError(assert.AnError)

				mock.ExpectRollback()

			},
			input: recordsrestapi.Record{
				Title:     "Donda",
				Artist:    "Kanye West",
				Year:      2011,
				Tracklist: []string{"No Church In The Wild", "Otis"},
				Credits:   []string{"Frank Ocean", "Beyoncé"},
				Duration:  "46:12",
			},
			wantErr: true,

			},

		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				tt.mock()
				got, err := repo.CreateRecord(tt.input)
				if (err != nil) != tt.wantErr {
					t.Errorf("CreateRecord() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("CreateRecord() = %v, want %v", got, tt.want)
				}
			})

		}

	}

	func TestArtistPostgres_GetRecord(t *testing.T) {
		db, mock, err := sqlmock.Newx()
		assert.NoError(t, err)
		defer db.Close()
	
		repo := RecordPostgres{db: db}
	
		mock.ExpectQuery(`SELECT r.id, r.title, a.name AS artist, r.year, r.tracklist, r.credits, r.duration 
			FROM records r 
			INNER JOIN artists a ON r.artist_id = a.id 
			WHERE r.id = \$1`).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "title", "artist", "year", "tracklist", "credits", "duration"}).
				AddRow(1, "Album Title", "Artist Name", 2022, pq.Array([]string{"Track 1", "Track 2"}), pq.Array([]string{"Credit 1", "Credit 2"}), "43:52"))
	
		record, err := repo.GetRecord(1)
		assert.NoError(t, err)
		assert.Equal(t, "Album Title", record.Title)
		assert.Equal(t, "Artist Name", record.Artist)
		assert.Equal(t, "43:52", record.Duration)
		assert.NoError(t, mock.ExpectationsWereMet())
	}
	
	func TestRecordPostgres_GetAllRecords(t *testing.T) {
		db, mock, err := sqlmock.Newx()
		assert.NoError(t, err)
		defer db.Close()
	
		repo := RecordPostgres{db: db}
	
		mock.ExpectQuery(`SELECT r.id, r.title, a.name AS artist, r.year, r.tracklist, r.credits, r.duration 
			FROM records r 
			INNER JOIN artists a ON r.artist_id = a.id`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "title", "artist", "year", "tracklist", "credits", "duration"}).
				AddRow(1, "Album Title", "Artist Name", 2022, pq.Array([]string{"Track 1", "Track 2"}), pq.Array([]string{"Credit 1", "Credit 2"}), "43:52").
				AddRow(2, "Another Album", "Another Artist", 2021, pq.Array([]string{"Track A", "Track B"}), pq.Array([]string{"Credit A", "Credit B"}), "41:30"))
	
		records, err := repo.GetAllRecords()
		assert.NoError(t, err)
		assert.Len(t, records, 2)
		assert.Equal(t, "Album Title", records[0].Title)
		assert.Equal(t, "Another Album", records[1].Title)
		assert.NoError(t, mock.ExpectationsWereMet())
	}
	
	func TestArtistPostgres_DeleteRecord(t *testing.T) {
		db, mock, err := sqlmock.Newx()
		assert.NoError(t, err)
		defer db.Close()
	
		repo := RecordPostgres{db: db}
	
		mock.ExpectExec("DELETE FROM records WHERE id = \\$1").
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 1))
	
		err = repo.DeleteRecord(1)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	}
	
	func TestArtistPostgres_UpdateRecord(t *testing.T) {
		db, mock, err := sqlmock.Newx()
		assert.NoError(t, err)
		defer db.Close()
	
		repo := RecordPostgres{db: db}
	
		tests := []struct {
			name    string
			mock    func()
			inputID int
			input   recordsrestapi.Record
			wantErr bool
		}{
			{
				name: "Success",
				mock: func() {
					mock.ExpectExec(`UPDATE records SET title = \$1, artist_id = \(SELECT id FROM artists WHERE name = \$2\), year = \$3, tracklist = \$4, credits = \$5, duration = \$6 WHERE id = \$7`).
						WithArgs("Album Title", "Artist Name", 2022, pq.Array([]string{"Track 1", "Track 2"}), pq.Array([]string{"Credit 1", "Credit 2"}), "43:52", 1).
						WillReturnResult(sqlmock.NewResult(0, 1))
				},
				inputID: 1,
				input: recordsrestapi.Record{
					Title:     "Album Title",
					Artist:    "Artist Name",
					Year:      2022,
					Tracklist: []string{"Track 1", "Track 2"},
					Credits:   []string{"Credit 1", "Credit 2"},
					Duration:  "43:52",
				},
				wantErr: false,
			},
		}
	
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				tt.mock()
				err := repo.UpdateRecord(tt.inputID, tt.input)
	
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
	
				assert.NoError(t, mock.ExpectationsWereMet())
			})
		}
	}
	
	func TestArtistPostgres_PatchRecord(t *testing.T) {
		db, mock, err := sqlmock.Newx()
		assert.NoError(t, err)
		defer db.Close()
	
		repo := RecordPostgres{db: db}
	
		tests := []struct {
			name    string
			mock    func()
			inputID int
			updates map[string]interface{}
			wantErr bool
		}{
			{
				name: "Success",
				mock: func() {
					mock.ExpectExec("UPDATE records SET title = \\$2 WHERE id = \\$1").
						WithArgs(1, "New Album").
						WillReturnResult(sqlmock.NewResult(0, 1))
				},
				inputID: 1,
				updates: map[string]interface{}{"title": "New Album"},
				wantErr: false,
			},
		}
	
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				tt.mock()
				err := repo.PatchRecord(tt.inputID, tt.updates)
	
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
	
				assert.NoError(t, mock.ExpectationsWereMet())
			})
		}
	}
	
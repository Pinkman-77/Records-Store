package repository

import (
	"database/sql"
	"testing"

	recordsrestapi "github.com/Pinkman-77/records-restapi"
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

		

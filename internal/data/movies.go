package data

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/mahmoud-shabban/greenlight/internal/validator"
)

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int64     `json:"year,omitempty"`
	Runtime   Runtime   `json:"runtime,omitempty"`
	Genres    []string  `json:"genres,omitempty"`
	Version   int32     `json:"version"`
}

type MovideModel struct {
	DB *sql.DB
}

func ValidateMovie(validations *validator.Validator, movie *Movie) {

	validations.Check(movie.Title != "", "title", "title must not be empty")
	validations.Check(len(movie.Title) <= 500, "title", "title must be <= 500 bytes long")

	validations.Check(movie.Runtime.Duration != 0, "runtime", "runtime must be provided")
	validations.Check(movie.Runtime.Duration > 0, "runtime", "runtime must be positive integer")

	validations.Check(movie.Year != 0, "year", "year must be provided")
	validations.Check(movie.Year >= 1888, "year", "year must be greater than 1888")
	validations.Check(movie.Year <= int64(time.Now().Year()), "year", fmt.Sprintf("year must be less than or equal %d", time.Now().Year()))

	validations.Check(movie.Genres != nil, "genres", "genres must provided")
	validations.Check(1 <= len(movie.Genres) && len(movie.Genres) <= 5, "genres", "genres must contain between 1 and 5 genres")
	validations.Check(validator.Unique(movie.Genres), "genres", "genres must be unique")

}

func (m MovideModel) Insert(movie *Movie) error {
	stmt := `
			INSERT INTO movies (title, year, runtime, genres)
			VALUES($1, $2, $3,$4)
			RETURNING id, created_at, version
	`

	args := []any{movie.Title, movie.Year, movie.Runtime.Duration, pq.Array(movie.Genres)}

	return m.DB.QueryRow(stmt, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

func (m MovideModel) Get(id int64) (*Movie, error) {
	return &Movie{}, nil
}

func (m MovideModel) Update(movie *Movie) error {
	return nil
}

func (m MovideModel) Delete(id int64) error {
	return nil

}

package data

import (
	"context"
	"database/sql"
	"errors"
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

	validations.Check(movie.Runtime != 0, "runtime", "runtime must be provided")
	validations.Check(movie.Runtime > 0, "runtime", "runtime must be positive integer")

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

	args := []any{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}

	return m.DB.QueryRow(stmt, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

func (m MovideModel) Get(id int64) (*Movie, error) {

	if id < 1 {
		return nil, ErrRecoredNotFound
	}

	stmt := `
			SELECT id, created_at, title, year, runtime, genres, version
			FROM movies
			WHERE id = $1
		`

	result := Movie{}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()
	err := m.DB.QueryRowContext(ctx, stmt, id).Scan(
		// err := m.DB.QueryRow(stmt, id).Scan(
		// &[]byte{},
		&result.ID,
		&result.CreatedAt,
		&result.Title,
		&result.Year,
		&result.Runtime,
		pq.Array(&result.Genres),
		&result.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecoredNotFound
		default:
			return nil, err
		}
	}

	return &result, nil
}

func (m MovideModel) Update(movie *Movie) error {

	stmt := `
		UPDATE movies
		SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1
		WHERE id = $5 AND version = $6
		RETURNING version
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()
	// err := m.DB.QueryRow(
	err := m.DB.QueryRowContext(
		ctx,
		stmt,
		movie.Title,
		movie.Year,
		movie.Runtime,
		pq.Array(movie.Genres),
		movie.ID,
		movie.Version,
	).Scan(&movie.Version)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil

}

func (m MovideModel) Delete(id int64) error {

	stmt := `
		DELETE FROM movies
		WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()
	result, err := m.DB.ExecContext(ctx, stmt, id)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecoredNotFound
	}

	return nil
}

func (m MovideModel) GetAll(title string, genres []string, filrters Filters) ([]*Movie, Metadata, error) {
	stmt := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), id, created_at, title, year, runtime, genres, version
		FROM movies
		WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (genres @> $2 OR $2 = '{}')
		ORDER BY %s %s, id ASC
		LIMIT %d
		OFFSET (%d - 1) * %d
		`,
		filrters.sortColumn(),
		filrters.sortDirection(),
		filrters.PageSize,
		filrters.Page,
		filrters.PageSize,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	rows, err := m.DB.QueryContext(ctx, stmt, title, pq.Array(genres))

	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	movies := []*Movie{}

	for rows.Next() {
		var movie Movie

		err := rows.Scan(
			&totalRecords,
			&movie.ID,
			&movie.CreatedAt,
			&movie.Title,
			&movie.Year,
			&movie.Runtime,
			pq.Array(&movie.Genres),
			&movie.Version,
		)

		if err != nil {
			return nil, Metadata{}, err
		}

		movies = append(movies, &movie)
	}

	meta := calculateMetadata(totalRecords, filrters.Page, filrters.PageSize)

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}
	return movies, meta, nil
}

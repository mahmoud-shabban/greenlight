package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecoredNotFound = errors.New("record not found")
)

type Models struct {
	Movies MovideModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Movies: MovideModel{DB: db},
	}
}

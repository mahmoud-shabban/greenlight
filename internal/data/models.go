package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecoredNotFound = errors.New("record not found")
	ErrEditConflict    = errors.New("edit conflict")
)

type Models struct {
	Movies      MovideModel
	Users       UserModel
	Tokens      TokenModel
	Permissions PermissionModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Movies:      MovideModel{DB: db},
		Users:       UserModel{DB: db},
		Tokens:      TokenModel{DB: db},
		Permissions: PermissionModel{DB: db},
	}
}

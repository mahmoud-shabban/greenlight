package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/mahmoud-shabban/greenlight/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	Activated bool      `json:"activated"`
	Version   int       `json:"version"`
}

var (
	ErrDublicateEmail = errors.New("duplicate email")
)

type UserModel struct {
	DB *sql.DB
}

type password struct {
	plainText *string
	hash      []byte
}

func (p *password) Set(plainText string) error {

	hash, err := bcrypt.GenerateFromPassword([]byte(plainText), 12)

	if err != nil {
		return err
	}

	p.plainText = &plainText
	p.hash = hash

	return nil
}

func (p *password) Matches(plainText string) (bool, error) {

	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plainText))

	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil

}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidatePasswordPlainText(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) <= 72, "password", "must not be more than 75 bytes long")
	v.Check(len(password) >= 8, "password", "must not be less than 8 bytes long")
}

func ValidateUser(v *validator.Validator, user *User) {

	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) <= 500, "name", "mut not be more thant 500 bytes long")

	ValidateEmail(v, user.Email)
	ValidatePasswordPlainText(v, *user.Password.plainText)

	if user.Password.hash == nil {
		panic("missing password hash for user")
	}

}

func (u *UserModel) Insert(user *User) error {
	stmt := `
		INSERT INTO users(name,email, password_hash, activated)
		VALUES($1, $2, $3, $4)
		RETURNING id, created_at, version
	`

	args := []any{user.Name, user.Email, user.Password.hash, user.Activated}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	err := u.DB.QueryRowContext(ctx, stmt, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDublicateEmail
		default:
			return err
		}
	}

	return nil
}

func (u *UserModel) GetByEmail(email string) (*User, error) {

	stmt := `
		SELECT id, created_at, name, email, password_hash, activated, version 
		FROM USERS 
		WHERE email = $1
	`

	user := User{}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := u.DB.QueryRowContext(ctx, stmt, email).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
		&user.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecoredNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (u *UserModel) Update(user *User) error {

	stmt := `
		UPDATE users
		SET name = $1, email = $2, password_hash = $3, activated = $4, version = version + 1
		WHERE id = $5 AND version = $6
		RETURNING version
	`

	args := []any{user.Name, user.Email, user.Password.hash, user.Activated, user.ID, user.Version}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	err := u.DB.QueryRowContext(ctx, stmt, args...).Scan(&user.Version)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDublicateEmail
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		}
	}
	return nil
}

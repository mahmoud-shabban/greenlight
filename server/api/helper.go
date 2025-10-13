package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	// _ "github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

func (app *Application) readIDParam(params httprouter.Params) (int64, error) {
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)

	if err != nil || id <= 0 {
		// fmt.Fprintf(w, "movie id must be positive integer\n")
		// http.NotFound(w, r)
		return 0, fmt.Errorf("invalid id parameter")
	}

	return id, nil

}

type envelope map[string]any

func (app *Application) writeJson(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")

	if err != nil {
		return err
	}

	js = append(js, '\n')
	for k, v := range headers {
		w.Header()[k] = v
	}

	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func (app *Application) readJson(w http.ResponseWriter, r *http.Request, dest any) error {

	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dest)

	if err != nil {
		var syntaxError *json.SyntaxError
		var umMarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return fmt.Errorf("body contains badly-formed JSON")
		case errors.As(err, &umMarshalTypeError):
			if umMarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", umMarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", umMarshalTypeError.Offset)
		case errors.Is(err, io.EOF):
			return fmt.Errorf("body mut not be empty")
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			field := strings.TrimPrefix(err.Error(), "json: unknown field")
			return fmt.Errorf("body contains unknown field %s", field)
		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err

		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return fmt.Errorf("body must only contain single JSON value")
	}
	return nil
}

func (app *Application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	err = db.PingContext(ctx)
	defer cancel()

	if err != nil {
		return nil, err
	}
	return db, nil

}

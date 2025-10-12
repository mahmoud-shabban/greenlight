package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/mahmoud-shabban/greenlight/internal/data"
	"github.com/mahmoud-shabban/greenlight/internal/validator"
)

func (app *Application) createMovieHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var input struct {
		Title   string       `json:"title"`
		Year    int64        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}
	err := app.readJson(w, r, &input)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	movie := data.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	validations := validator.New()
	data.ValidateMovie(validations, &movie)
	if !validations.Valid() {
		app.faildValidationResponse(w, r, validations.Errors)
		return
	}

	err = app.models.Movies.Insert(&movie)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)

	headers.Set("Location", fmt.Sprintf("v1/movies/%d", movie.ID))
	// w.WriteHeader(http.StatusCreated)
	err = app.writeJson(w, http.StatusCreated, envelope{"movie": input}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	// fmt.Fprintf(w, "%+v\n", input)
}

func (app *Application) showMovieHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	// params = httprouter.ParamsFromContext(r.Context())

	id, err := app.readIDParam(params)

	if err != nil {
		// http.Error(w, err.Error(), http.StatusNotFound)
		app.notFoundResponse(w, r)
		return
	}

	movie := data.Movie{
		ID:        int64(id),
		CreatedAt: time.Now(),
		Title:     "Test Movie",
		Runtime:   data.Runtime{Duration: 102, Unit: "mins"},
		Genres:    []string{"Action", "Comedy"},
		Version:   1,
	}

	err = app.writeJson(w, http.StatusOK, envelope{"movie": movie}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

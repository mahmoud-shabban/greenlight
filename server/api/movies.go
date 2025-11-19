package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mahmoud-shabban/greenlight/internal/data"
	"github.com/mahmoud-shabban/greenlight/internal/validator"
)

func (app *Application) createMovieHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	_, span := app.config.tracer.Start(r.Context(), "createMovie")
	defer span.End()

	var input struct {
		Title   string       `json:"title"`
		Year    int64        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}

	span.AddEvent("read body data")
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

	span.AddEvent("validating data")
	validations := validator.New()
	data.ValidateMovie(validations, &movie)
	if !validations.Valid() {
		app.faildValidationResponse(w, r, validations.Errors)
		return
	}

	span.AddEvent("insert movie data")
	err = app.models.Movies.Insert(&movie)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)

	headers.Set("Location", fmt.Sprintf("v1/movies/%d", movie.ID))
	// w.WriteHeader(http.StatusCreated)

	span.AddEvent("sending response json")
	err = app.writeJson(w, http.StatusCreated, envelope{"movie": input}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	// fmt.Fprintf(w, "%+v\n", input)
}

func (app *Application) showMovieHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	// params = httprouter.ParamsFromContext(r.Context())

	_, span := app.config.tracer.Start(r.Context(), "show movie")
	defer span.End()

	id, err := app.readIDParam(params)

	if err != nil {
		// http.Error(w, err.Error(), http.StatusNotFound)
		app.notFoundResponse(w, r)
		return
	}

	// movie := data.Movie{
	// 	ID:        int64(id),
	// 	CreatedAt: time.Now(),
	// 	Title:     "Test Movie",
	// 	Runtime:   data.Runtime{Duration: 102, Unit: "mins"},
	// 	Genres:    []string{"Action", "Comedy"},
	// 	Version:   1,
	// }

	span.AddEvent("query movie with id")
	movie, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecoredNotFound):
			app.notFoundResponse(w, r)
			return
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	span.AddEvent("sending response")
	err = app.writeJson(w, http.StatusOK, envelope{"movie": movie}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *Application) updateMovieHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	_, span := app.config.tracer.Start(r.Context(), "update movie")
	defer span.End()

	id, err := app.readIDParam(params)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	span.AddEvent("query movie with id")
	// var movie data.Movie
	movie, err := app.models.Movies.Get(id)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecoredNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	input := struct {
		Title   *string       `json:"title"`
		Year    *int64        `json:"year"`
		Genres  []string      `json:"genres"`
		Runtime *data.Runtime `json:"runtime"`
	}{}

	span.AddEvent("read request body")
	err = app.readJson(w, r, &input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if input.Title != nil {
		movie.Title = *input.Title
	}

	if input.Year != nil {
		movie.Year = *input.Year
	}
	if input.Genres != nil {
		movie.Genres = input.Genres
	}

	if input.Runtime != nil {
		movie.Runtime = *input.Runtime
	}

	span.AddEvent("validating movie data")
	validations := validator.New()

	data.ValidateMovie(validations, movie)

	if !validations.Valid() {
		app.faildValidationResponse(w, r, validations.Errors)
		return
	}

	span.AddEvent("update datbase movie data")
	if err = app.models.Movies.Update(movie); err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	span.AddEvent("sending response")
	if err = app.writeJson(w, http.StatusOK, envelope{"movie": movie}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) deleteMovieHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	_, span := app.config.tracer.Start(r.Context(), "delete movie")
	defer span.End()

	id, err := app.readIDParam(params)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	span.AddEvent("delete db data")
	err = app.models.Movies.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecoredNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	span.AddEvent("sending response")
	err = app.writeJson(w, http.StatusOK, envelope{"message": "movie delete successfully"}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *Application) listMoviesHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	_, span := app.config.tracer.Start(r.Context(), "list_movie_handler")
	defer span.End()

	// span.SetName("list_movies_handler")
	var input struct {
		Tittle string
		Genres []string
		data.Filters
	}

	span.AddEvent("reading url query string and validating")
	v := validator.New()

	qs := r.URL.Query()

	input.Tittle = app.readString(qs, "title", "")
	input.Genres = app.readCSV(qs, "genres", []string{})
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")

	input.Filters.SortSafelist = []string{"id", "-id", "runtime", "-runtime", "year", "-year", "title", "-title"}

	data.ValidateFilters(v, input.Filters)

	if !v.Valid() {
		app.faildValidationResponse(w, r, v.Errors)
		return
	}

	span.AddEvent("getting all movies")
	movies, meta, err := app.models.Movies.GetAll(input.Tittle, input.Genres, input.Filters)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	span.AddEvent("sending response")

	err = app.writeJson(w, http.StatusOK, envelope{"metadata": meta, "movies": movies}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// fmt.Fprintf(w, "%+v\n", input)

}

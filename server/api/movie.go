package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/mahmoud-shabban/greenlight/internal/data"
)

func (app *Application) createMovieHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var input struct {
		Title   string   `json:"title"`
		Year    int64    `json:"year"`
		Runtime int64    `json:"runtime"`
		Genres  []string `json:"genres"`
	}

	body := r.Body

	defer body.Close()

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%+v\n", input)
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

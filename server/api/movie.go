package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/mahmoud-shabban/greenlight/internal/data"
)

func (app *Application) createMovieHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Fprintf(w, "create new movie\n")
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

package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *Application) createMovieHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Fprintf(w, "create new movie\n")
}

func (app *Application) showMovieHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	// params = httprouter.ParamsFromContext(r.Context())

	id, err := app.readIDParam(params)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "show movie with id: %d\n", id)
}

package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *Application) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	err := fmt.Sprintf("path %s not found", r.URL.Path)
	http.Error(w, err, http.StatusNotFound)
}

func (app *Application) routes() *httprouter.Router {

	router := &httprouter.Router{
		RedirectTrailingSlash: true,
		NotFound:              http.HandlerFunc(app.notFoundHandler),
	}

	router.GET("/v1/healthcheck", app.healthCheckeHandler)
	router.POST("/v1/movies", app.createMovieHandler)
	router.GET("/v1/movie/:id", app.showMovieHandler)

	return router

}

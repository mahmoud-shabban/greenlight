package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *Application) routerNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	app.notFoundResponse(w, r)
}

func (app *Application) routes() *httprouter.Router {

	router := &httprouter.Router{
		RedirectTrailingSlash:  true,
		HandleMethodNotAllowed: true,
		NotFound:               http.HandlerFunc(app.routerNotFoundHandler),
		MethodNotAllowed:       http.HandlerFunc(app.methodNotAllowedResponse),
	}

	router.GET("/v1/healthcheck", app.healthCheckeHandler)
	router.POST("/v1/movies", app.createMovieHandler)
	router.GET("/v1/movie/:id", app.showMovieHandler)

	return router

}

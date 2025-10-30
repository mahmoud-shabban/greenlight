package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *Application) routerNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	app.notFoundResponse(w, r)
}

func (app *Application) routes() http.Handler {

	router := httprouter.Router{
		RedirectTrailingSlash:  true,
		HandleMethodNotAllowed: true,
		NotFound:               http.HandlerFunc(app.routerNotFoundHandler),
		MethodNotAllowed:       http.HandlerFunc(app.methodNotAllowedResponse),
	}

	router.GET("/v1/healthcheck", app.healthCheckeHandler)
	router.POST("/v1/movies", app.createMovieHandler)
	router.GET("/v1/movies", app.listMoviesHandler)
	router.GET("/v1/movies/:id", app.showMovieHandler)
	router.PATCH("/v1/movies/:id", app.updateMovieHandler)
	router.DELETE("/v1/movies/:id", app.deleteMovieHandler)

	return app.recoverPanic(&router)

}

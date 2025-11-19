package main

import (
	"expvar"
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

	router.POST("/v1/movies", app.requirePermissions("movies:write", app.createMovieHandler))
	router.GET("/v1/movies", app.requirePermissions("movies:read", app.listMoviesHandler))
	router.GET("/v1/movies/:id", app.requirePermissions("movies:read", app.showMovieHandler))
	router.PATCH("/v1/movies/:id", app.requirePermissions("movies:write", app.updateMovieHandler))
	router.DELETE("/v1/movies/:id", app.requirePermissions("movies:write", app.deleteMovieHandler))

	router.POST("/v1/users", app.registerUserHandler)
	router.PUT("/v1/users/activated", app.activateUserHandler)

	router.POST("/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())
	return app.logRequest(app.metrics(app.recoverPanic(app.rateLimit(app.authenticate(app.enableCORS(&router))))))
}

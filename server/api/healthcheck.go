package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *Application) healthCheckeHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	data := map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     version,
	}

	headers := http.Header{
		"Content-Type": []string{"application/json"},
	}

	err := app.writeJson(w, http.StatusOK, data, headers)

	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (app *Application) readIDParam(params httprouter.Params) (int, error) {
	id, err := strconv.Atoi(params.ByName("id"))

	if err != nil || id <= 0 {
		// fmt.Fprintf(w, "movie id must be positive integer\n")
		// http.NotFound(w, r)
		return 0, fmt.Errorf("invalid id parameter")
	}

	return id, nil

}

func (app *Application) writeJson(w http.ResponseWriter, status int, data any, headers http.Header) error {
	js, err := json.Marshal(data)

	if err != nil {
		// app.logger.Error(err.Error())
		// http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
		return err
	}

	js = append(js, '\n')
	for k, v := range headers {
		w.Header()[k] = v

	}

	w.WriteHeader(status)
	w.Write(js)

	return nil
}

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

type envelope map[string]any

func (app *Application) writeJson(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")

	if err != nil {
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

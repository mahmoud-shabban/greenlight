package main

import (
	"fmt"
	"net/http"
)

func (app *application) healthChecker(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "status: avaliable\n")
	fmt.Fprintf(w, "environment: %s\n", app.config.env)
	fmt.Fprintf(w, "version: %s\n", version)
}

package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/mahmoud-shabban/greenlight/internal/data"
	"github.com/mahmoud-shabban/greenlight/internal/validator"
)

func (app *Application) createAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	_, span := app.config.tracer.Start(r.Context(), "create_auth_token")
	defer span.End()

	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	span.AddEvent("reading body")
	err := app.readJson(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	span.AddEvent("validation")
	v := validator.New()

	data.ValidateEmail(v, input.Email)
	data.ValidatePasswordPlainText(v, input.Password)

	if !v.Valid() {
		app.faildValidationResponse(w, r, v.Errors)
		return
	}

	span.AddEvent("query user")
	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecoredNotFound):
			app.invalidCredentialsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	span.AddEvent("compare passwords")
	match, err := user.Password.Matches(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if !match {
		app.invalidCredentialsResponse(w, r)
		return
	}

	span.AddEvent("comapreing password hash")
	token, err := app.models.Tokens.New(user.ID, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	span.AddEvent("sending response")
	err = app.writeJson(w, http.StatusCreated, envelope{"authentication_token": token}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

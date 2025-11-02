package main

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mahmoud-shabban/greenlight/internal/data"
	"github.com/mahmoud-shabban/greenlight/internal/validator"
)

func (app *Application) registerUserHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJson(w, r, &input)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()

	if data.ValidateUser(v, user); !v.Valid() {
		app.faildValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Users.Insert(user)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrDublicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.faildValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	err = app.writeJson(w, http.StatusCreated, envelope{"user": user}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

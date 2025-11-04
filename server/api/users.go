package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

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

	token, err := app.models.Tokens.New(user.ID, 3*24*time.Hour, data.ScopeActivation)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.background(func() {
		data := map[string]any{
			"userID":          user.ID,
			"activationToken": token.Plaintext,
		}
		err = app.mailer.Send(user.Email, "user_welcome.tmpl.html", data)
		if err != nil {
			app.logger.PrintError(err, nil)
		}

	})

	err = app.writeJson(w, http.StatusAccepted, envelope{"user": user}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *Application) activateUserHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var input struct {
		TokenPlaintext string `json:"token"`
	}

	err := app.readJson(w, r, &input)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	if data.ValidateTokenPlaintext(v, input.TokenPlaintext); !v.Valid() {
		app.faildValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetForToken(data.ScopeActivation, input.TokenPlaintext)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecoredNotFound):
			v.AddError("token", "invalid or expired activation token")
			app.faildValidationResponse(w, r, v.Errors)
		default:
			fmt.Println("111111111111111111111111111111111")
			app.serverErrorResponse(w, r, err)

		}
		return
	}

	user.Activated = true

	if app.models.Users.Update(user) != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			fmt.Println("************************")
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.models.Tokens.DeleteAllForUser(user.ID, data.ScopeActivation)
	if err != nil {
		fmt.Println("####################################")
		app.serverErrorResponse(w, r, err)
		return
	}

	if app.writeJson(w, http.StatusOK, envelope{"user": user}, nil) != nil {
		fmt.Println("8888888888888888888888888888888888")
		app.serverErrorResponse(w, r, err)
	}

}

package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/mahmoud-shabban/greenlight/internal/data"
	"github.com/mahmoud-shabban/greenlight/internal/validator"
)

func (app *Application) registerUserHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	ctx, span := app.config.tracer.Start(r.Context(), "register user")
	defer span.End()
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	span.AddEvent("reading request data and validating")
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

	span.AddEvent("setting user password")
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

	span.AddEvent("inser user into database")
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

	span.AddEvent("setting user permission in database")
	err = app.models.Permissions.AddForUser(user.ID, "movies:read")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	span.AddEvent("create and sending activation token")
	token, err := app.models.Tokens.New(ctx, user.ID, 3*24*time.Hour, data.ScopeActivation, app.config.tracer)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.background(func() {
		span.AddEvent("background task: sending activation mail")
		data := map[string]any{
			"userID":          user.ID,
			"activationToken": token.Plaintext,
		}
		err = app.mailer.Send(user.Email, "user_welcome.tmpl.html", data)
		if err != nil {
			app.logger.PrintError(err, nil)
		}

	})

	span.AddEvent("sending response")
	err = app.writeJson(w, http.StatusAccepted, envelope{"user": user}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *Application) activateUserHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	_, span := app.config.tracer.Start(r.Context(), "activate user")
	defer span.End()

	var input struct {
		TokenPlaintext string `json:"token"`
	}

	span.AddEvent("reading request data and validating")
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

	span.AddEvent("get all user activation tokens")
	user, err := app.models.Users.GetForToken(data.ScopeActivation, input.TokenPlaintext)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecoredNotFound):
			v.AddError("token", "invalid or expired activation token")
			app.faildValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)

		}
		return
	}

	user.Activated = true

	span.AddEvent("updateing activation status in database")
	if app.models.Users.Update(user) != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	span.AddEvent("deleteing user activation tokens from database")
	err = app.models.Tokens.DeleteAllForUser(user.ID, data.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	span.AddEvent("sending response")
	if app.writeJson(w, http.StatusOK, envelope{"user": user}, nil) != nil {
		app.serverErrorResponse(w, r, err)
	}

}

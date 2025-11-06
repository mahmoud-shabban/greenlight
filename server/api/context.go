package main

import (
	"context"
	"net/http"

	"github.com/mahmoud-shabban/greenlight/internal/data"
)

type contextKey string

const userContextKey = contextKey("user")

func (app *Application) contextSetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func (app *Application) contextGetUser(r *http.Request) *data.User {
	u, ok := r.Context().Value(userContextKey).(*data.User)

	if !ok {
		panic("missing user value in request context")
	}

	return u
}

package api

import (
	"context"
	"github.com/giffone/forum-authentication/internal/object"
	"net/http"
)

type Handler interface {
	Register(ctx context.Context, router *http.ServeMux, session Middleware)
}

type Middleware interface { // middleware for handlers
	CreateSession(ctx context.Context, w http.ResponseWriter, id int, method string) object.Status
	Skip(ctx context.Context, fn func(context.Context, Middleware,
		http.ResponseWriter, *http.Request)) http.HandlerFunc
	CheckSession(ctx context.Context, fn func(context.Context,
		*object.Cookie, object.Status, http.ResponseWriter,
		*http.Request)) http.HandlerFunc
	EndSession(w http.ResponseWriter) object.Status
}

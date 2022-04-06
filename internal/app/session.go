package app

import (
	"github.com/giffone/forum-authentication/internal/adapters/api"
	middleware2 "github.com/giffone/forum-authentication/internal/adapters/api/middleware"
	"github.com/giffone/forum-authentication/internal/adapters/repository"
	"github.com/giffone/forum-authentication/internal/service"
	"github.com/giffone/forum-authentication/internal/service/middleware"
)

func (a *App) middlewareService(repo repository.Repo) (service.Middleware, api.Middleware) {
	srv := middleware.NewService(repo)
	return srv, middleware2.NewMiddleware(srv)
}

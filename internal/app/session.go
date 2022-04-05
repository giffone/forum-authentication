package app

import (
	"github.com/giffone/forum-authentication/internal/adapters/api"
	"github.com/giffone/forum-authentication/internal/adapters/api/middleware"
	"github.com/giffone/forum-authentication/internal/adapters/repository"
	"github.com/giffone/forum-authentication/internal/service/session"
)

func (a *App) session(repo repository.Repo) api.Session {
	srv := session.NewService(repo)
	return middleware.NewSession(srv)
}

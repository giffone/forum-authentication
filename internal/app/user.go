package app

import (
	"github.com/giffone/forum-authentication/internal/adapters/api"
	user2 "github.com/giffone/forum-authentication/internal/adapters/api/user"
	"github.com/giffone/forum-authentication/internal/adapters/repository"
	"github.com/giffone/forum-authentication/internal/service/user"
)

func (a *App) user(repo repository.Repo, ses api.Middleware) {
	srv := user.NewService(repo)
	user2.NewHandler(srv).Register(a.ctx, a.router, ses)
}

package app

import (
	"github.com/giffone/forum-authentication/internal/adapters/api"
	"github.com/giffone/forum-authentication/internal/adapters/api/home"
	"github.com/giffone/forum-authentication/internal/service"
)

func (a *App) home(srvPost service.Post, srvCategory service.Category, sMid api.Middleware) {
	home.NewHandler(srvPost, srvCategory).Register(a.ctx, a.router, sMid)
}

package app

import (
	"github.com/giffone/forum-authentication/internal/adapters/api"
	authentication2 "github.com/giffone/forum-authentication/internal/adapters/api/authentication"
	"github.com/giffone/forum-authentication/internal/adapters/authentication"
	"github.com/giffone/forum-authentication/internal/service"
)

func (a *App) authentication(auth *authentication.Auth,
	srvUser service.User, sMid service.Middleware, aMid api.Middleware) {
	authentication2.NewHandler(auth, srvUser, sMid).Register(a.ctx, a.router, aMid)
}

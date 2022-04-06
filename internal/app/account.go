package app

import (
	"github.com/giffone/forum-authentication/internal/adapters/api"
	"github.com/giffone/forum-authentication/internal/adapters/api/account"
	"github.com/giffone/forum-authentication/internal/service"
)

func (a *App) account(srvPost service.Post, srvCategory service.Category,
	srvComment service.Comment, srvRatio service.Ratio, sMid api.Middleware) {
	account.NewHandler(srvPost, srvCategory, srvComment, srvRatio).Register(a.ctx, a.router, sMid)
}

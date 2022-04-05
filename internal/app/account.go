package app

import (
	"github.com/giffone/forum-authentication/internal/adapters/api"
	"github.com/giffone/forum-authentication/internal/adapters/api/account"
	"github.com/giffone/forum-authentication/internal/service"
)

func (a *App) account(srvPost service.Post, srvCategory service.Category,
	srvComment service.Comment, apiRatio api.Ratio, ses api.Session) {
	account.NewHandler(srvPost, srvCategory, srvComment, apiRatio).Register(a.ctx, a.router, ses)
}

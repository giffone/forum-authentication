package app

import (
	"github.com/giffone/forum-authentication/internal/adapters/api"
	post2 "github.com/giffone/forum-authentication/internal/adapters/api/post"
	"github.com/giffone/forum-authentication/internal/adapters/repository"
	"github.com/giffone/forum-authentication/internal/service"
	"github.com/giffone/forum-authentication/internal/service/post"
)

func (a *App) post(repo repository.Repo, srvCategory service.Category, srvComment service.Comment,
	srvRatio service.Ratio, sMid service.Middleware, apiMid api.Middleware) service.Post {
	srv := post.NewService(repo, srvCategory, srvRatio, sMid)
	post2.NewHandler(srv, srvCategory, srvComment, srvRatio).Register(a.ctx, a.router, apiMid)
	return srv
}

package app

import (
	"github.com/giffone/forum-authentication/internal/adapters/api"
	post2 "github.com/giffone/forum-authentication/internal/adapters/api/post"
	"github.com/giffone/forum-authentication/internal/adapters/repository"
	"github.com/giffone/forum-authentication/internal/service"
	"github.com/giffone/forum-authentication/internal/service/post"
)

func (a *App) post(repo repository.Repo, srvCategory service.Category,
	srvRatio service.Ratio) service.Post {
	return post.NewService(repo, srvRatio, srvCategory)
}

func (a *App) post2(srv service.Post, srvCategory service.Category,
	srvComment service.Comment, apiRatio api.Ratio, srvMid service.Middleware, apiMid api.Middleware) {
	post2.NewHandler(srv, srvCategory, srvComment, srvMid, apiRatio).Register(a.ctx, a.router, apiMid)

}

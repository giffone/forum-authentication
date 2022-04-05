package app

import (
	"github.com/giffone/forum-authentication/internal/adapters/api"
	ratio2 "github.com/giffone/forum-authentication/internal/adapters/api/ratio"
	"github.com/giffone/forum-authentication/internal/adapters/repository"
	"github.com/giffone/forum-authentication/internal/service"
	"github.com/giffone/forum-authentication/internal/service/ratio"
)

func (a *App) ratio(repo repository.Repo) (service.Ratio, api.Ratio) {
	srvRatio := ratio.NewService(repo)
	return srvRatio, ratio2.NewRatio(srvRatio)
}

func (a *App) ratio2(srv service.Ratio, srvPost service.Post, srvComment service.Comment) {
	srv.AddService(srvPost, srvComment)
}

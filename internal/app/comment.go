package app

import (
	"github.com/giffone/forum-authentication/internal/adapters/repository"
	"github.com/giffone/forum-authentication/internal/service"
	"github.com/giffone/forum-authentication/internal/service/comment"
)

func (a *App) comment(repo repository.Repo, srvPost service.Post, srvRatio service.Ratio) service.Comment {
	srv := comment.NewService(repo, srvPost, srvRatio)
	return srv
}

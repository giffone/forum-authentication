package app

import (
	"github.com/giffone/forum-authentication/internal/adapters/repository"
	"github.com/giffone/forum-authentication/internal/service"
	"github.com/giffone/forum-authentication/internal/service/category"
)

func (a *App) category(repo repository.Repo) service.Category {
	return category.NewService(repo)
}

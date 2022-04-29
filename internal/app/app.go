package app

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/giffone/forum-authentication/internal/adapters/authentication"
	"github.com/giffone/forum-authentication/internal/adapters/repository"
	"github.com/giffone/forum-authentication/internal/constant"
	"net/http"
)

type App struct {
	router *http.ServeMux
	ctx    context.Context
}

func NewApp(ctx context.Context) *App {
	return &App{
		ctx:    ctx,
		router: http.NewServeMux(),
	}
}

func (a *App) Run(driver string) (*sql.DB, *http.ServeMux, string) {
	repo := switcher(a.ctx, driver)
	db, port, _ := repo.ExportSettings()

	home := fmt.Sprintf("%s%s", constant.HomePage, port)
	tokens := authentication.NewSocialToken(home)

	srvMiddleware, apiMiddleware := a.middlewareService(repo)
	srvUser := a.user(repo, apiMiddleware, tokens)
	a.authentication(tokens, srvUser, srvMiddleware, apiMiddleware)
	srvCategory := a.category(repo)
	srvRatio := a.ratio(repo, srvMiddleware)
	srvComment := a.comment(repo, srvRatio, srvMiddleware)
	srvPost := a.post(repo, srvCategory, srvComment, srvRatio, srvMiddleware, apiMiddleware)
	a.home(srvPost, srvCategory, apiMiddleware)
	a.account(srvPost, srvCategory, srvComment, srvRatio, apiMiddleware)

	dir := http.Dir("internal/web/assets")
	dirHandler := http.StripPrefix("/assets/", http.FileServer(dir))
	a.router.Handle("/assets/", dirHandler)

	// FOR TEST ONLY
	_, _, schema := repo.ExportSettings()
	repository.NewLoremIpsum().Run(db, schema)
	// FOR TEST ONLY

	return db, a.router, port
}

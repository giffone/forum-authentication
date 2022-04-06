package app

import (
	"context"
	"database/sql"
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

	srvCategory := a.category(repo)
	srvMiddleware, apiMiddleware := a.middlewareService(repo)
	a.user(repo, apiMiddleware)
	srvRatio, apiRatio := a.ratio(repo)
	srvPost := a.post(repo, srvCategory, srvRatio)
	srvComment := a.comment(repo, srvPost, srvRatio)
	a.post2(srvPost, srvCategory, srvComment, apiRatio, srvMiddleware, apiMiddleware)
	a.home(srvPost, srvCategory, apiMiddleware)
	a.account(srvPost, srvCategory, srvComment, apiRatio, apiMiddleware)
	a.ratio2(srvRatio, srvPost, srvComment)

	dir := http.Dir("internal/web/assets")
	dirHandler := http.StripPrefix("/assets/", http.FileServer(dir))
	a.router.Handle("/assets/", dirHandler)

	// FOR TEST ONLY
	// _, _, schema := repo.ExportSettings()
	// repository.NewLoremIpsum().Run(db, schema)
	// FOR TEST ONLY

	return db, a.router, port
}

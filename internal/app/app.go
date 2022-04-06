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

	srvMiddleware, apiMiddleware := a.middlewareService(repo)
	a.user(repo, apiMiddleware)
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
	//_, _, schema := repo.ExportSettings()
	//repository.NewLoremIpsum().Run(db, schema)
	// FOR TEST ONLY

	return db, a.router, port
}

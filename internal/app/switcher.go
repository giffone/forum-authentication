package app

import (
	"context"
	"github.com/giffone/forum-authentication/internal/adapters/repository"
	"github.com/giffone/forum-authentication/internal/adapters/repository/mysqldb"
	"github.com/giffone/forum-authentication/internal/adapters/repository/sqlitedb"
	"log"
)

func switcher(ctx context.Context, driver string) repository.Repo {
	switch driver {
	case "mysql":
		return repository.NewRepoTx(ctx, &mysqldb.MySql{})
	case "sqlite3":
		return repository.NewRepo(ctx, &sqlitedb.Lite{})
	default:
		log.Fatalf("switcher: unknow driver %s\n", driver)
	}
	return nil
}

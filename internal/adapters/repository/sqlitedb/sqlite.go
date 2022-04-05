package sqlitedb

import (
	"context"
	"database/sql"
	"fmt"
	"forum/internal/adapters/repository"
	"forum/internal/adapters/repository/sqlitedb/schema"
	"forum/internal/object"
	"forum/pkg/paths"
	"log"
	"os"
)

type Lite struct {
	db *sql.DB
	c  *repository.Configuration
	q  *object.Query
}

func (l *Lite) Connect() *repository.Configuration {
	l.c = &repository.Configuration{
		Name:   "database-sqlite3.db",
		Path:   "db/database-sqlite3.db",
		PathB:  "db/backup/database-sqlite3.db",
		Driver: "sqlite3",
		Port:   ":3306",
	}
	return l.c
}

func (l *Lite) Query() *object.Query {
	l.q = &object.Query{
		Schema: schema.Query(),
	}
	return l.q
}

func (l *Lite) DataBase(ctx context.Context) *sql.DB {
	paths.CreatePaths("db/", "db/backup")
	if paths.NotExist(l.c.Path) { // if main base not exist
		l.create() // create new base

		if l.BackupExist() { // have backup
			l.backup()
		}
	}
	l.open()
	l.make() // make tables without restore from backup
	repository.InsertSrc(l.db, l.q.Schema)
	return l.db
}

func (l *Lite) create() {
	log.Printf("Creating %s...\n", l.c.Name)

	file, err := os.Create(l.c.Path)
	if err != nil {
		log.Fatalf("create db: can not create base - %v\n", err)
	}
	err = file.Close()
	if err != nil {
		return
	}

	log.Printf("%s created\n", l.c.Name)
}

func (l *Lite) open() {
	var err error

	l.db, err = sql.Open(l.c.Driver, l.c.Path)
	if err != nil {
		log.Fatalf("base: open: %v\n", err)
	}
}

func (l *Lite) make() {
	for _, table := range repository.MakeTables() {
		l.tables(table)
	}
}

func (l *Lite) tables(table string) {
	value, ok := l.q.Schema[table]
	if !ok {
		log.Fatalf("%s was not created, no such query to make a table. Fatal exit\n", table)
	}

	val := fmt.Sprintf(value, table)
	_, err := l.db.Exec(val)
	if err != nil {
		log.Fatalf("tables: execute %s: %v. Fatal exit\n", table, err) // if can not build table in Db, stop program
	}
}

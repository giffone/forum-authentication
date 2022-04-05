package mysqldb

import (
	"context"
	"database/sql"
	"fmt"
	"forum/internal/adapters/repository"
	"forum/internal/adapters/repository/sqlitedb/schema"
	"forum/internal/object"
	"log"
)

type MySql struct {
	Db *sql.DB
	c  *repository.Configuration
	q  *object.Query
}

func (ms *MySql) Connect() *repository.Configuration {
	ms.c = &repository.Configuration{
		Name:       "database-mysql.db",
		Path:       "db/database-mysql.db",
		PathB:      "db/backup/database-mysql.db",
		Driver:     "mysql",
		Port:       ":3306",
		Connection: "admin:admin@tcp(localhost:3306)/forum_db", //<username>:<pw>@tcp(<HOST>:<port>)/<dbname>
	}
	return ms.c
}

func (ms *MySql) Query() *object.Query {
	ms.q = &object.Query{
		Schema: schema.Query(), // can use query from sqlite
	}
	return ms.q
}

func (ms *MySql) DataBase(ctx context.Context) *sql.DB {
	ms.open(ctx)
	ms.make(ctx) // make tables without restore from backup
	repository.InsertSrc(ms.Db, ms.q.Schema)
	return ms.Db
}

func (ms *MySql) open(ctx context.Context) {
	var err error
	ms.Db, err = sql.Open(ms.c.Driver, ms.c.Connection)
	if err != nil {
		log.Fatalf("base: open: %v\n", err)
	}

	err = ms.Db.PingContext(ctx)
	if err != nil {
		log.Fatalf("base: ping: %v\n", err)
	}
}

func (ms *MySql) make(ctx context.Context) {
	tx, err := ms.Db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatalf("function begin tx: %v", err)
	}

	for _, table := range repository.MakeTables() {
		ms.tables(tx, ctx, table)
	}

	if err = tx.Commit(); err != nil {
		log.Fatal(err)
	}
}

func (ms *MySql) tables(tx *sql.Tx, ctx context.Context, table string) {
	value, ok := ms.q.Schema[table]
	if !ok {
		log.Fatalf("%s was not created, no such query to make a table. Fatal exit\n", table)
	}

	val := fmt.Sprintf(value, table)

	_, err := tx.ExecContext(ctx, val)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Fatalf("tables: execute %s: unable to rollback: %v", table, rollbackErr)
		}
		log.Fatalf("tables: execute %s: %v. Fatal exit\n", table, err)
	}
}

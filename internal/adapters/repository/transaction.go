package repository

import (
	"context"
	"database/sql"
	"forum/internal/constant"
	"forum/internal/object"
	"log"
)

func TxBegin(ctx context.Context, db *sql.DB) (*sql.Tx, object.Status) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, object.StatusByCodeAndLog(constant.Code500,
			err, "transaction begin:")
	}
	return tx, nil
}

func TxRollBack(tx *sql.Tx) {
	if err := tx.Rollback(); err != nil {
		log.Printf("transaction: can not rollback: %v", err)
	}
}

func TxCommit(tx *sql.Tx) {
	if err := tx.Commit(); err != nil {
		log.Printf("transaction commit: %v", err)
	}
}

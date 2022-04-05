package repository

import (
	"context"
	"database/sql"
	"fmt"
	"forum/internal/constant"
	"forum/internal/object"
)

func Query(ctx context.Context, db *sql.DB, query string, keys []interface{}) (*sql.Rows, object.Status) {
	rows, err := db.QueryContext(ctx, query, keys...)
	if err != nil {
		mess := fmt.Sprintf("query:\n%s", query)
		return nil, object.StatusByCodeAndLog(constant.Code500, err, mess)
	}
	return rows, nil
}

func QueryTx(ctx context.Context, tx *sql.Tx, query string, keys []interface{}) (*sql.Rows, object.Status) {
	rows, err := tx.QueryContext(ctx, query, keys...)
	if err != nil {
		mess := fmt.Sprintf("queryTx:\n%s", query)
		return nil, object.StatusByCodeAndLog(constant.Code500, err, mess)
	}
	return rows, nil
}

func QueryRow(ctx context.Context, db *sql.DB, query string, keys []interface{}) (*sql.Row, object.Status) {
	row := db.QueryRowContext(ctx, query, keys...)
	if err := row.Err(); err != nil {
		mess := fmt.Sprintf("queryRow:\n%s", query)
		return nil, object.StatusByCodeAndLog(constant.Code500, err, mess)
	}
	return row, nil
}

func QueryRowTx(ctx context.Context, tx *sql.Tx, query string, keys []interface{}) (*sql.Row, object.Status) {
	row := tx.QueryRowContext(ctx, query, keys...)
	if err := row.Err(); err != nil {
		mess := fmt.Sprintf("queryRowTx:\n%s", query)
		return nil, object.StatusByCodeAndLog(constant.Code500, err, mess)
	}
	return row, nil
}

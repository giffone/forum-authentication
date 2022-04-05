package repository

import (
	"database/sql"
	"forum/internal/constant"
	"forum/internal/object"
	"forum/internal/object/model"
	"log"
)

func Rows(rows *sql.Rows, m model.Models) object.Status {
	for rows.Next() {
		keys := m.NewList()
		if err := rows.Scan(keys...); err != nil {
			if err == sql.ErrNoRows {
				break
			}
			return object.StatusByCodeAndLog(constant.Code500, err, "rows:")
		}
	}
	if err := rows.Err(); err != nil {
		return object.StatusByCodeAndLog(constant.Code500, err, "rows: end with:")
	}
	return nil
}

func Row(row *sql.Row, m model.Model) object.Status {
	keys := m.New()
	if err := row.Scan(keys...); err != nil {
		if err != nil {
			if err == sql.ErrNoRows {
				return nil
			} else {
				return object.StatusByCodeAndLog(constant.Code500, err, "row:")
			}
		}
	}
	return nil
}

func CloseRows(rows *sql.Rows) {
	if err := rows.Close(); err != nil {
		log.Printf("close rows: %v", err)
	}
}

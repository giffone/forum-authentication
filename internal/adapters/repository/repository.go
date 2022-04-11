package repository

import (
	"context"
	"database/sql"
	"github.com/giffone/forum-authentication/internal/constant"
	"github.com/giffone/forum-authentication/internal/object"
	"github.com/giffone/forum-authentication/internal/object/dto"
	"github.com/giffone/forum-authentication/internal/object/model"
)

type repo struct {
	conf   *Configuration
	schema *object.Query
	db     *sql.DB
}

func NewRepo(ctx context.Context, new New) Repo {
	return &repo{
		conf:   new.Connect(),
		schema: new.Query(),
		db:     new.DataBase(ctx),
	}
}

func (r *repo) Create(ctx context.Context, d dto.DTO) (int, object.Status) {
	// prepare query for db
	q := d.Create().MakeQuery(r.schema)

	stmt, err := r.db.PrepareContext(ctx, q.Query)
	if err != nil {
		return 0,
			object.StatusByCodeAndLog(constant.Code500,
				err, "create: stmt")
	}
	// apply query
	res, err := stmt.ExecContext(ctx, q.Fields...)
	if err != nil {
		return 0,
			object.StatusByCodeAndLog(constant.Code500,
				err, "create: exec")
	}

	// get id of new record
	id, err := res.LastInsertId()
	if err != nil {
		return 0,
			object.StatusByCodeAndLog(constant.Code500,
				err, "create: last inserted id")
	}
	err = stmt.Close()
	if err != nil {
		return 0,
			object.StatusByCodeAndLog(constant.Code500,
				err, "create: close stmt")
	}
	return int(id), nil
}

func (r *repo) Delete(ctx context.Context, d dto.DTO) object.Status {
	// prepare query for db
	q := d.Delete().MakeQuery(r.schema)
	stmt, err := r.db.PrepareContext(ctx, q.Query)
	if err != nil {
		return object.StatusByCodeAndLog(constant.Code500,
			err, "delete: stmt")
	}
	_, err = stmt.ExecContext(ctx, q.Fields...)
	if err != nil {
		return object.StatusByCodeAndLog(constant.Code500,
			err, "delete")
	}
	err = stmt.Close()
	if err != nil {
		return object.StatusByCodeAndLog(constant.Code500,
			err, "delete: close stmt")
	}
	return nil
}

func (r *repo) GetList(ctx context.Context, m model.Models) object.Status {
	// prepare query for db
	q := m.GetList().MakeQuery(r.schema)
	rows, sts := Query(ctx, r.db, q.Query, q.Fields)
	if sts != nil {
		return sts // 500
	}
	defer CloseRows(rows)
	//panic("stop")
	sts = Rows(rows, m)
	if sts != nil {
		return sts
	}
	return nil
}

func (r *repo) GetOne(ctx context.Context, m model.Model) object.Status {
	// prepare query for db
	q := m.Get().MakeQuery(r.schema)

	row, sts := QueryRow(ctx, r.db, q.Query, q.Fields)
	if sts != nil {
		return sts
	}
	sts = Row(row, m)
	if sts != nil {
		return sts
	}
	return nil
}

func (r *repo) ExportSettings() (*sql.DB, string, *object.Query) {
	return r.db, r.conf.Port, r.schema
}

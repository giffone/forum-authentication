package repository

import (
	"context"
	"database/sql"
	"forum/internal/constant"
	"forum/internal/object"
	"forum/internal/object/dto"
	"forum/internal/object/model"
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

	// apply query
	res, err := r.db.ExecContext(ctx, q.Query, q.Fields...)
	if err != nil {
		return 0,
			object.StatusByCodeAndLog(constant.Code500,
				err, "create")
	}

	// get id of new record
	id, err := res.LastInsertId()
	if err != nil {
		return 0,
			object.StatusByCodeAndLog(constant.Code500,
				err, "create: last inserted id")
	}
	return int(id), nil
}

func (r *repo) Delete(ctx context.Context, d dto.DTO) object.Status {
	// prepare query for db
	q := d.Delete().MakeQuery(r.schema)
	_, err := r.db.ExecContext(ctx, q.Query, q.Fields...)
	if err != nil {
		return object.StatusByCodeAndLog(constant.Code500,
			err, "delete")
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

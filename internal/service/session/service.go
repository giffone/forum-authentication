package session

import (
	"context"
	"github.com/giffone/forum-authentication/internal/adapters/repository"
	"github.com/giffone/forum-authentication/internal/object"
	"github.com/giffone/forum-authentication/internal/object/dto"
	"github.com/giffone/forum-authentication/internal/object/model"
	"github.com/giffone/forum-authentication/internal/service"
	"github.com/giffone/forum-authentication/pkg/datef"
	"log"
)

type sSession struct {
	repo repository.Repo
}

func NewService(repo repository.Repo) service.Session {
	return &sSession{
		repo: repo,
	}
}

func (ss *sSession) Create(ctx context.Context, d *dto.Session) (int, object.Status) {
	// if session already exist, delete it
	dDelete := dto.NewSession(nil, nil, d.Obj.Ck)
	if sts := ss.repo.Delete(ctx, dDelete); sts != nil {
		return 0, sts
	}
	// create session
	return ss.repo.Create(ctx, d)
}

func (ss *sSession) Check(ctx context.Context, d *dto.Session) (interface{}, object.Status) {
	// make new model session
	session := model.NewSession(nil, d.Obj.Ck)
	// get session from db
	sts := ss.repo.GetOne(ctx, session)
	if sts != nil {
		return nil, sts
	}
	// session not match
	if session.UUID != d.Obj.Ck.SessionUUID {
		log.Printf("uuid did not match db: %s dto: %v", session.UUID, d.Obj.Ck)
		return nil, nil
	}
	// session expire
	if datef.Expire(session.Expire) {
		sts = ss.repo.Delete(ctx, d)
	}
	return session, sts
}

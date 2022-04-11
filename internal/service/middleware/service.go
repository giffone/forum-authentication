package middleware

import (
	"context"
	"github.com/giffone/forum-authentication/internal/adapters/repository"
	"github.com/giffone/forum-authentication/internal/constant"
	"github.com/giffone/forum-authentication/internal/object"
	"github.com/giffone/forum-authentication/internal/object/dto"
	"github.com/giffone/forum-authentication/internal/object/model"
	"github.com/giffone/forum-authentication/internal/service"
	"github.com/giffone/forum-authentication/pkg/datef"
	"log"
	"strconv"
)

type sMiddleware struct {
	repo repository.Repo
}

func NewService(repo repository.Repo) service.Middleware {
	return &sMiddleware{
		repo: repo,
	}
}

func (smw *sMiddleware) CreateSession(ctx context.Context, d *dto.Session) (int, object.Status) {
	// if middleware already exist, delete it
	dDelete := dto.NewSession(nil, nil, d.Obj.Ck)
	if sts := smw.repo.Delete(ctx, dDelete); sts != nil {
		return 0, sts
	}
	// create middleware
	return smw.repo.Create(ctx, d)
}

func (smw *sMiddleware) CheckSession(ctx context.Context, d *dto.Session) (interface{}, object.Status) {
	// make new model middleware
	session := model.NewSession(nil, d.Obj.Ck)
	// get middleware from db
	sts := smw.repo.GetOne(ctx, session)
	if sts != nil {
		return nil, sts
	}
	// middleware not match
	if session.UUID != d.Obj.Ck.SessionUUID {
		log.Printf("uuid did not match db: %s dto: %v", session.UUID, d.Obj.Ck)
		return nil, nil
	}
	// middleware expire
	if datef.Expire(session.Expire) {
		sts = smw.repo.Delete(ctx, d)
	}
	return session, sts
}

func (smw *sMiddleware) GetID(ctx context.Context, d *dto.CheckID) (int, object.Status) {
	var value interface{}
	if d.Atoi {
		if d.IDString == "" {
			return 0, object.StatusByCodeAndLog(constant.Code500,
				nil, "check id: atoi = true, but IDString empty")
		}
		idInt, err := strconv.Atoi(d.IDString)
		if err != nil || idInt == 0 {
			return 0, object.StatusByCodeAndLog(constant.Code400,
				err, "check id: atoi")
		}
		value = idInt
	} else {
		value = d.ID
	}
	who := model.NewCheckID(nil, nil, nil)
	who.MakeKeys(d.Who, value)
	sts := smw.repo.GetOne(ctx, who)
	if sts != nil {
		return 0, sts
	}
	if who.ID == 0 {
		return 0, object.StatusByCode(constant.Code400)
	}
	return who.ID, nil
}

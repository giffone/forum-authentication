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
	repo      repository.Repo
	sPost     service.Post
	sCategory service.Category
}

func NewService(repo repository.Repo, sPost service.Post, sCategory service.Category) service.Middleware {
	return &sMiddleware{
		repo:      repo,
		sPost:     sPost,
		sCategory: sCategory,
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

func (smw *sMiddleware) Check(ctx context.Context, d *dto.CheckID) ([]int, object.Status) {
	var ids []int
	for i := 0; i < len(d.ID); i++ {
		id, err := strconv.Atoi(d.ID[i])
		if err != nil {
			return nil, object.StatusByCodeAndLog(constant.Code500,
				err, "check id: atoi")
		}
		who := model.NewCheckID(nil, nil, nil)
		who.MakeKeys(d.Who, id)

		sts := smw.repo.GetList(ctx, who)
		if sts != nil {
			return nil, sts
		}
		if len(who.Slice) == 0 {
			return nil, object.StatusByCode(constant.Code400)
		}
		ids = append(ids, id)
	}
	return ids, nil
	// be careful send []int=nil without sts with error, will panic!
}

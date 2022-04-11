package user

import (
	"context"
	"github.com/giffone/forum-authentication/internal/adapters/repository"
	"github.com/giffone/forum-authentication/internal/constant"
	"github.com/giffone/forum-authentication/internal/object"
	"github.com/giffone/forum-authentication/internal/object/dto"
	"github.com/giffone/forum-authentication/internal/object/model"
	"github.com/giffone/forum-authentication/internal/service"
	"golang.org/x/crypto/bcrypt"
)

type sUser struct {
	repo repository.Repo
}

func NewService(repo repository.Repo) service.User {
	return &sUser{repo: repo}
}

func (su *sUser) Create(ctx context.Context, d *dto.User) (int, object.Status) {
	id, sts := su.repo.Create(ctx, d)
	if sts != nil {
		return 0, object.StatusByText(nil, constant.AlreadyExist, "login or password")
	}
	return id, nil
}

func (su *sUser) CheckLoginPassword(ctx context.Context, d *dto.User) (int, object.Status) {
	m := model.NewUser(nil, nil)
	m.MakeKeys(constant.FieldLogin, d.Login)
	sts := su.repo.GetOne(ctx, m)
	if sts != nil {
		return 0, sts
	}
	if m.ID == 0 { // if did not find login
		return 0, object.StatusByText(nil, constant.WrongEnter,
			"login did not founded or password")
	}
	err := bcrypt.CompareHashAndPassword([]byte(m.Password), []byte(d.Password))
	if err != nil { // passwords did not match
		return 0, object.StatusByText(err, constant.WrongEnter, "login or password")
	}
	return m.ID, nil
}

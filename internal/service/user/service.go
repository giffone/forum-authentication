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
	"strconv"
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
		return 0, object.StatusByText(constant.AlreadyExist, "login or password", nil)
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
		return 0, object.StatusByText(constant.WrongEnter,
			"login did not founded or password", nil)
	}
	err := bcrypt.CompareHashAndPassword([]byte(m.Password), []byte(d.Password))
	if err != nil { // passwords did not match
		return 0, object.StatusByText(constant.WrongEnter, "login or password", err)
	}
	return m.ID, nil
}

func (su *sUser) Check(ctx context.Context, slice []string) ([]int, object.Status) {
	var idUser []int
	for i := 0; i < len(slice); i++ {
		id, err := strconv.Atoi(slice[i])
		if err != nil {
			return nil, object.StatusByCodeAndLog(constant.Code500,
				err, "check user: atoi")
		}
		posts := model.NewPosts(nil, nil)
		posts.MakeKeys(constant.KeyID, id)

		sts := su.repo.GetList(ctx, posts)
		if sts != nil {
			return nil, sts
		}
		if len(posts.Slice) == 0 {
			return nil, object.StatusByCode(constant.Code400)
		}
		idUser = append(idUser, id)
	}
	return idUser, nil
}

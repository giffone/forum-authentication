package service

import (
	"context"
	"github.com/giffone/forum-authentication/internal/object"
	"github.com/giffone/forum-authentication/internal/object/dto"
	"github.com/giffone/forum-authentication/internal/object/model"
)

type Post interface { // use cases
	Create(ctx context.Context, dto *dto.Post) (int, object.Status)
	Get(ctx context.Context, m model.Models) (interface{}, object.Status)
	GetChan(ctx context.Context, m model.Models) (interface{}, object.Status)
}

type User interface { // use cases
	Create(ctx context.Context, dto *dto.User) (int, object.Status)
	CheckLoginPassword(ctx context.Context, dto *dto.User) (int, object.Status)
	Check(ctx context.Context, slice []string) ([]int, object.Status)
}

type Category interface { // use cases
	Create(ctx context.Context, dto *dto.Category) object.Status
	Delete(ctx context.Context, id int) object.Status
	GetList(ctx context.Context, m model.Models) (interface{}, object.Status)
	GetFor(ctx context.Context, pc model.PostOrComment) object.Status
	GetForChan(ctx context.Context, pc model.PostOrComment, channel chan object.Status)
}

type Ratio interface { // use cases
	Create(ctx context.Context, d *dto.Ratio) (int, object.Status)
	CountFor(ctx context.Context, pc model.PostOrComment) object.Status
	Liked(ctx context.Context, pc model.PostOrComment) object.Status
	CountForChan(ctx context.Context, pc model.PostOrComment, channel chan object.Status)
	LikedChan(ctx context.Context, pc model.PostOrComment, channel chan object.Status)
}

type Middleware interface { // use cases
	CreateSession(ctx context.Context, dto *dto.Session) (int, object.Status)
	CheckSession(ctx context.Context, dto *dto.Session) (interface{}, object.Status)
	Check(ctx context.Context, d *dto.CheckID) ([]int, object.Status)
}

type Comment interface {
	Create(ctx context.Context, dto *dto.Comment) (int, object.Status)
	Delete(ctx context.Context, id int) object.Status
	Get(ctx context.Context, m model.Models) (interface{}, object.Status)
	GetChan(ctx context.Context, m model.Models) (interface{}, object.Status)
	Check(ctx context.Context, slice []string) ([]int, object.Status)
}

package ratio

import (
	"context"
	"github.com/giffone/forum-authentication/internal/adapters/api"
	"github.com/giffone/forum-authentication/internal/object"
	"github.com/giffone/forum-authentication/internal/object/dto"
	"github.com/giffone/forum-authentication/internal/service"
	"net/http"
)

type ratio struct {
	service service.Ratio
}

func NewRatio(service service.Ratio) api.Ratio {
	return &ratio{
		service: service,
	}
}

func (rt *ratio) Rate(ctx context.Context, ck *object.Cookie, r *http.Request) object.Status {
	// create DTO with a new rate
	like := dto.NewRatio(nil, nil, ck)
	// add request data to DTO and check err
	if r.Method == "POST" {
		if !like.AddByPOST(r) {
			return like.Obj.Sts
		}
	} else if r.Method == "GET" {
		if !like.AddByGET(r) {
			return like.Obj.Sts
		}
	}
	// create like in database
	_, sts := rt.service.Create(ctx, like)
	if sts != nil {
		return sts
	}
	return nil
}

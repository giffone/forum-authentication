package post

import (
	"context"
	"github.com/giffone/forum-authentication/internal/adapters/repository"
	"github.com/giffone/forum-authentication/internal/constant"
	"github.com/giffone/forum-authentication/internal/object"
	"github.com/giffone/forum-authentication/internal/object/dto"
	"github.com/giffone/forum-authentication/internal/object/model"
	"github.com/giffone/forum-authentication/internal/service"
	"log"
)

type sPost struct {
	repo        repository.Repo
	sLike       service.Ratio
	sCategory   service.Category
	sMiddleware service.Middleware
}

func NewService(repo repository.Repo, sLike service.Ratio,
	sCategory service.Category, sMiddleware service.Middleware) service.Post {
	return &sPost{
		repo:        repo,
		sLike:       sLike,
		sCategory:   sCategory,
		sMiddleware: sMiddleware,
	}
}

func (sp *sPost) Create(ctx context.Context, d *dto.Post) (int, object.Status) {
	ctx2, cancel := context.WithTimeout(ctx, constant.TimeLimitDB)
	defer cancel()

	// check valid categories
	if len(d.Categories.Slice) > 0 {
		categories := dto.NewCheckID(constant.KeyCategory, d.Categories.Slice)
		ids, sts := sp.sMiddleware.Check(ctx, categories)
		if sts != nil {
			return 0, sts
		}
		if ids != nil {
			d.Categories.ID = ids
		}
	}
	// create post
	id, sts := sp.repo.Create(ctx2, d)
	if sts != nil {
		return 0, sts
	}
	// remember id new created post
	d.Categories.Post = id
	if lSlice := len(d.Categories.ID); lSlice != 0 {
		for i := 0; i < lSlice; i++ {
			// current id category to add
			d.Categories.Category = d.Categories.ID[i]
			// create category
			_, sts = sp.repo.Create(ctx2, d.Categories)
			if sts != nil {
				return 0, sts
			}
		}
	}
	return id, nil
}

func (sp *sPost) Delete(ctx context.Context, id int) *object.Statuses {
	return nil
}

func (sp *sPost) Get(ctx context.Context, m model.Models) (interface{}, object.Status) {
	ctx2, cancel := context.WithTimeout(ctx, constant.TimeLimitDB)
	defer cancel()

	sts := sp.repo.GetList(ctx2, m)
	if sts != nil {
		return nil, sts
	}

	posts := m.Return().Posts

	lSlice := len(posts.Slice)
	if lSlice == 0 {
		return posts.IfNil(), nil
	}

	// checks if authorized user liked post
	if posts.Ck.Session {
		// checks liked post or not
		sts = sp.sLike.Liked(ctx, posts)
		if sts != nil {
			return nil, sts
		}
	}
	// checks categories for post
	sts = sp.sCategory.GetFor(ctx, posts)
	if sts != nil {
		return nil, sts
	}
	// count likes/dislikes
	sts = sp.sLike.CountFor(ctx, posts)
	if sts != nil {
		return nil, sts
	}
	return posts.Slice, nil
}

func (sp *sPost) GetChan(ctx context.Context, m model.Models) (interface{}, object.Status) {
	ctx2, cancel := context.WithTimeout(ctx, constant.TimeLimitDB)
	defer cancel()

	sts := sp.repo.GetList(ctx2, m)
	if sts != nil {
		return nil, sts
	}

	posts := m.Return().Posts

	lSlice := len(posts.Slice)
	if lSlice == 0 {
		return posts.IfNil(), nil
	}

	channel := make(chan object.Status)
	// checks if authorized user liked post
	if posts.Ck.Session {
		// checks liked post or not
		go sp.sLike.LikedChan(ctx, posts, channel)
	} else {
		channel <- nil
	}
	// checks categories for post
	go sp.sCategory.GetForChan(ctx, posts, channel)
	// count likes/dislikes
	go sp.sLike.CountForChan(ctx, posts, channel)

	err1 := <-channel
	err2 := <-channel
	err3 := <-channel

	if err1 != nil || err2 != nil || err3 != nil {
		log.Printf("err1: %v\nerr2: %v\nerr3: %v\n", err1, err2, err3)
		return nil, sts
	}
	return posts.Slice, nil
}

package comment

import (
	"context"
	"github.com/giffone/forum-authentication/internal/adapters/repository"
	"github.com/giffone/forum-authentication/internal/constant"
	"github.com/giffone/forum-authentication/internal/object"
	"github.com/giffone/forum-authentication/internal/object/dto"
	"github.com/giffone/forum-authentication/internal/object/model"
	"github.com/giffone/forum-authentication/internal/service"
)

type sComment struct {
	repo        repository.Repo
	sRatio      service.Ratio
	sMiddleware service.Middleware
}

func NewService(repo repository.Repo, sRatio service.Ratio,
	sMiddleware service.Middleware) service.Comment {
	return &sComment{
		repo:        repo,
		sRatio:      sRatio,
		sMiddleware: sMiddleware,
	}
}

func (sc *sComment) Create(ctx context.Context, d *dto.Comment) (int, object.Status) {
	ctx2, cancel := context.WithTimeout(ctx, constant.TimeLimitDB)
	defer cancel()

	// check valid postID
	postID := dto.NewCheckID(constant.KeyPost, []string{d.Obj.Ck.PostString})
	ids, sts := sc.sMiddleware.CheckID(ctx, postID)
	if sts != nil {
		return 0, sts
	}
	if ids != nil {
		d.Obj.Ck.Post = ids[0]
	} else {
		return 0, object.StatusByCode(constant.Code400)
	}
	// create comment
	id, sts := sc.repo.Create(ctx2, d)
	if sts != nil {
		return 0, sts
	}
	return id, nil
}

func (sc *sComment) Delete(ctx context.Context, id int) object.Status {
	return nil
}

func (sc *sComment) Get(ctx context.Context, m model.Models) (interface{}, object.Status) {
	ctx2, cancel := context.WithTimeout(ctx, constant.TimeLimitDB)
	defer cancel()

	sts := sc.repo.GetList(ctx2, m)
	if sts != nil {
		return nil, sts
	}

	comments := m.Return().Comments
	lSlice := len(comments.Slice)
	if lSlice == 0 {
		return comments.IfNil(), nil
	}
	// checks if authorized user liked comment
	if comments.Ck.Session {
		// checks liked comment or not
		sts = sc.sRatio.Liked(ctx2, comments)
		if sts != nil {
			return nil, sts
		}
	}
	// checks need refer to ... or not
	if comments.St.Refers {
		// make refer to own post
		sts = sc.refer(ctx2, comments)
		if sts != nil {
			return nil, sts
		}
	}
	// count likes/dislikes
	sts = sc.sRatio.CountFor(ctx2, comments)
	if sts != nil {
		return nil, sts
	}
	return comments.Slice, nil
}

func (sc *sComment) GetChan(ctx context.Context, m model.Models) (interface{}, object.Status) {
	ctx2, cancel := context.WithTimeout(ctx, constant.TimeLimitDB)
	defer cancel()

	sts := sc.repo.GetList(ctx2, m)
	if sts != nil {
		return nil, sts
	}

	comments := m.Return().Comments
	lSlice := len(comments.Slice)
	if lSlice == 0 {
		return comments.IfNil(), nil
	}

	channel := make(chan object.Status)
	// checks if authorized user liked comment
	if comments.Ck.Session {
		// checks liked comment or not
		go sc.sRatio.LikedChan(ctx2, comments, channel)
	} else {
		channel <- nil
	}
	// checks need refer to ... or not
	if comments.St.Refers {
		// make refer to own post
		go sc.referChan(ctx2, comments, channel)
	} else {
		channel <- nil
	}
	// count likes/dislikes
	go sc.sRatio.CountForChan(ctx2, comments, channel)

	err1 := <-channel
	err2 := <-channel
	err3 := <-channel

	if err1 != nil || err2 != nil || err3 != nil {
		return nil, sts
	}
	return comments.Slice, nil
}

func (sc *sComment) refer(ctx context.Context, c *model.Comments) object.Status {
	for i := 0; i < len(c.Slice); i++ {
		p := model.NewPost(nil, nil)
		p.MakeKeys(constant.KeyPost, c.Slice[i].Post)
		sts := sc.repo.GetOne(ctx, p)
		if sts != nil {
			return sts
		}
		c.Slice[i].Title = p
	}
	return nil
}

func (sc *sComment) referChan(ctx context.Context, c *model.Comments, channel chan object.Status) {
	for i := 0; i < len(c.Slice); i++ {
		p := model.NewPost(nil, nil)
		p.MakeKeys(constant.KeyPost, c.Slice[i].Post)
		sts := sc.repo.GetOne(ctx, p)
		if sts != nil {
			channel <- sts
			return
		}
		c.Slice[i].Title = p
	}
	channel <- nil
}
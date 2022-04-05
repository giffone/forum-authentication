package ratio

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

type sRatio struct {
	repo     repository.Repo
	sPost    service.Post
	sComment service.Comment
}

func NewService(repo repository.Repo) service.Ratio {
	return &sRatio{
		repo: repo,
	}
}

func (sr *sRatio) AddService(sPost service.Post, sComment service.Comment) {
	sr.sPost = sPost
	sr.sComment = sComment
}

func (sr *sRatio) Create(ctx context.Context, d *dto.Ratio) (int, object.Status) {
	ctx2, cancel := context.WithTimeout(ctx, constant.TimeLimitDB)
	defer cancel()

	if d.Obj.Ck.PostString != "" {
		// check postID valid (refer number)
		idPost, sts := sr.sPost.Check(ctx, []string{d.Obj.Ck.PostString})
		if sts != nil {
			return 0, sts
		}
		if len(idPost) == 0 {
			return 0, object.StatusByCode(constant.Code400)
		}
		d.Obj.Ck.Post = idPost[0]
	}

	like := model.NewLike(nil, d.Obj.Ck)
	post := false
	// post
	if id, ok := d.PostOrComm[constant.KeyPost]; ok {
		post = true
		// check id for valid
		idPost, sts := sr.sPost.Check(ctx2, []string{id})

		if sts != nil {
			return 0, sts
		}
		if len(idPost) == 0 {
			return 0, object.StatusByCode(constant.Code400)
		}
		// post - keys for get likes from db
		like.MakeKeys(constant.KeyPost, d.Obj.Ck.User, idPost[0])
		// post - keys for create like in db
		d.MakeKeys(constant.KeyPost, d.Obj.Ck.User, idPost[0])
		// comment
	} else if id, ok := d.PostOrComm[constant.KeyComment]; ok {
		// check id for valid
		idComm, sts := sr.sComment.Check(ctx2, []string{id})
		if sts != nil {
			return 0, sts
		}
		if len(idComm) == 0 {
			return 0, object.StatusByCode(constant.Code400)
		}
		// comment - keys for get likes from db
		like.MakeKeys(constant.KeyComment, d.Obj.Ck.User, idComm[0])
		// comment - keys for create like in db
		d.MakeKeys(constant.KeyComment, d.Obj.Ck.User, idComm[0])
	} else {
		// if not post or comment
		return 0, object.StatusByCode(constant.Code400)
	}
	//check like exist by user_id and post_id/comment_id
	sts := sr.repo.GetOne(ctx2, like)
	if sts != nil {
		return 0, sts
	}
	// not exist
	if like.PostOrComm == 0 {
		// keys - here need only key and ignore value [in DTO method create]
		id, sts := sr.repo.Create(ctx2, d)
		if sts != nil {
			return 0, sts
		}
		// post_id for return page (redirect)
		return id, nil
	}
	// DTO for delete object
	dDelete := dto.NewRatio(nil, nil, d.Obj.Ck)
	// make keys for delete by id
	if post {
		dDelete.MakeKeys(constant.KeyPost, like.PostOrComm)
	} else {
		dDelete.MakeKeys(constant.KeyComment, like.PostOrComm)
	}
	// delete
	sts = sr.repo.Delete(ctx2, dDelete)
	if sts != nil {
		return 0, sts
	}
	// is same - was like and new like (not create new)
	if like.Like == d.Like {
		// post_id for return page (redirect)
		return d.Obj.Ck.Post, nil
	}
	// is not same - create new
	id, sts := sr.repo.Create(ctx2, d)
	if sts != nil {
		return 0, sts
	}
	// post_id for return page (redirect)
	return id, nil
}

func (sr *sRatio) CountFor(ctx context.Context, pc model.PostOrComment) object.Status {
	for i := 0; i < pc.LSlice(); i++ {
		id := pc.PostOrCommentID(i)
		likesCount := model.NewLikesCount(pc.Settings().ClearKey(), pc.Cookie()) // auto insert session
		likesCount.MakeKeys(pc.KeyRole(), id)
		// for make map["post"]id
		likesCount.PostOrComm = id
		sts := sr.repo.GetList(ctx, likesCount)
		if sts != nil {
			return sts
		}
		lSlice := len(likesCount.Slice)
		if lSlice == 0 {
			pc.Add(constant.KeyLike, i, likesCount.IfNil())
		} else {
			// like or dislike only, need to show another with 0
			if lSlice == 1 {
				if likesCount.Slice[0].Body == constant.FieldLike {
					likesCount.Slice = append(likesCount.Slice, likesCount.DislikeNil())
				} else {
					likesCount.Slice = append(likesCount.Slice, likesCount.LikeNil())
				}
			}
			pc.Add(constant.KeyLike, i, likesCount.Slice)
		}
	}
	return nil
}

func (sr *sRatio) Liked(ctx context.Context, pc model.PostOrComment) object.Status {
	user := pc.Cookie().User
	for i := 0; i < pc.LSlice(); i++ {
		id := pc.PostOrCommentID(i)
		like := model.NewLike(nil, nil)
		like.MakeKeys(pc.KeyLiked(), user, id)
		sts := sr.repo.GetOne(ctx, like)
		if sts != nil {
			return sts
		}
		pc.Add(constant.KeyRated, i, like.Body)
	}
	return nil
}

func (sr *sRatio) CountForChan(ctx context.Context, pc model.PostOrComment, channel chan object.Status) {
	log.Println("in CountForChan")
	for i := 0; i < pc.LSlice(); i++ {
		id := pc.PostOrCommentID(i)
		likesCount := model.NewLikesCount(pc.Settings().ClearKey(), pc.Cookie()) // auto insert session
		likesCount.MakeKeys(pc.KeyRole(), id)
		// for make map["post"]id
		likesCount.PostOrComm = id
		sts := sr.repo.GetList(ctx, likesCount)
		if sts != nil {
			log.Println("err CountForChan")
			channel <- sts
			return
		}
		lSlice := len(likesCount.Slice)
		if lSlice == 0 {
			pc.Add(constant.KeyLike, i, likesCount.IfNil())
		} else {
			// like or dislike only, need to show another with 0
			if lSlice == 1 {
				if likesCount.Slice[0].Body == constant.FieldLike {
					likesCount.Slice = append(likesCount.Slice, likesCount.DislikeNil())
				} else {
					likesCount.Slice = append(likesCount.Slice, likesCount.LikeNil())
				}
			}
			pc.Add(constant.KeyLike, i, likesCount.Slice)
		}
	}
	log.Println("out CountForChan")
	channel <- nil
}

func (sr *sRatio) LikedChan(ctx context.Context, pc model.PostOrComment, channel chan object.Status) {
	log.Println("in LikedChan")
	user := pc.Cookie().User
	for i := 0; i < pc.LSlice(); i++ {
		id := pc.PostOrCommentID(i)
		like := model.NewLike(nil, nil)
		like.MakeKeys(pc.KeyLiked(), user, id)
		sts := sr.repo.GetOne(ctx, like)
		if sts != nil {
			log.Println("err LikedChan")
			channel <- sts
			return
		}
		pc.Add(constant.KeyRated, i, like.Body)
	}
	log.Println("out LikedChan")
	channel <- nil
}

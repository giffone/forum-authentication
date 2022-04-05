package model

import (
	"github.com/giffone/forum-authentication/internal/constant"
	"github.com/giffone/forum-authentication/internal/object"
)

type LCount struct {
	Like       int    // like_id or dislike_id
	Body       string // like or dislike (name)
	Count      int    // number of likes/dislikes
	Session    bool   // logged user or not
	AllPost    bool   // one post window or not
	PostOrComm int    // current post_id or comment_id
}

func (lk *LCount) ifNil() *LCount {
	return &LCount{
		Body: "not rated yet",
	}
}

type LikesCount struct {
	Slice      []*LCount
	PostOrComm int // current post_id or comment_id
	St         *object.Settings
	Ck         *object.Cookie
}

func NewLikesCount(st *object.Settings, ck *object.Cookie) *LikesCount {
	lk := new(LikesCount)
	if st == nil {
		lk.St = &object.Settings{
			Key: make(map[string][]interface{}),
		}
	} else {
		lk.St = st
	}
	if ck == nil {
		lk.Ck = new(object.Cookie)
	} else {
		lk.Ck = ck
	}
	return lk
}

func (lk *LikesCount) MakeKeys(key string, data ...interface{}) {
	if key != "" {
		lk.St.Key[key] = data
	} else {
		lk.St.Key[constant.KeyPost] = []interface{}{0}
	}
}

// GetList configs query
func (lk *LikesCount) GetList() *object.QuerySettings {
	qs := new(object.QuerySettings)
	qs.QueryName = constant.QueSelectLikeCountBy
	if value, ok := lk.St.Key[constant.KeyPost]; ok {
		qs.QueryFields = []interface{}{
			constant.TabPostsLikes,
			constant.TabPostsLikes,
			constant.TabPostsLikes,
			constant.FieldPost,
		}
		qs.Fields = value
	} else if value, ok := lk.St.Key[constant.KeyComment]; ok {
		qs.QueryFields = []interface{}{
			constant.TabCommentsLikes,
			constant.TabCommentsLikes,
			constant.TabCommentsLikes,
			constant.FieldComment,
		}
		qs.Fields = value
	} else {
		qs.QueryFields = []interface{}{
			constant.TabPostsLikes,
			constant.TabPostsLikes,
			constant.TabPostsLikes,
			constant.FieldPost,
		}
		qs.Fields = []interface{}{0} // for null list
	}
	return qs
}

// NewList returns fields for adding new data in "query rows"
func (lk *LikesCount) NewList() []interface{} {
	count := new(LCount)
	// true - it will show buttons
	// false - no buttons, only count
	count.Session = lk.Ck.Session
	// one post - it will access to save likes
	// all post - it will not access to save likes
	count.AllPost = lk.St.AllPost
	count.PostOrComm = lk.PostOrComm
	lk.Slice = append(lk.Slice, count)
	return []interface{}{
		&count.Like,
		&count.Body,
		&count.Count,
	}
}

func (lk *LikesCount) IfNil() interface{} {
	if lk.Ck.Session && !lk.St.AllPost {
		return []*LCount{
			lk.LikeNil(),
			lk.DislikeNil(),
		}
	}
	return []*LCount{new(LCount).ifNil()}
}

func (lk *LikesCount) Return() *Buf {
	return &Buf{LikesCount: lk}
}

func (lk *LikesCount) LikeNil() *LCount {
	return &LCount{
		Like:       1,
		Body:       constant.FieldLike,
		Count:      0,
		Session:    lk.Ck.Session,
		AllPost:    lk.St.AllPost,
		PostOrComm: lk.PostOrComm,
	}
}

func (lk *LikesCount) DislikeNil() *LCount {
	return &LCount{
		Like:       2,
		Body:       constant.FieldDislike,
		Count:      0,
		Session:    lk.Ck.Session,
		AllPost:    lk.St.AllPost,
		PostOrComm: lk.PostOrComm,
	}
}

package model

import (
	"github.com/giffone/forum-authentication/internal/constant"
	"github.com/giffone/forum-authentication/internal/object"
	"time"
)

type Like struct {
	PostOrComm int       // current post_id or comment_id
	Like       int       // like_id or dislike_id
	Body       string    // like or dislike (name)
	Created    time.Time // created date
	St         *object.Settings
	Ck         *object.Cookie
}

func NewLike(st *object.Settings, ck *object.Cookie) *Like {
	l := new(Like)
	if st == nil {
		l.St = &object.Settings{
			Key: make(map[string][]interface{}),
		}
	} else {
		l.St = st
	}
	if ck == nil {
		l.Ck = new(object.Cookie)
	} else {
		l.Ck = ck
	}
	return l
}

func (l *Like) MakeKeys(key string, data ...interface{}) {
	if key != "" {
		l.St.Key[key] = data
	} else {
		l.St.Key[constant.KeyPost] = []interface{}{0}
	}
}

func (l *Like) Get() *object.QuerySettings {
	qs := new(object.QuerySettings)
	if value, ok := l.St.Key[constant.KeyPost]; ok {
		qs.QueryName = constant.QueSelectLikeBy
		qs.QueryFields = []interface{}{
			constant.TabPostsLikes,
			constant.TabPostsLikes,
			constant.TabPostsLikes,
			constant.TabPostsLikes,
			constant.TabPostsLikes,
			constant.TabPostsLikes,
			constant.FieldUser,
			constant.TabPostsLikes,
			constant.FieldPost,
		}
		qs.Fields = value
	} else if value, ok := l.St.Key[constant.KeyComment]; ok {
		qs.QueryName = constant.QueSelectLikeBy
		qs.QueryFields = []interface{}{
			constant.TabCommentsLikes,
			constant.TabCommentsLikes,
			constant.TabCommentsLikes,
			constant.TabCommentsLikes,
			constant.TabCommentsLikes,
			constant.TabCommentsLikes,
			constant.FieldUser,
			constant.TabCommentsLikes,
			constant.FieldComment,
		}
		qs.Fields = value
	} else if value, ok := l.St.Key[constant.KeyPostRated]; ok {
		qs.QueryName = constant.QueSelectLikedOrNot
		qs.QueryFields = []interface{}{
			constant.TabPostsLikes,
			constant.TabPostsLikes,
			constant.TabPostsLikes,
			constant.FieldUser,
			constant.TabPostsLikes,
			constant.FieldPost,
		}
		qs.Fields = value
	} else if value, ok := l.St.Key[constant.KeyCommentRated]; ok {
		qs.QueryName = constant.QueSelectLikedOrNot
		qs.QueryFields = []interface{}{
			constant.TabCommentsLikes,
			constant.TabCommentsLikes,
			constant.TabCommentsLikes,
			constant.FieldUser,
			constant.TabCommentsLikes,
			constant.FieldComment,
		}
		qs.Fields = value
	} else { // for null list
		qs.QueryName = constant.QueSelectLikeBy
		qs.QueryFields = []interface{}{
			constant.TabPostsLikes,
			constant.TabPostsLikes,
			constant.TabPostsLikes,
			constant.FieldUser,
			constant.TabPostsLikes,
			constant.FieldPost,
		}
		qs.Fields = []interface{}{0, 0}
	}
	return qs
}

func (l *Like) New() []interface{} {
	if _, ok := l.St.Key[constant.KeyPostRated]; ok {
		return []interface{}{
			&l.Body,
		}
	} else if _, ok := l.St.Key[constant.KeyCommentRated]; ok {
		return []interface{}{
			&l.Body,
		}
	}
	return []interface{}{
		&l.PostOrComm,
		&l.Like,
		&l.Body,
		&l.Created,
	}
}

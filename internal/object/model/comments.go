package model

import (
	"github.com/giffone/forum-authentication/internal/constant"
	"github.com/giffone/forum-authentication/internal/object"
)

type Comments struct {
	Slice []*Comment
	St    *object.Settings
	Ck    *object.Cookie
}

func NewComments(st *object.Settings, ck *object.Cookie) *Comments {
	c := new(Comments)
	if st == nil {
		c.St = &object.Settings{
			Key: make(map[string][]interface{}),
		}
	} else {
		c.St = st
	}
	if ck == nil {
		c.Ck = new(object.Cookie)
	} else {
		c.Ck = ck
	}
	return c
}

func (c *Comments) MakeKeys(key string, data ...interface{}) {
	if key != "" {
		c.St.Key[key] = data
	} else {
		c.St.Key[constant.KeyComment] = []interface{}{0}
	}
}

func (c *Comments) GetList() *object.QuerySettings {
	qs := new(object.QuerySettings)
	if value, ok := c.St.Key[constant.KeyPost]; ok {
		qs.QueryName = constant.QueSelectCommentsBy
		qs.QueryFields = []interface{}{
			constant.TabComments,
			constant.FieldPost,
		}
		if value == nil {
			qs.Fields = []interface{}{c.Ck.Post}
		} else {
			qs.Fields = value
		}
	} else if value, ok := c.St.Key[constant.KeyUser]; ok {
		qs.QueryName = constant.QueSelectCommentsBy
		qs.QueryFields = []interface{}{
			constant.TabComments,
			constant.FieldUser,
		}
		if value == nil {
			qs.Fields = []interface{}{c.Ck.User}
		} else {
			qs.Fields = value
		}
	} else if value, ok := c.St.Key[constant.KeyRated]; ok {
		qs.QueryName = constant.QueSelectCommentsRatedBy
		qs.QueryFields = []interface{}{
			constant.TabCommentsLikes,
			constant.FieldUser,
		}
		if value == nil {
			qs.Fields = []interface{}{c.Ck.User}
		} else {
			qs.Fields = value
		}
	} else if value, ok := c.St.Key[constant.KeyID]; ok {
		qs.QueryName = constant.QueSelect
		qs.QueryFields = []interface{}{
			constant.TabComments,
			constant.TabComments,
			constant.FieldID,
		}
		if value == nil {
			qs.Fields = []interface{}{0}
		} else {
			qs.Fields = value
		}
	}
	return qs
}

func (c *Comments) NewList() []interface{} {
	comment := new(Comment)
	c.Slice = append(c.Slice, comment)
	if _, ok := c.St.Key[constant.KeyRated]; ok {
		return []interface{}{
			&comment.ID,
			&comment.Login,
			&comment.Body,
			&comment.Created,
			&comment.Post,
			&comment.Liked,
		}
	} else if _, ok := c.St.Key[constant.KeyID]; ok {
		return []interface{}{
			&comment.ID,
		}
	}
	return []interface{}{
		&comment.ID,
		&comment.Login,
		&comment.Body,
		&comment.Created,
		&comment.Post,
	}
}

func (c *Comments) IfNil() interface{} {
	return []*Comment{new(Comment).ifNil()}
}

func (c *Comments) Return() *Buf {
	return &Buf{Comments: c}
}

func (c *Comments) LSlice() int {
	return len(c.Slice)
}

func (c *Comments) PostOrCommentID(index int) int {
	return c.Slice[index].ID
}

func (c *Comments) Add(key string, index int, data interface{}) {
	switch key {
	case constant.KeyLike:
		c.Slice[index].Likes = data
	case constant.KeyRated:
		c.Slice[index].Liked = data
	}
}

func (c *Comments) Cookie() *object.Cookie {
	return c.Ck
}

func (c *Comments) Settings() *object.Settings {
	return c.St
}

func (c *Comments) KeyRole() string {
	return constant.KeyComment
}

func (c *Comments) KeyLiked() string {
	return constant.KeyCommentRated
}

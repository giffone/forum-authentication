package dto

import (
	"github.com/giffone/forum-authentication/internal/constant"
	"github.com/giffone/forum-authentication/internal/object"
	"net/http"
	"strings"
	"time"
)

type Comment struct {
	Body string
	Obj  object.Obj
}

func NewComment(st *object.Settings, sts *object.Statuses, ck *object.Cookie) *Comment {
	c := new(Comment)
	c.Obj.NewObjects(st, sts, ck)
	return c
}

func (c *Comment) Add(r *http.Request) bool {
	// get user id
	sts := c.Obj.Ck.CookieUserIDRead(r)
	if sts != nil {
		c.Obj.Sts = sts.Status()
		return false
	}
	// get post id
	sts = c.Obj.Ck.CookiePostIDRead(r)
	if sts != nil {
		c.Obj.Sts = sts.Status()
		return false
	}
	c.Body = r.PostFormValue("body text")
	return true
}

func (c *Comment) Valid() bool {
	// delete space for check an any symbol
	body := strings.TrimSpace(c.Body)
	if body == "" {
		c.Obj.Sts.StatusByText(nil, constant.TooShort,
			"comment", "one")
		return false
	}
	return true
}

func (c *Comment) Create() *object.QuerySettings {
	return &object.QuerySettings{
		QueryName: constant.QueInsert4,
		QueryFields: []interface{}{
			constant.TabComments,
			constant.FieldUser,
			constant.FieldPost,
			constant.FieldBody,
			constant.FieldCreated,
		},
		Fields: []interface{}{c.Obj.Ck.User,
			c.Obj.Ck.Post,
			c.Body,
			time.Now(),
		},
	}
}

func (c *Comment) Delete() *object.QuerySettings {
	return &object.QuerySettings{}
}

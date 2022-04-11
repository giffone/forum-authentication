package dto

import (
	"github.com/giffone/forum-authentication/internal/constant"
	"github.com/giffone/forum-authentication/internal/object"
	"net/http"
	"strings"
	"time"
)

type Post struct {
	Title      string
	Body       string
	Categories *CategoryPost
	Obj        object.Obj
}

func NewPost(st *object.Settings, sts *object.Statuses, ck *object.Cookie) *Post {
	p := new(Post)
	p.Categories = NewCategoryPost()
	p.Obj.NewObjects(st, sts, ck)
	return p
}

func (p *Post) Add(r *http.Request) bool {
	// get user id
	sts := p.Obj.Ck.CookieUserIDRead(r)
	if sts != nil {
		p.Obj.Sts = sts.Status()
		return false
	}
	p.Title = r.PostFormValue("title")
	p.Body = r.PostFormValue("body text")
	p.Categories.Slice = r.PostForm["categories"]
	return true
}

func (p *Post) Valid() bool {
	if p.Obj.Sts.StatusBody != "" {
		return false
	}
	// delete space for check an any symbol
	body := strings.TrimSpace(p.Body)
	if body == "" {
		p.Obj.Sts.StatusByText(nil, constant.TooShort, "post", "one")
		return false
	}
	body = strings.TrimSpace(p.Title)
	if body == "" {
		if len(p.Body) > 20 {
			p.Title = p.Body[0:19] + "..."
		} else {
			p.Title = p.Body
		}
	}
	return true
}

func (p *Post) Create() *object.QuerySettings {
	return &object.QuerySettings{
		QueryName: constant.QueInsert4,
		QueryFields: []interface{}{
			constant.TabPosts,
			constant.FieldUser,
			constant.FieldTitle,
			constant.FieldBody,
			constant.FieldCreated,
		},
		Fields: []interface{}{
			p.Obj.Ck.User,
			p.Title,
			p.Body,
			time.Now(),
		},
	}
}

func (p *Post) Delete() *object.QuerySettings {
	return &object.QuerySettings{}
}

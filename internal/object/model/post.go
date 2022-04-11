package model

import (
	"github.com/giffone/forum-authentication/internal/constant"
	"github.com/giffone/forum-authentication/internal/object"
	"time"
)

type Post struct {
	ID         int
	Title      string
	Body       string
	User       int
	Name       string
	Created    time.Time
	Categories interface{}
	Likes      interface{}
	Comments   *Comments
	Liked      interface{}
	St         *object.Settings
	Ck         *object.Cookie
}

func NewPost(st *object.Settings, ck *object.Cookie) *Post {
	p := new(Post)
	if st == nil {
		p.St = &object.Settings{
			Key: make(map[string][]interface{}),
		}
	} else {
		p.St = st
	}
	if ck == nil {
		p.Ck = new(object.Cookie)
	} else {
		p.Ck = ck
	}
	return p
}

func (p *Post) MakeKeys(key string, data ...interface{}) {
	if key != "" {
		p.St.Key[key] = data
	} else {
		p.St.Key[constant.KeyPost] = []interface{}{0}
	}
}

func (p *Post) Get() *object.QuerySettings {
	qs := new(object.QuerySettings)
	qs.QueryName = constant.QueSelectPostsBy
	qs.QueryFields = []interface{}{
		constant.TabPosts,
		constant.FieldID,
	}
	if value, ok := p.St.Key[constant.KeyPost]; ok {
		qs.Fields = value
	} else {
		qs.Fields = []interface{}{0} // for null list
	}
	return qs
}

func (p *Post) New() []interface{} {
	return []interface{}{
		&p.ID,
		&p.Title,
		&p.Body,
		&p.User,
		&p.Name,
		&p.Created,
	}
}

func (p *Post) IfNil() interface{} {
	return p.ifNil()
}

func (p *Post) ifNil() *Post {
	return &Post{
		Title:   "no posts created",
		Body:    "sorry, empty here",
		Created: time.Now(),
		User:    1,
		Name:    "Admin",
	}
}

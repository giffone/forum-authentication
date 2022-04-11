package model

import (
	"github.com/giffone/forum-authentication/internal/constant"
	"github.com/giffone/forum-authentication/internal/object"
)

type Posts struct {
	Slice []*Post
	St    *object.Settings
	Ck    *object.Cookie
}

func NewPosts(st *object.Settings, ck *object.Cookie) *Posts {
	p := new(Posts)
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

func (p *Posts) MakeKeys(key string, data ...interface{}) {
	if key != "" {
		p.St.Key[key] = data
	} else {
		p.St.Key[constant.KeyPost] = []interface{}{0}
	}
}

func (p *Posts) GetList() *object.QuerySettings {
	if len(p.St.Key) == 0 {
		return &object.QuerySettings{
			QueryName: constant.QueSelectPosts,
		}
	}
	qs := new(object.QuerySettings)
	if value, ok := p.St.Key[constant.KeyCategory]; ok {
		qs.QueryName = constant.QueSelectPostsAndCategoryBy
		qs.QueryFields = []interface{}{
			constant.TabPostsCategories,
			constant.FieldCategory,
		}
		qs.Fields = value
	} else if value, ok := p.St.Key[constant.KeyPost]; ok {
		qs.QueryName = constant.QueSelectPostsBy
		qs.QueryFields = []interface{}{
			constant.TabPosts,
			constant.FieldID,
		}
		if value == nil {
			qs.Fields = []interface{}{p.Ck.Post}
		} else {
			qs.Fields = value
		}
	} else if value, ok := p.St.Key[constant.KeyUser]; ok {
		qs.QueryName = constant.QueSelectPostsBy
		qs.QueryFields = []interface{}{
			constant.TabPosts,
			constant.FieldUser,
		}
		if value == nil {
			qs.Fields = []interface{}{p.Ck.User}
		} else {
			qs.Fields = value
		}
	} else if value, ok := p.St.Key[constant.KeyID]; ok {
		qs.QueryName = constant.QueSelect
		qs.QueryFields = []interface{}{
			constant.TabPosts,
			constant.TabPosts,
			constant.FieldID,
		}
		if value == nil {
			qs.Fields = []interface{}{0}
		} else {
			qs.Fields = value
		}
	} else if value, ok := p.St.Key[constant.KeyRated]; ok {
		qs.QueryName = constant.QueSelectPostsRatedBy
		qs.QueryFields = []interface{}{
			constant.TabPostsLikes,
			constant.FieldUser,
		}
		if value == nil {
			qs.Fields = []interface{}{p.Ck.User}
		} else {
			qs.Fields = value
		}
	}
	return qs
}

func (p *Posts) NewList() []interface{} {
	post := new(Post)
	p.Slice = append(p.Slice, post)
	// for account handler
	if _, ok := p.St.Key[constant.KeyComment]; ok {
		comment := new(Comment)
		comments := new(Comments)
		comments.Slice = append(comments.Slice, comment)
		post.Comments = comments
		return []interface{}{
			&post.ID,
			&post.Title,
			&post.Body,
			&post.User,
			&post.Name,
			&post.Created,
			&comment.ID,
			&comment.Name,
			&comment.Body,
			&comment.Created,
		}
	} else if _, ok := p.St.Key[constant.KeyRated]; ok {
		return []interface{}{
			&post.ID,
			&post.Title,
			&post.Body,
			&post.User,
			&post.Name,
			&post.Created,
			&post.Liked,
		}
	} else if _, ok := p.St.Key[constant.KeyID]; ok {
		return []interface{}{
			&post.ID,
		}
	}
	return []interface{}{
		&post.ID,
		&post.Title,
		&post.Body,
		&post.User,
		&post.Name,
		&post.Created,
	}
}

func (p *Posts) IfNil() interface{} {
	return []*Post{new(Post).ifNil()}
}

func (p *Posts) Return() *Buf {
	return &Buf{Posts: p}
}

func (p *Posts) LSlice() int {
	return len(p.Slice)
}

func (p *Posts) PostOrCommentID(index int) int {
	return p.Slice[index].ID
}

// Add adding information to slice post/comment by index
func (p *Posts) Add(key string, index int, data interface{}) {
	switch key {
	case constant.KeyLike:
		p.Slice[index].Likes = data
	case constant.KeyRated:
		p.Slice[index].Liked = data
	case constant.KeyCategory:
		p.Slice[index].Categories = data
	}
}

func (p *Posts) Cookie() *object.Cookie {
	return p.Ck
}

func (p *Posts) Settings() *object.Settings {
	return p.St
}

func (p *Posts) KeyRole() string {
	return constant.KeyPost
}

func (p *Posts) KeyLiked() string {
	return constant.KeyPostRated
}

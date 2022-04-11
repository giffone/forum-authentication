package model

import (
	"time"
)

type Comment struct {
	ID      int
	Title   *Post // refer to post
	Body    string
	User    int
	Name    string
	Created time.Time
	Likes   interface{}
	Liked   interface{}
	Post    int
}

func (c *Comment) IfNil() interface{} {
	return c.ifNil()
}

func (c *Comment) ifNil() *Comment {
	return &Comment{
		ID: 0,
	}
}

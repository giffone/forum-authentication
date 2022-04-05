package model

import (
	"github.com/giffone/forum-authentication/internal/constant"
	"github.com/giffone/forum-authentication/internal/object"
	"time"
)

type Session struct {
	User   string
	UUID   string
	Expire time.Time
	St     *object.Settings // settings
	Ck     *object.Cookie
}

func NewSession(st *object.Settings, ck *object.Cookie) *Session {
	s := new(Session)
	if st == nil {
		s.St = &object.Settings{
			Key: make(map[string][]interface{}),
		}
	} else {
		s.St = st
	}
	if ck == nil {
		s.Ck = new(object.Cookie)
	} else {
		s.Ck = ck
	}
	return s
}

func (s *Session) Get() *object.QuerySettings {
	return &object.QuerySettings{
		QueryName: constant.QueSelectSessionBy,
		QueryFields: []interface{}{
			constant.FieldUser,
			constant.FieldUUID,
		},
		Fields: []interface{}{
			s.Ck.User,
			s.Ck.SessionUUID,
		},
	}
}

func (s *Session) New() []interface{} {
	return []interface{}{
		&s.User,
		&s.UUID,
		&s.Expire,
	}
}
